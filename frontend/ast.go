package frontend

import (
	"belt/reporter"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type AstType int

const (
	ANFile AstType = iota
	
	ANFndecl
	ANClosure
	ANFnArg

	ANExpr
	ANExprOp1
	ANExprOp2
	ANExprGroup
	ANExprLiteral
	ANExprVar
	ANExprFncall
	ANExprBuiltinCorePrint

	ANStmt
	ANStmtLet
	ANStmtReturn
	ANStmtBreak
	ANStmtContinue

	// ANBuiltinName

	ANValType
	ANTypeBinary
	ANTypeVar
	ANTypeStruct
	ANTypeEnum
	ANTypeFnType
	ANTypeTuple

	ANBlock
)

type AstExprType int

const (
	ANEOp2 AstExprType = iota
	ANEOp1
	ANEGroup
	ANELiteral
	ANEVar
	ANEIfElse
	ANEWhile
	ANEForIn
	ANEClosure
	ANEBlock
	ANEBuiltinCorePrint
	ANEFncall
)

type AstValTypeType int

const (
	ANTBinary AstValTypeType = iota
	ANTVar
	ANTUnknown
	ANTStruct
	ANTEnum
	ANTFnType
	ANTTuple
	// ANTTuple
)

type AstStmtType int

const (
	ANSExpr AstStmtType = iota
	ANSLet
	ANSReturn
	ANSBreak
	ANSFndecl
	ANSContinue
)

func ident(x uint) {
	if x > 0 {
		fmt.Printf("%v+ ", strings.Repeat("  ", int(x-1)))
	}
}

var debugNilString = color.New(color.Italic, color.Bold, color.FgHiCyan).Sprintf("nil")

func debug_token(x uint, tok *Token) {
	ident(x)
	if tok != nil {
		fmt.Printf("Token { %v, `%v` }\n", tok.ttype, tok.value)
	} else {
		fmt.Printf("Token { %v }\n", debugNilString)
	}
}

type AstNode interface {
	ANType() AstType
	Debug(uint)
	Where() reporter.Where
}

type AstFile struct {
	Items []AstItem
}

func (a *AstFile) ANType() AstType {
	return ANFile
}

func (a *AstFile) Debug(x uint) {
	ident(x)
	fmt.Printf("AstFile\n")
	for _, item := range a.Items {
		item.Debug(x + 1)
	}
}

func (a *AstFile) Where() reporter.Where {
	return reporter.WhereNew(1, 1, 0, 0)
}

type AstItem = AstStmt

type AstFnDecl struct {
	Tok_fn      *Token
	Name        *Token
	Tok_lbrace  *Token
	Args        []AstFnArg
	Tok_rbrace  *Token
	Tok_thinarr *Token
	Ret_t       AstValType
	Body        AstBlock
}

func (a *AstFnDecl) ANType() AstType {
	return ANFndecl
}

func (a *AstFnDecl) Debug(x uint) {
	ident(x)
	fmt.Printf("AstFnDecl\n")
	debug_token(x+1, a.Tok_fn)
	debug_token(x+1, a.Name)
	debug_token(x+1, a.Tok_lbrace)
	for _, arg:= range a.Args {
		arg.Debug(x + 1)
	}
	debug_token(x+1, a.Tok_rbrace)
	a.Ret_t.Debug(x + 1)
	a.Body.Debug(x + 1)
}

func (a *AstFnDecl) Where() reporter.Where {
	return a.Tok_fn.where.Merge(&a.Body.Tok_rbra.where)
}

type AstClosure struct {
	Tok_lbor    *Token
	Args        []AstFnArg
	Tok_rbor    *Token
	Tok_thinarr *Token
	Ret_t       AstValType
	Body        AstExpr
}

func (a *AstClosure) ANType() AstType {
	return ANFndecl
}

func (a *AstClosure) Debug(x uint) {
	ident(x)
	fmt.Printf("AstClosure\n")
	debug_token(x+1, a.Tok_lbor)
	ident(x + 1)
	fmt.Printf("Arguments\n")
	for _, arg := range a.Args {
		arg.Debug(x + 2)
	}
	debug_token(x+1, a.Tok_rbor)
	a.Ret_t.Debug(x + 1)
	a.Body.Debug(x + 1)
}

func (a *AstClosure) Where() reporter.Where {
	body := a.Body.Item.Where()
	return a.Tok_lbor.where.Merge(&body)
}

type AstFnArg struct {
	Name      *Token
	Tok_colon *Token
	Atype     AstValType
	Tok_comma *Token
}

func (a *AstFnArg) ANType() AstType {
	return ANFnArg
}

func (a *AstFnArg) Debug(x uint) {
	ident(x)
	fmt.Printf("AstFnArg\n")
	debug_token(x+1, a.Name)
	debug_token(x+1, a.Tok_colon)
	a.Atype.Debug(x + 1)
	debug_token(x+1, a.Tok_comma)
}

func (a *AstFnArg) Where() reporter.Where {
	return a.Name.where
}

type AstValType struct {
	Vttype AstValTypeType
	Item   AstNode
}

func (a *AstValType) ANType() AstType {
	return ANValType
}

func (a *AstValType) Debug(x uint) {
	ident(x)
	fmt.Printf("AstType\n")
	a.Item.Debug(x + 1)
}

func (a *AstValType) Where() reporter.Where {
	return a.Item.Where()
}

func (a *AstValType) IsLlType() bool {
	switch a.Vttype {
	case ANTBinary, ANTTuple:
		return true
	default:
		return false
	}
}

type AstValTypeBinary struct {
	Tok_type *Token
}

func (a *AstValTypeBinary) ANType() AstType {
	return ANTypeBinary
}

func (a *AstValTypeBinary) Debug(x uint) {
	ident(x)
	fmt.Printf("AstTypeBinary\n")
	debug_token(x + 1, a.Tok_type)
}

func (a *AstValTypeBinary) Where() reporter.Where {
	return a.Tok_type.where
}

type AstValTypeVar struct {
	Ident *Token
	Real *AstValType
}

func (a *AstValTypeVar) ANType() AstType {
	return ANTypeVar
}

func (a *AstValTypeVar) Debug(x uint) {
	ident(x)
	fmt.Printf("AstTypeVar\n")
	debug_token(x + 1, a.Ident)
	if a.Real != nil {
		a.Real.Debug(x + 1)
	}
}

func (a *AstValTypeVar) Where() reporter.Where {
	if a.Ident != nil {
		return a.Ident.Where()
	}
	return reporter.FakeWhere()
}
/*
type AstValTypeStruct struct {
	Vttype AstValTypeType
	Item   AstNode
}

func (a *AstValTypeStruct) ANType() AstType {
	return ANTypeStruct
}

func (a *AstValTypeStruct) Debug(x uint) {
	ident(x)
	fmt.Printf("AstTypeStruct\n")
	a.Item.Debug(x + 1)
}

func (a *AstValTypeStruct) Where() reporter.Where {
	return a.Item.Where()
}

type AstValTypeEnum struct {
	Vttype AstValTypeType
	Item   AstNode
}

func (a *AstValTypeEnum) ANType() AstType {
	return ANTypeEnum
}

func (a *AstValTypeEnum) Debug(x uint) {
	ident(x)
	fmt.Printf("AstTypeEnum\n")
	a.Item.Debug(x + 1)
}

func (a *AstValTypeEnum) Where() reporter.Where {
	return a.Item.Where()
}
*/

type AstValTypeFnType struct {
	Tok_fn      *Token
	Tok_lbrace  *Token
	Types       []AstValType
	Tok_rbrace  *Token
	Tok_thinarr *Token
	Ret_t       AstValType
}

func (a *AstValTypeFnType) ANType() AstType {
	return ANTypeFnType
}

func (a *AstValTypeFnType) Debug(x uint) {
	ident(x)
	fmt.Printf("AstTypeFnType\n")
	debug_token(x + 1, a.Tok_fn)
	debug_token(x + 1, a.Tok_lbrace)
	ident(x + 1)
	fmt.Printf("Arguments\n")
	for _, item := range a.Types {
		item.Debug(x + 2)
	}
	debug_token(x + 1, a.Tok_rbrace)
	debug_token(x + 1, a.Tok_thinarr)
	a.Ret_t.Debug(x + 1)
}

func (a *AstValTypeFnType) Where() reporter.Where {
	ty := a.Ret_t.Where()
	return a.Tok_fn.where.Merge(&ty)
}

type AstValTypeTuple struct {
	Tok_lbrace  *Token
	Types       []AstValType
	Tok_rbrace  *Token
}

func AstTupleNew(types ...AstValType) AstValType {
	return AstValType{
		Vttype: ANTTuple,
		Item: &AstValTypeTuple{
			Types: types,
		},
	}
}

func (a *AstValTypeTuple) ANType() AstType {
	return ANTypeTuple
}

func (a *AstValTypeTuple) Debug(x uint) {
	ident(x)
	fmt.Printf("AstTypeTuple\n")
	debug_token(x + 1, a.Tok_lbrace)
	ident(x + 1)
	fmt.Printf("Types\n")
	for _, item := range a.Types {
		item.Debug(x + 2)
	}
	debug_token(x + 1, a.Tok_rbrace)
}

func (a *AstValTypeTuple) Where() reporter.Where {
	return a.Tok_lbrace.where.Merge(&a.Tok_rbrace.where)
}

type AstExpr struct {
	Etype AstExprType
	Item  AstNode
}

func (a *AstExpr) ANType() AstType {
	return ANExpr
}

func (a *AstExpr) Debug(x uint) {
	ident(x)
	fmt.Printf("AstExpr\n")
	a.Item.Debug(x + 1)
}

func (a *AstExpr) Where() reporter.Where {
	return a.Item.Where()
}

type AstExprOp1 struct {
	Tok_op *Token
	Expr AstExpr
}

func (a *AstExprOp1) ANType() AstType {
	return ANExprOp1
}

func (a *AstExprOp1) Debug(x uint) {
	ident(x)
	fmt.Printf("AstExprOperator 1\n")
	debug_token(x + 1, a.Tok_op)
	a.Expr.Debug(x + 1)
}

func (a *AstExprOp1) Where() reporter.Where {
	expr := a.Expr.Where()
	return a.Tok_op.where.Merge(&expr)
}

type AstExprGroup struct {
	Tok_lbrace *Token
	Expr AstExpr
	Tok_rbrace *Token
}

func (a *AstExprGroup) ANType() AstType {
	return ANExprGroup
}

func (a *AstExprGroup) Debug(x uint) {
	ident(x)
	fmt.Printf("AstExprGroup\n")
	debug_token(x + 1, a.Tok_lbrace)
	a.Expr.Debug(x + 1)
	debug_token(x + 1, a.Tok_rbrace)
}

func (a *AstExprGroup) Where() reporter.Where {
	return a.Tok_lbrace.where.Merge(&a.Tok_rbrace.where)
}

type AstExprLiteral struct {
	Value *Token
}

func (a *AstExprLiteral) ANType() AstType {
	return ANExprLiteral
}

func (a *AstExprLiteral) Debug(x uint) {
	ident(x)
	fmt.Printf("AstExprLiteral\n")
	debug_token(x + 1, a.Value)
}

func (a *AstExprLiteral) Where() reporter.Where {
	return a.Value.Where()
}

type AstExprVar struct {
	Ident *Token
}

func (a *AstExprVar) ANType() AstType {
	return ANExprVar
}

func (a *AstExprVar) Debug(x uint) {
	ident(x)
	fmt.Printf("AstExprVariable\n")
	debug_token(x + 1, a.Ident)
}

func (a *AstExprVar) Where() reporter.Where {
	return a.Ident.Where()
}

type AstExprFncall struct {
	Callable AstExpr
	Tok_lbrace *Token
	Args []AstExpr
	Tok_rbrace *Token
}

func (a *AstExprFncall) ANType() AstType {
	return ANExprFncall
}

func (a *AstExprFncall) Debug(x uint) {
	ident(x)
	fmt.Printf("ANExprFncall\n")
	a.Callable.Debug(x + 1)
	debug_token(x + 1, a.Tok_lbrace)
	ident(x + 1)
	fmt.Printf("Arguments\n")
	for _, arg := range a.Args {
		arg.Debug(x + 2)
	}
	debug_token(x + 1, a.Tok_rbrace)
}

func (a *AstExprFncall) Where() reporter.Where {
	callable := a.Callable.Where()
	return a.Tok_rbrace.where.Merge(&callable)
}

type AstExprBuiltinCorePrint struct {
	Tok_kcoreprint *Token
	Expr AstExpr
}

func (a *AstExprBuiltinCorePrint) ANType() AstType {
	return ANExprBuiltinCorePrint
}

func (a *AstExprBuiltinCorePrint) Debug(x uint) {
	ident(x)
	fmt.Printf("AstExprBuiltinCorePrint\n")
	debug_token(x + 1, a.Tok_kcoreprint)
	a.Expr.Debug(x + 1)
}

func (a *AstExprBuiltinCorePrint) Where() reporter.Where {
	expr := a.Expr.Where()
	return a.Tok_kcoreprint.where.Merge(&expr)
}

type AstExprOp2 struct {
	Lhs AstExpr
	Op *Token
	Rhs AstExpr
}

func (a *AstExprOp2) ANType() AstType {
	return ANExprOp2
}

func (a *AstExprOp2) Debug(x uint) {
	ident(x)
	fmt.Printf("AstExprOperator 2\n")
	a.Lhs.Debug(x + 1)
	debug_token(x + 1, a.Op)
	a.Rhs.Debug(x + 1)
}

func (a *AstExprOp2) Where() reporter.Where {
	lhs := a.Lhs.Where()
	rhs := a.Rhs.Where()
	return lhs.Merge(&rhs)
}

type AstStmt struct {
	Stype AstStmtType
	Item  AstNode
}

func (a *AstStmt) ANType() AstType {
	return ANStmt
}

func (a *AstStmt) Debug(x uint) {
	ident(x)
	fmt.Printf("AstStmt\n")
	a.Item.Debug(x + 1)
}

func (a *AstStmt) Where() reporter.Where {
	return a.Item.Where()
}

type AstBlock struct {
	Tok_lbra *Token
	Items    []AstItem
	Tok_rbra *Token
}

func (a *AstBlock) ANType() AstType {
	return ANBlock
}

func (a *AstBlock) Debug(x uint) {
	ident(x)
	fmt.Printf("AstBlock\n")
	debug_token(x+1, a.Tok_lbra)
	ident(x + 1)
	fmt.Printf("Items\n")
	for _, item := range a.Items {
		item.Debug(x + 2)
	}
	debug_token(x+1, a.Tok_rbra)
}

func (a *AstBlock) Where() reporter.Where {
	return a.Tok_lbra.where.Merge(&a.Tok_rbra.where)
}

type AstLetStmt struct {
	Tok_let *Token
	Name *Token
	Tok_colon *Token
	Vtype AstValType
	Tok_assign *Token
	Expr *AstExpr
}

func (a *AstLetStmt) ANType() AstType {
	return ANStmtLet
}

func (a *AstLetStmt) Debug(x uint) {
	ident(x)
	fmt.Printf("AstLetStmt\n")
	debug_token(x + 1, a.Tok_let)
	debug_token(x + 1, a.Name)
	debug_token(x + 1, a.Tok_colon)
	a.Vtype.Debug(x + 1)
	debug_token(x + 1, a.Tok_assign)
	if a.Expr != nil {
		a.Expr.Debug(x + 1)
	}
}

func (a *AstLetStmt) Where() reporter.Where {
	if a.Tok_assign != nil {
		expr := a.Expr.Where()
		return a.Tok_let.where.Merge(&expr)
	}
	if a.Tok_colon != nil {
		vtype := a.Vtype.Where()
		return a.Tok_let.where.Merge(&vtype)
	}
	return a.Tok_let.where.Merge(&a.Name.where)
}

type AstReturnStmt struct {
	Tok_return *Token
	Expr AstExpr
}

func (a *AstReturnStmt) ANType() AstType {
	return ANStmtReturn
}

func (a *AstReturnStmt) Debug(x uint) {
	ident(x)
	fmt.Printf("AstReturnStmt\n")
	debug_token(x + 1, a.Tok_return)
	a.Expr.Debug(x + 1)
}

func (a *AstReturnStmt) Where() reporter.Where {
	expr := a.Expr.Where()
	return a.Tok_return.where.Merge(&expr)
}

type AstBreakStmt struct {
	Tok_break *Token
}

func (a *AstBreakStmt) ANType() AstType {
	return ANStmtBreak
}

func (a *AstBreakStmt) Debug(x uint) {
	ident(x)
	fmt.Printf("AstBreakStmt\n")
	debug_token(x + 1, a.Tok_break)
}

func (a *AstBreakStmt) Where() reporter.Where {
	return a.Tok_break.where
}

type AstContinueStmt struct {
	tok_continue *Token
}

func (a *AstContinueStmt) ANType() AstType {
	return ANStmtContinue
}

func (a *AstContinueStmt) Debug(x uint) {
	ident(x)
	fmt.Printf("AstContinueStmt\n")
	debug_token(x + 1, a.tok_continue)
}

func (a *AstContinueStmt) Where() reporter.Where {
	return a.tok_continue.where
}

type AstUnkownType struct {
	For reporter.Where
	Rtype *AstValType
}

func (a *AstUnkownType) ANType() AstType {
	return ANStmtReturn
}

func (a *AstUnkownType) Debug(x uint) {
	ident(x)
	fmt.Printf("AstUnknownType\n")
	if a.Rtype != nil {
		a.Rtype.Debug(x + 1)
	}
}

func (a *AstUnkownType) Where() reporter.Where {
	return a.For
}

func ANTUnknownNew(where reporter.Where) AstValType {
	return AstValType{
		Vttype: ANTUnknown,
		Item: &AstUnkownType{
			For: where,
			Rtype: nil,
		},
	}
}