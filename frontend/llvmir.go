package frontend

import (
	"belt/reporter"
	"belt/utils"
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type AstLLVMBuilder struct {
	ast      AstFile
	m        *ir.Module
	file     *utils.File
	pool     map[string]*ir.Global
	builtins map[string]*ir.Func
}

func AstLLVMBuilderNew(ast AstFile, file *utils.File) AstLLVMBuilder {
	return AstLLVMBuilder{
		ast:      ast,
		file:     file,
		pool:     make(map[string]*ir.Global),
		builtins: make(map[string]*ir.Func),
	}
}

func (b *AstLLVMBuilder) Build() *ir.Module {
	m := ir.NewModule()
	b.m = m

	builtins := map[string]*ir.Func{}
	b.builtins = builtins

	// extern functions
	builtins["puts"] = m.NewFunc("puts", types.I32, ir.NewParam("", types.I8Ptr))

	// bl_builtins
	builtins["__bl_str_connect"] = m.NewFunc("__bl_str_connect", types.I8Ptr, ir.NewParam("", types.I8Ptr), ir.NewParam("", types.I8Ptr))

	b.BuildFile(&b.ast)

	return m
}

func (b *AstLLVMBuilder) BuildFile(a *AstFile) {
	for _, item := range a.Items {
		b.BuildItemGlobal(&item)
	}
}

func (b *AstLLVMBuilder) BuildItemGlobal(a *AstItem) {
	switch a.Stype {
	case ANSFndecl:
		b.BuildFndeclGlobal(a.Item.(*AstFnDecl))
	default:
		err := reporter.Error(
			a.Item.Where(),
			"unexpected item in the top-level scope",
		)
		reporter.Report(&err, b.file)
	}
}

func (b *AstLLVMBuilder) BuildStmtFnLocal(fn *Func, a *AstStmt) *value.Value {
	switch a.Stype {
	case ANSLet:
		stmt := a.Item.(*AstLetStmt)
		ty := b.BuildType(&stmt.Vtype)
		addr := fn.block.NewAlloca(ty)
		if stmt.Expr != nil {
			value := b.BuildExprFnLocal(fn, stmt.Expr)
			fn.block.NewStore(value, addr)
		}
		fn.addIdent(stmt.Name.value, &Identifier{
			Value: addr,
		})
	case ANSExpr:
		val := b.BuildExprFnLocal(fn, a.Item.(*AstExpr))
		return &val
	case ANSReturn:
		stmt := a.Item.(*AstReturnStmt)
		ret := b.BuildExprFnLocal(fn, &stmt.Expr)
		fn.block.Term = &ir.TermRet{
			X: ret,
		}
	case ANSBreak:
	case ANSContinue:
	default:
		err := reporter.Error(
			a.Item.Where(),
			"unexpected item in the local scope",
		)
		reporter.Report(&err, b.file)
	}
	return nil
}

type Func struct {
	fn     *ir.Func
	params []*ir.Param
	block  *ir.Block
	idents map[string]*Identifier
}

type Identifier struct {
	Value value.Value
}

func FuncNew(fn *ir.Func, params []*ir.Param) Func {
	return Func{
		fn:     fn,
		params: params,
		block:  fn.NewBlock("entry"),
		idents: make(map[string]*Identifier, 0),
	}
}

func (b *AstLLVMBuilder) BuildFndeclGlobal(a *AstFnDecl) {
	params := make([]*ir.Param, 0)
	for _, arg := range a.Args {
		params = append(params, ir.NewParam(arg.Name.value, b.BuildType(&arg.Atype)))
	}
	llir_fn := b.m.NewFunc(a.Name.value, b.BuildType(&a.Ret_t), params...)
	fn := FuncNew(llir_fn, params)
	b.BuildFnBody(&fn, a)
}

func (b *AstLLVMBuilder) BuildType(a *AstValType) types.Type {
	switch a.Vttype {
	case ANTBinary:
		switch a.Item.(*AstValTypeBinary).Tok_type.ttype {
		case KTInt:
			return types.I64
		case KTFloat:
			return types.FP128
		case KTBool:
			return types.I1
		case KTString:
			return types.I8Ptr
		default:
			panic("reaching an unreachable code! something went wrong")
		}
	case ANTStruct:
		panic("not implemented yet")
	case ANTEnum:
		panic("not implemented yet")
	case ANTUnknown:
		ty := a.Item.(*AstUnkownType)
		if ty.Rtype != nil {
			return b.BuildType(ty.Rtype)
		}
		panic("reaching an unreachable code! something went wrong")
	case ANTVar:
		ty := a.Item.(*AstValTypeVar)
		if ty.Real != nil {
			return b.BuildType(ty.Real)
		}
		panic("reaching an unreachable code! something went wrong")
	default:
		panic("reaching an unreachable code! something went wrong")
	}
}

func (fn *Func) addIdent(name string, ident *Identifier) {
	fn.idents[name] = ident
}

func (fn *Func) lookupIdent(name string) *Identifier {
	it, ok := fn.idents[name]
	if !ok {
		panic("reaching an unreachable code! something went wrong")
	}
	return it
}

func (b *AstLLVMBuilder) BuildBlockFnLocal(fn *Func, a *AstBlock) value.Value {
	/*for _, item := range a.Items {
		b.BuildStmtFnLocal(fn, &item)
	}*/
	block := b.BuildExprFnLocalBlock(fn, a)
	return block.Ret
}

func (b *AstLLVMBuilder) BuildFnBody(fn *Func, a *AstFnDecl) {
	for _, param := range fn.params {
		addr := b.BuildFnParam(fn.block, param)
		fn.addIdent(param.Name(), &Identifier{
			Value: addr,
		})
	}

	ret := b.BuildBlockFnLocal(fn, &a.Body)
	fn.block.Term = ir.NewRet(ret)
}

func (b *AstLLVMBuilder) BuildFnParam(block *ir.Block, param *ir.Param) value.Value {
	addr := block.NewAlloca(param.Type())
	block.NewStore(param, addr)
	return addr
}

func (b *AstLLVMBuilder) BuildExprFnLocal(fn *Func, e *AstExpr) value.Value {
	switch e.Etype {
	case ANELiteral:
		return b.BuildExprFnLocalLiteral(fn, e.Item.(*AstExprLiteral))
	case ANEVar:
		return b.BuildExprFnLocalVar(fn, e.Item.(*AstExprVar))
	case ANEOp1:
		return b.BuildExprFnLocalOp1(fn, e.Item.(*AstExprOp1))
	case ANEOp2:
		return b.BuildExprFnLocalOp2(fn, e.Item.(*AstExprOp2))
	case ANEIfElse:
		return b.BuildExprFnLocalIfElse(fn, e.Item.(*AstExprIfElse))
	case ANEBlock:
		return b.BuildExprFnLocalBlock(fn, e.Item.(*AstBlock)).Ret
	case ANEClosure:
		return b.BuildExprFnLocalClosure(fn, e.Item.(*AstClosure))
	case ANEFncall:
		return b.BuildExprFnLocalFncall(fn, e.Item.(*AstExprFncall))
	case ANEForIn:
		panic("not implemented yet")
	case ANEWhile:
		return b.BuildExprFnLocalWhile(fn, e.Item.(*AstExprWhile))
	case ANEGroup:
		return b.BuildExprFnLocal(fn, &e.Item.(*AstExprGroup).Expr)
	case ANEBuiltinCorePrint:
		panic("not implemented yet")
	}
	c := constant.NewInt(types.I64, 114514)
	return c
}

var pool_id uint

func GetPoolId() string {
	id := fmt.Sprintf("$const_%v", pool_id)
	pool_id += 1
	return id
}

func (b *AstLLVMBuilder) GetFromPoolOrCreate(src constant.Constant) *ir.Global {
	key := src.String()
	def, ok := b.pool[key]
	if !ok {
		def = b.m.NewGlobalDef(GetPoolId(), src)
		b.pool[key] = def
	}
	return def
}

func (b *AstLLVMBuilder) GetPool() map[string]*ir.Global {
	return b.pool
}

func (b *AstLLVMBuilder) BuildExprFnLocalLiteral(fn *Func, e *AstExprLiteral) value.Value {
	switch e.Value.ttype {
	case LlBool:
		return constant.NewBool(e.Value.AsBool(b.file))
	case LlFloat:
		return constant.NewFloat(types.FP128, e.Value.AsFloat(b.file))
	case LlInt:
		return constant.NewInt(types.I64, int64(e.Value.AsInt(b.file)))
	case LlNil:
		return constant.NewInt(types.I64, 0)
	case LlString:
		str := constant.NewCharArrayFromString(e.Value.value)
		def := b.GetFromPoolOrCreate(str)
		zero := constant.NewInt(types.I64, 0)
		return constant.NewGetElementPtr(def.ContentType, def, zero, zero)
	default:
		panic("reaching an unreachable code! something went wrong")
	}
}

func (b *AstLLVMBuilder) BuildExprFnLocalVar(fn *Func, e *AstExprVar) value.Value {
	ident := fn.lookupIdent(e.Ident.value)
	return fn.block.NewLoad(ident.Value.Type().(*types.PointerType).ElemType, ident.Value)
}

type Block struct {
	Block *ir.Block
	Ret   value.Value
}

func (b *AstLLVMBuilder) BuildExprFnLocalBlock(fn *Func, e *AstBlock) Block {
	var val *value.Value
	var ret value.Value
	for _, item := range e.Items {
		val = b.BuildStmtFnLocal(fn, &item)
	}
	if val == nil {
		ret = constant.NewStruct(types.NewStruct())
	} else {
		ret = *val
	}
	return Block{
		Block: fn.block,
		Ret:   ret,
	}
}

func (b *AstLLVMBuilder) BuildExprFnLocalClosure(fn *Func, e *AstClosure) value.Value {
	panic("not implemented yet")
}

func (b *AstLLVMBuilder) BuildExprFnLocalFncall(fn *Func, e *AstExprFncall) value.Value {
	panic("not implemented yet")
}

func (b *AstLLVMBuilder) BuildExprFnLocalOp1(fn *Func, e *AstExprOp1) value.Value {
	panic("not implemented yet")
}

func (b *AstLLVMBuilder) BuildExprFnLocalOp2(fn *Func, e *AstExprOp2) value.Value {
	lhs := b.BuildExprFnLocal(fn, &e.Lhs)
	rhs := b.BuildExprFnLocal(fn, &e.Rhs)
	switch e.Op.ttype {
	case OAdd:
		return fn.block.NewAdd(lhs, rhs)
	case OAddf:
		return fn.block.NewFAdd(lhs, rhs)
	case OConnect:
		return fn.block.NewCall(b.builtins["__bl_str_connect"], lhs, rhs)
	case OSub:
		return fn.block.NewSub(lhs, rhs)
	case OSubf:
		return fn.block.NewFSub(lhs, rhs)
	case OMul:
		return fn.block.NewMul(lhs, rhs)
	case OMulf:
		return fn.block.NewFMul(lhs, rhs)
	case ODiv:
		return fn.block.NewSDiv(lhs, rhs)
	case ODivf:
		return fn.block.NewFDiv(lhs, rhs)
	case OEq:
		return fn.block.NewICmp(enum.IPredEQ, lhs, rhs)
	case ONeq:
		return fn.block.NewICmp(enum.IPredNE, lhs, rhs)
	case OGrt:
		return fn.block.NewICmp(enum.IPredSGT, lhs, rhs)
	case OGeq:
		return fn.block.NewICmp(enum.IPredSGE, lhs, rhs)
	case OLes:
		return fn.block.NewICmp(enum.IPredSLT, lhs, rhs)
	case OLeq:
		return fn.block.NewICmp(enum.IPredSLE, lhs, rhs)
	case OAnd:
		return fn.block.NewAnd(lhs, rhs)
	case OOr:
		return fn.block.NewOr(lhs, rhs)
	case OBXor:
		return fn.block.NewXor(lhs, rhs)
	case OBAnd:
		panic("not implemented yet")
	case OBOr:
		panic("not implemented yet")
	case OMovl:
		panic("not implemented yet")
	case OMovr:
		panic("not implemented yet")
	case OMember:
		panic("not implemented yet")
	case OLookup:
		panic("not implemented yet")
	case OAssign:
		panic("not implemented yet")
	default:
		panic("reaching an unreachable code! something went wrong")
	}
}

func (b *AstLLVMBuilder) BuildExprFnLocalIfElse(fn *Func, e *AstExprIfElse) value.Value {
	if e.Case2 == nil {
		return b.BuildExprFnLocalIf(fn, e)
	}
	cond := b.BuildExprFnLocal(fn, &e.Cond)
	case1 := fn.fn.NewBlock("ifelse_true_case")
	case2 := fn.fn.NewBlock("ifelse_false_case")
	fn.block.Term = ir.NewCondBr(cond, case1, case2)
	outer := fn.fn.NewBlock("ifelse_outer")
	case1.Term = &ir.TermBr{
		Target: outer,
	}
	case2.Term = &ir.TermBr{
		Target: outer,
	}
	ret := fn.block.NewAlloca(b.BuildType(&e.Type))
	fn.block = case1
	case1_val := b.BuildBlockFnLocal(fn, &e.Case1)
	fn.block.NewStore(case1_val, ret)
	fn.block = case2
	case2_val := b.BuildBlockFnLocal(fn, e.Case2)
	fn.block.NewStore(case2_val, ret)
	fn.block = outer
	return ret
}

func (b *AstLLVMBuilder) BuildExprFnLocalIf(fn *Func, e *AstExprIfElse) value.Value {
	panic("not implemented yet")
}

func (b *AstLLVMBuilder) BuildExprFnLocalWhile(fn *Func, e *AstExprWhile) value.Value {
	panic("not implemented yet")
}