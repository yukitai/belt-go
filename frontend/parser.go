package frontend

import (
	"belt/reporter"
	"belt/utils"
	"fmt"
)

type Parser struct {
	file *utils.File
	tokens TokenStream
}

func ParserNew(file *utils.File, tokens TokenStream) Parser {
	return Parser {
		file, tokens,
	}
}

func (p *Parser) ParseFile() AstFile {
	items := make([]AstItem, 0)
	for ; !p.tokens.IsEof(); {
		items = append(items, p.ParseItem())
	}
	return AstFile{
		Items: items,
	}
}

func (p *Parser) ParseItem() AstItem {
	return p.ParseStmt()
}

func (p *Parser) ParseStmt() AstStmt {
	tok := p.tokens.Peek()
	switch tok.ttype {
	case KLet:
		stmt := p.ParseLetStmt()
		return AstStmt{
			Stype: ANSExpr,
			Item: &stmt,
		}
	case KBreak:
		stmt := p.ParseBreakStmt()
		return AstStmt{
			Stype: ANSExpr,
			Item: &stmt,
		}
	case KContinue:
		stmt := p.ParseContinueStmt()
		return AstStmt{
			Stype: ANSExpr,
			Item: &stmt,
		}
	case KReturn:
		stmt := p.ParseReturnStmt()
		return AstStmt{
			Stype: ANSExpr,
			Item: &stmt,
		}
	default:
		expr := p.ParseExpr()
		return AstStmt{
			Stype: ANSExpr,
			Item: &expr,
		}
	}
}

func (p *Parser) ParseLetStmt() AstLetStmt {
	tok_let := p.tokens.Next()
	name := p.tokens.AssertNextOrReport(Ident, p.file)
	tok_colon := p.tokens.AssertNext(Colon)
	var vtype AstValType
	if tok_colon != nil {
		vtype = p.ParseValType()
	} else {
		vtype = ANTUnknownNew()
	}
	tok_assign := p.tokens.AssertNext(OAssign)
	var expr *AstExpr
	if tok_assign != nil {
		e := p.ParseExpr()
		expr = &e
	}
	return AstLetStmt{
		Tok_let: &tok_let,
		Name: &name,
		Tok_colon: tok_colon,
		Vtype: vtype,
		Tok_assign: tok_assign,
		Expr: expr,
	}
}

func (p *Parser) ParseBreakStmt() AstBreakStmt {
	tok_break := p.tokens.Next()
	return AstBreakStmt{
		Tok_break: &tok_break,
	}
}

func (p *Parser) ParseContinueStmt() AstContinueStmt {
	tok_continue := p.tokens.Next()
	return AstContinueStmt{
		tok_continue: &tok_continue,
	}
}

func (p *Parser) ParseReturnStmt() AstReturnStmt {
	tok_return := p.tokens.Next()
	expr := p.ParseExpr()
	return AstReturnStmt{
		Tok_return: &tok_return,
		Expr: expr,
	}
}

func (p *Parser) ParseExpr() AstExpr {
	return p.ParseExprBinary()
}

func (p *Parser) ParseExprBinary() AstExpr {
	tok := p.tokens.Next()
	switch tok.ttype {
	case LlInt, LlFloat, LlString, LlBool, LlNil:
		return AstExpr{
			Etype: ANELiteral,
			Item: &AstExprLiteral{
				Value: &tok,
			},
		}
	case Ident:
		return AstExpr{
			Etype: ANEVar,
			Item: &AstExprVar{
				Ident: &tok,
			},
		}
	default:
		err := reporter.Error(
			tok.where,
			fmt.Sprintf("unexpected %v", tok.ttype.ToString()),
		)
		reporter.Report(&err, p.file)
		utils.Exit(1)
		panic("reaching an unreachable code! something went wrong")
	}
}

func (p *Parser) ParseValType() AstValType {
	panic("not implemented yet")
}