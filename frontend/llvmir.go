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
	file *utils.File
}

func AstLLVMBuilderNew(ast AstFile, file *utils.File) AstLLVMBuilder {
	return AstLLVMBuilder{
		ast, file,
	}
}

func (b *AstLLVMBuilder) Build() *ir.Module {
	m := ir.NewModule()

	builtins := map[string]*ir.Func{}

	// extern functions
	builtins["puts"] = m.NewFunc("puts", types.I32, ir.NewParam("", types.NewPointer(types.I8)))

	b.BuildFile(m, &b.ast)

	return m
}

func (b *AstLLVMBuilder) BuildFile(m *ir.Module, a *AstFile) {
	for _, item := range a.Items {
		b.BuildItemGlobal(m, &item)
	}
}

func (b *AstLLVMBuilder) BuildItemGlobal(m *ir.Module, a *AstItem) {
	switch a.Stype {
	case ANSFndecl:
		b.BuildFndeclGlobal(m, a.Item.(*AstFnDecl))
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
		addr := fn.block.NewAlloca(b.BuildType(&stmt.Vtype))
		if stmt.Expr != nil {
			value := b.BuildExprFnLocal(fn, stmt.Expr)
			fn.block.NewStore(value, addr)
		}
	case ANSExpr:
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
	idents map[string]value.Value
}

func FuncNew(fn *ir.Func, params []*ir.Param) Func {
	return Func{
		fn: fn,
		params: params,
		block: fn.NewBlock("entry"),
		idents: make(map[string]value.Value, 0),
	}
}

func (b *AstLLVMBuilder) BuildFndeclGlobal(m *ir.Module, a *AstFnDecl) {
	params := make([]*ir.Param, 0)
	for _, arg := range a.Args {
		params = append(params, ir.NewParam(arg.Name.value, b.BuildType(&arg.Atype)))
	}
	llir_fn := m.NewFunc(a.Name.value, b.BuildType(&a.Ret_t), params...)
	fn := FuncNew(llir_fn, params)
	b.BuildFnBody(&fn, a)
}

func (b *AstLLVMBuilder) BuildType(a *AstValType) types.Type {
	switch a.Vttype {
	default:
		return types.Void
	}
}

func (b *AstLLVMBuilder) BuildBlockFnLocal(fn *Func, a *AstBlock) {
	for _, item := range a.Items {
		b.BuildStmtFnLocal(fn, &item)
	}
}

func (b *AstLLVMBuilder) BuildFnBody(fn *Func, a *AstFnDecl) {
	for _, param := range fn.params {
		addr := b.BuildFnParam(fn.block, param)
		fn.idents[param.LocalName] = addr
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
	c := constant.NewInt(types.I64, 114514)
	return c
}