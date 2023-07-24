package frontend

import (
	"fmt"
	"strings"
)

type AstType int

const (
	ANFile AstType = iota
	ANFndecl
	ANFnArg
	ANExpr
	ANStmt
	ANBuiltinName
	ANBuiltinStmt
	ANValType
	ANBlock
)

type AstExprType int

const (
	ANEOp2 AstExprType = iota
	ANEOp1
	ANEGroup
)

type AstValTypeType int

const (
	ANTBinary AstValTypeType = iota
	ANTVar
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

func debug_token(x uint, tok *Token) {
	if tok != nil {
		ident(x)
		fmt.Printf("Token { %v, %v }", tok.ttype, tok.value)
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
		Item := a.Items[i]
		Item.Debug(x + 1)
	}
}

type AstItem struct {
	Antype AstType
	Item   AstNode
}

func (a *AstItem) ANType() AstType {
	return a.Antype
}

func (a *AstItem) Debug(x uint) {
	ident(x)
	fmt.Printf("AstItem\n")
	a.Item.Debug(x + 1)
}

type AstFnDecl struct {
	Tok_fn     *Token
	Name       *Token
	Tok_lbrace *Token
	Args       []AstFnArg
	Tok_rbrace *Token
	Ret_t      *AstValType
	Body       AstBlock
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
	Item AstNode
}

func (a *AstValType) ANType() AstType {
	return ANValType
}

func (a *AstValType) Debug(x uint) {
	ident(x)
	fmt.Printf("AstType\n")
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
		Item := a.Items[i]
		Item.Debug(x + 1)
	}
	debug_token(x+1, a.Tok_rbra)
}
