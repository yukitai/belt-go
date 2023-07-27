package frontend

import (
	"belt/reporter"
	"belt/utils"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type AstLLVMBuilder struct {
	ast AstFile
	m *ir.Module
	file *utils.File
}

func AstLLVMBuilderNew(ast AstFile, file *utils.File) AstLLVMBuilder {
	return AstLLVMBuilder{
		ast: ast,
		file: file,
	}
}

func (b *AstLLVMBuilder) Build() *ir.Module {
	m := ir.NewModule()
	b.m = m

	builtins := map[string]*ir.Func{}

	// extern functions
	builtins["puts"] = m.NewFunc("puts", types.I32, ir.NewParam("", types.NewPointer(types.I8)))

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

func (b *AstLLVMBuilder) BuildStmtFnLocal(fn *Func, a *AstStmt) {
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
			Type: &ty,
		})
	case ANSExpr:
		b.BuildExprFnLocal(fn, a.Item.(*AstExpr))
	case ANSReturn:
	case ANSBreak:
	case ANSContinue:
	default:
		err := reporter.Error(
			a.Item.Where(),
			"unexpected item in the local scope",
		)
		reporter.Report(&err, b.file)
	}
}

type Func struct {
	fn     *ir.Func
	params []*ir.Param
	block  *ir.Block
	idents map[string]*Identifier
}

type Identifier struct {
	Type  *types.Type
	Value value.Value
}

func FuncNew(fn *ir.Func, params []*ir.Param) Func {
	return Func{
		fn: fn,
		params: params,
		block: fn.NewBlock("entry"),
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

func (b *AstLLVMBuilder) BuildBlockFnLocal(fn *Func, a *AstBlock) {
	for _, item := range a.Items {
		b.BuildStmtFnLocal(fn, &item)
	}
}

func (b *AstLLVMBuilder) BuildFnBody(fn *Func, a *AstFnDecl) {
	for _, param := range fn.params {
		addr := b.BuildFnParam(fn.block, param)
		fn.addIdent(param.Name(), &Identifier{
			Value: addr,
			Type: &param.Typ,
		})
	}

	b.BuildBlockFnLocal(fn, &a.Body)
	
	term := ir.NewRet(nil)
	fn.block.Term = term
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
		panic("not implemented yet")
	case ANEBlock:
		return b.BuildExprFnLocalBlock(fn, e.Item.(*AstBlock))
	case ANEClosure:
		return b.BuildExprFnLocalClosure(fn, e.Item.(*AstClosure))
	case ANEFncall:
		return b.BuildExprFnLocalFncall(fn, e.Item.(*AstExprFncall))
	case ANEForIn:
		panic("not implemented yet")
	case ANEWhile:
		panic("not implemented yet")
	case ANEGroup:
		return b.BuildExprFnLocal(fn, &e.Item.(*AstExprGroup).Expr)
	case ANEBuiltinCorePrint:
		panic("not implemented yet")
	}
	c := constant.NewInt(types.I64, 114514)
	return c
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
		return constant.NewAddrSpaceCast(str, types.I8Ptr)
	default:
		panic("reaching an unreachable code! something went wrong")
	}
}

func (b *AstLLVMBuilder) BuildExprFnLocalVar(fn *Func, e *AstExprVar) value.Value {
	ident :=fn.lookupIdent(e.Ident.value)
	return fn.block.NewLoad(*ident.Type, ident.Value).Src
}

func (b *AstLLVMBuilder) BuildExprFnLocalBlock(fn *Func, e *AstBlock) value.Value {
	panic("not implemented yet")
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
	panic("not implemented yet")
}