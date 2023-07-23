package frontend

type AstType int

type AstNode interface {
	ANType() AstType
}

const (
	ANFile AstType = iota
	ANFndecl
	ANExpr
	ANStmt
	ANBuiltinName
	ANBuiltinStmt
)

