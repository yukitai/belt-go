package frontend

import (
	"belt/reporter"
	"fmt"
)

type TyInfer struct {
	analyzer *Analyzer
	globals  map[string]*AstValType
	funcs    []*AstFnDecl
}

func TyInferNew(a *Analyzer) TyInfer {
	return TyInfer{
		analyzer: a,
		globals:  make(map[string]*AstValType),
		funcs:    make([]*AstFnDecl, 0),
	}
}

func (i *TyInfer) InitGlobal() {
	for _, item := range i.analyzer.ast.Items {
		switch item.Stype {
		case ANSFndecl:
			fn := item.Item.(*AstFnDecl)
			types := make([]AstValType, 0)
			for _, arg := range fn.Args {
				types = append(types, arg.Atype)
			}
			i.globals[fn.Name.value] = &AstValType{
				Vttype: ANTFnType,
				Item: &AstValTypeFnType{
					Types: types,
					Ret_t: fn.Ret_t,
				},
			}
			i.funcs = append(i.funcs, fn)
		default:
		}
	}
}

func (i *TyInfer) InferAll() {
	for _, fn := range i.funcs {
		i.InferFunc(fn)
	}
}

func (i *TyInfer) InferFunc(fn *AstFnDecl) {
	ufs := UFSNew()
	unit := AstTupleNew()
	last := ufs.Extend(&unit) // Empty Ret Type
	ret_t := ufs.Extend(&fn.Ret_t)
	idents := map[string]uint{
		"$$ret": ret_t,
	}
	for _, arg := range fn.Args {
		idents[arg.Name.value] = ufs.Extend(&arg.Atype)
	}
	for _, item := range fn.Body.Items {
		last = i.InferStmt(&ufs, idents, &item)
	}
	ufs.Merge(last, ret_t)

	ufs.MakeEffect(i.analyzer)
	// fmt.Printf("Parents: %v\nRanks:   %v\nValues:  %v\n", ufs.parents, ufs.ranks, ufs.values)
}

func (i *TyInfer) InferStmt(ufs *UnionFindSet, idents map[string]uint, stmt *AstStmt) uint {
	switch stmt.Stype {
	case ANSExpr:
		expr_s := stmt.Item.(*AstExpr)
		return i.InferExpr(ufs, idents, expr_s)
	case ANSLet:
		let_s := stmt.Item.(*AstLetStmt)
		val_t := ufs.Extend(&let_s.Vtype)
		if let_s.Expr != nil {
			ufs.Merge(val_t, i.InferExpr(ufs, idents, let_s.Expr))
		}
		idents[let_s.Name.value] = val_t
		return 0
	case ANSReturn: /* @todo replace this type with never(!) */
		ret_s := stmt.Item.(*AstReturnStmt)
		ufs.Merge(idents["$$ret"], i.InferExpr(ufs, idents, &ret_s.Expr))
		return 0
	default:
		return 0
	}
}

var type_int = AstValType{
	Vttype: ANTBinary,
	Item: &AstValTypeBinary{
		Tok_type: &Token{
			ttype: KTInt,
			value: "int",
		},
	},
}

var type_float = AstValType{
	Vttype: ANTBinary,
	Item: &AstValTypeBinary{
		Tok_type: &Token{
			ttype: KTFloat,
			value: "float",
		},
	},
}

var type_bool = AstValType{
	Vttype: ANTBinary,
	Item: &AstValTypeBinary{
		Tok_type: &Token{
			ttype: KTBool,
			value: "bool",
		},
	},
}

var type_string = AstValType{
	Vttype: ANTBinary,
	Item: &AstValTypeBinary{
		Tok_type: &Token{
			ttype: KTString,
			value: "string",
		},
	},
}

func (i *TyInfer) InferExpr(ufs *UnionFindSet, idents map[string]uint, expr *AstExpr) uint {
	switch expr.Etype {
	case ANELiteral:
		expr := expr.Item.(*AstExprLiteral)
		switch expr.Value.ttype {
		case LlInt:
			return ufs.Extend(&type_int)
		case LlFloat:
			return ufs.Extend(&type_float)
		case LlBool:
			return ufs.Extend(&type_bool)
		case LlString:
			return ufs.Extend(&type_string)
		default:
			panic("reaching an unreachable code! something went wrong")
		}
	case ANEVar:
		expr := expr.Item.(*AstExprVar)
		var_t, ok := idents[expr.Ident.value]
		if !ok {
			err := reporter.Error(
				expr.Where(),
				fmt.Sprintf("cannot resolve name `%v`", expr.Ident.value),
			)
			reporter.Report(&err, i.analyzer.file)
			i.analyzer.has_err = true
		}
		return var_t
	case ANEGroup:
		expr := expr.Item.(*AstExprGroup)
		return i.InferExpr(ufs, idents, &expr.Expr)
	case ANEOp1:
		panic("not implemented yet")
	case ANEOp2:
		expr := expr.Item.(*AstExprOp2)
		lhs_t := i.InferExpr(ufs, idents, &expr.Lhs)
		rhs_t := i.InferExpr(ufs, idents, &expr.Rhs)
		var expr_t uint
		switch expr.Op.ttype {
		case OAdd, OSub, OMul, ODiv:
			ufs.Merge(lhs_t, rhs_t)
			expr_t = ufs.Extend(&type_int)
			ufs.Merge(expr_t, lhs_t)
		case OAddf, OSubf, OMulf, ODivf:
			ufs.Merge(lhs_t, rhs_t)
			expr_t = ufs.Extend(&type_float)
			ufs.Merge(expr_t, lhs_t)
		case OConnect:
			ufs.Merge(lhs_t, rhs_t)
			expr_t = ufs.Extend(&type_string)
			ufs.Merge(expr_t, lhs_t)
		case OEq, ONeq, OGrt, OGeq, OLes, OLeq:
			ufs.Merge(lhs_t, rhs_t)
			expr_t = ufs.Extend(&type_bool)
		case OAnd, OOr:
			ufs.Merge(lhs_t, rhs_t)
			expr_t = ufs.Extend(&type_bool)
			ufs.Merge(expr_t, lhs_t)
		case OAssign:
			expr_t = ufs.ExtendTVar()
			ufs.Merge(expr_t, rhs_t)
		default:
			panic("not implemented yet")
		}
		return expr_t
	default:
		panic("not implemented yet")
	}
}