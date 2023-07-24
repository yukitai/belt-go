package frontend

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type AstType int

const (
	ANFile AstType = iota
	
	ANFndecl
	ANFnArg

	ANExpr
	ANExprOp1
	ANExprOp2
	ANExprGroup
	ANExprLiteral
	ANExprVar

	ANStmt
	ANBuiltinStmt
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

	ANBlock
)

type AstExprType int

const (
	ANEOp2 AstExprType = iota
	ANEOp1
	ANEGroup
	ANELiteral
	ANEVar
)

type AstValTypeType int

const (
	ANTBinary AstValTypeType = iota
	ANTVar
	ANTUnknown
	ANTStruct
	ANTEnum
)

type AstStmtType int

const (
	ANSExpr AstStmtType = iota
	ANSBuiltinStmt
	ANSLet
	ANSReturn
	ANSBreak
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
		fmt.Printf("Token { %v, %v }\n", tok.ttype, tok.value)
	} else {
		fmt.Printf("Token { %v }\n", debugNilString)
	}
}

type AstNode interface {
	ANType() AstType
	Debug(uint)
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
	for i := range a.Items {
		item := a.Items[i]
		item.Debug(x + 1)
	}
}

type AstItem = AstStmt

type AstFnDecl struct {
	Tok_fn      *Token
	Name        *Token
	Tok_lbrace  *Token
	Args        []AstFnArg
	Tok_rbrace  *Token
	Tok_thinarr *Token
	Ret_t       *AstValType
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
	for i := range a.Args {
		arg := a.Args[i]
		arg.Debug(x + 1)
	}
	debug_token(x+1, a.Tok_rbrace)
	a.Ret_t.Debug(x + 1)
	a.Body.Debug(x + 1)
}

type AstFnArg struct {
	Name      *Token
	Tok_colon *Token
	Atype     *AstValType
}

func (a *AstFnArg) ANType() AstType {
	return ANFnArg
}

func (a *AstFnArg) Debug(x uint) {
	ident(x)
	fmt.Printf("AstFnArg\n")
	debug_token(x+1, a.Name)
	debug_token(x+1, a.Tok_colon)
	if a.Atype != nil {
		a.Atype.Debug(x + 1)
	}
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

type AstExprOp2 struct {
	Lhs AstExpr
	Tok_op *Token
	Rhs AstExpr
}

func (a *AstExprOp2) ANType() AstType {
	return ANExprOp2
}

func (a *AstExprOp2) Debug(x uint) {
	ident(x)
	fmt.Printf("AstExprOperator 2\n")
	a.Lhs.Debug(x + 1)
	debug_token(x + 1, a.Tok_op)
	a.Rhs.Debug(x + 1)
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
	for i := range a.Items {
		item := a.Items[i]
		item.Debug(x + 1)
	}
	debug_token(x+1, a.Tok_rbra)
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

type AstUnkownType struct {
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

func ANTUnknownNew() AstValType {
	return AstValType{
		Vttype: ANTUnknown,
		Item: &AstUnkownType{
			Rtype: nil,
		},
	}
}