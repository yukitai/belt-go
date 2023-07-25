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
	case KFn:
		stmt := p.ParseFndeclStmt()
		return AstStmt{
			Stype: ANSFndecl,
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

func (p *Parser) ParseFndeclStmt() AstFnDecl {
	tok_fn := p.tokens.Next()
	name := p.tokens.AssertNext(Ident)
	tok_lbrace := p.tokens.AssertNextOrReport(LBrace, p.file)
	args := make([]AstFnArg, 0)
	var tok_rbrace *Token
	for ; !p.tokens.IsEof(); {
		tok := p.tokens.Peek()
		if tok.ttype == RBrace {
			tok_rbrace = tok
			p.tokens.Forward()
			break
		}
		aname := p.tokens.AssertNextOrReport(Ident, p.file)
		tok_colon := p.tokens.AssertNext(Colon)
		var atype AstValType
		if tok_colon != nil {
			atype = p.ParseValType()
		} else {
			atype = ANTUnknownNew()
		}
		tok_comma := p.tokens.AssertNext(Comma)
		args = append(args, AstFnArg{
			Name: aname,
			Tok_colon: tok_colon,
			Atype: atype,
			Tok_comma: tok_comma,
		})
	}
	tok_thinarr := p.tokens.AssertNext(ThinArr)
	var ret_t AstValType
	if tok_thinarr != nil {
		ret_t = p.ParseValType()
	} else {
		ret_t = ANTUnknownNew()
	}
	body := p.ParseBlock()
	return AstFnDecl{
		Tok_fn: tok_fn,
		Name: name,
		Tok_lbrace: tok_lbrace,
		Args: args,
		Tok_rbrace: tok_rbrace,
		Tok_thinarr: tok_thinarr,
		Ret_t: ret_t,
		Body: body,
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
		Tok_let: tok_let,
		Name: name,
		Tok_colon: tok_colon,
		Vtype: vtype,
		Tok_assign: tok_assign,
		Expr: expr,
	}
}

func (p *Parser) ParseBreakStmt() AstBreakStmt {
	tok_break := p.tokens.Next()
	return AstBreakStmt{
		Tok_break: tok_break,
	}
}

func (p *Parser) ParseContinueStmt() AstContinueStmt {
	tok_continue := p.tokens.Next()
	return AstContinueStmt{
		tok_continue: tok_continue,
	}
}

func (p *Parser) ParseReturnStmt() AstReturnStmt {
	tok_return := p.tokens.Next()
	expr := p.ParseExpr()
	return AstReturnStmt{
		Tok_return: tok_return,
		Expr: expr,
	}
}

func (p *Parser) ParseExpr() AstExpr {
	return p.ParseExprOr()
}

func (p *Parser) ParseExprOr() AstExpr {
	lhs := p.ParseExprAnd()
loop:
	for ; !p.tokens.IsEof(); {
		tok := p.tokens.Peek()
		switch tok.ttype {
		case OOr:
			p.tokens.Forward()
			rhs := p.ParseExprAnd()
			lhs = AstExpr{
				Etype: ANEOp2,
				Item: &AstExprOp2{
					Lhs: lhs,
					Op: tok,
					Rhs: rhs,
				},
			}
		default:
			break loop
		}
	}
	return lhs
}

func (p *Parser) ParseExprAnd() AstExpr {
	lhs := p.ParseExprBOr()
loop:
	for ; !p.tokens.IsEof(); {
		tok := p.tokens.Peek()
		switch tok.ttype {
		case OAnd:
			p.tokens.Forward()
			rhs := p.ParseExprBOr()
			lhs = AstExpr{
				Etype: ANEOp2,
				Item: &AstExprOp2{
					Lhs: lhs,
					Op: tok,
					Rhs: rhs,
				},
			}
		default:
			break loop
		}
	}
	return lhs
}

func (p *Parser) ParseExprBOr() AstExpr {
	lhs := p.ParseExprBXor()
loop:
	for ; !p.tokens.IsEof(); {
		tok := p.tokens.Peek()
		switch tok.ttype {
		case OBOr:
			p.tokens.Forward()
			rhs := p.ParseExprBXor()
			lhs = AstExpr{
				Etype: ANEOp2,
				Item: &AstExprOp2{
					Lhs: lhs,
					Op: tok,
					Rhs: rhs,
				},
			}
		default:
			break loop
		}
	}
	return lhs
}

func (p *Parser) ParseExprBXor() AstExpr {
	lhs := p.ParseExprBAnd()
loop:
	for ; !p.tokens.IsEof(); {
		tok := p.tokens.Peek()
		switch tok.ttype {
		case OBXor:
			p.tokens.Forward()
			rhs := p.ParseExprBAnd()
			lhs = AstExpr{
				Etype: ANEOp2,
				Item: &AstExprOp2{
					Lhs: lhs,
					Op: tok,
					Rhs: rhs,
				},
			}
		default:
			break loop
		}
	}
	return lhs
}

func (p *Parser) ParseExprBAnd() AstExpr {
	lhs := p.ParseExprEqual()
loop:
	for ; !p.tokens.IsEof(); {
		tok := p.tokens.Peek()
		switch tok.ttype {
		case OBAnd:
			p.tokens.Forward()
			rhs := p.ParseExprEqual()
			lhs = AstExpr{
				Etype: ANEOp2,
				Item: &AstExprOp2{
					Lhs: lhs,
					Op: tok,
					Rhs: rhs,
				},
			}
		default:
			break loop
		}
	}
	return lhs
}

func (p *Parser) ParseExprEqual() AstExpr {
	lhs := p.ParseExprCompare()
loop:
	for ; !p.tokens.IsEof(); {
		tok := p.tokens.Peek()
		switch tok.ttype {
		case OEq, ONeq:
			p.tokens.Forward()
			rhs := p.ParseExprCompare()
			lhs = AstExpr{
				Etype: ANEOp2,
				Item: &AstExprOp2{
					Lhs: lhs,
					Op: tok,
					Rhs: rhs,
				},
			}
		default:
			break loop
		}
	}
	return lhs
}

func (p *Parser) ParseExprCompare() AstExpr {
	lhs := p.ParseExprMov()
loop:
	for ; !p.tokens.IsEof(); {
		tok := p.tokens.Peek()
		switch tok.ttype {
		case OLes, OLeq, OGrt, OGeq:
			p.tokens.Forward()
			rhs := p.ParseExprMov()
			lhs = AstExpr{
				Etype: ANEOp2,
				Item: &AstExprOp2{
					Lhs: lhs,
					Op: tok,
					Rhs: rhs,
				},
			}
		default:
			break loop
		}
	}
	return lhs
}

func (p *Parser) ParseExprMov() AstExpr {
	lhs := p.ParseExprAdd()
loop:
	for ; !p.tokens.IsEof(); {
		tok := p.tokens.Peek()
		switch tok.ttype {
		case OMovl, OMovr:
			p.tokens.Forward()
			rhs := p.ParseExprAdd()
			lhs = AstExpr{
				Etype: ANEOp2,
				Item: &AstExprOp2{
					Lhs: lhs,
					Op: tok,
					Rhs: rhs,
				},
			}
		default:
			break loop
		}
	}
	return lhs
}

func (p *Parser) ParseExprAdd() AstExpr {
	lhs := p.ParseExprMul()
loop:
	for ; !p.tokens.IsEof(); {
		tok := p.tokens.Peek()
		switch tok.ttype {
		case OAdd, OAddf, OSub, OSubf, OConnect:
			p.tokens.Forward()
			rhs := p.ParseExprMul()
			lhs = AstExpr{
				Etype: ANEOp2,
				Item: &AstExprOp2{
					Lhs: lhs,
					Op: tok,
					Rhs: rhs,
				},
			}
		default:
			break loop
		}
	}
	return lhs
}

func (p *Parser) ParseExprMul() AstExpr {
	lhs := p.ParseExprUnary()
loop:
	for ; !p.tokens.IsEof(); {
		tok := p.tokens.Peek()
		switch tok.ttype {
		case OMul, OMulf, ODiv, ODivf:
			p.tokens.Forward()
			rhs := p.ParseExprUnary()
			lhs = AstExpr{
				Etype: ANEOp2,
				Item: &AstExprOp2{
					Lhs: lhs,
					Op: tok,
					Rhs: rhs,
				},
			}
		default:
			break loop
		}
	}
	return lhs
}

func (p *Parser) ParseExprUnary() AstExpr {
	return p.ParseExprBinary()
}

func (p *Parser) ParseExprBinary() AstExpr {
	tok := p.tokens.Next()
	switch tok.ttype {
	case LlInt, LlFloat, LlString, LlBool, LlNil:
		return AstExpr{
			Etype: ANELiteral,
			Item: &AstExprLiteral{
				Value: tok,
			},
		}
	case Ident:
		return AstExpr{
			Etype: ANEVar,
			Item: &AstExprVar{
				Ident: tok,
			},
		}
	case LBrace:
		tok_lbrace := tok
		expr := p.ParseExpr()
		tok_rbrace := p.tokens.AssertNextOrReport(RBrace, p.file)
		return AstExpr{
			Etype: ANEGroup,
			Item: &AstExprGroup{
				Tok_lbrace: tok_lbrace,
				Expr: expr,
				Tok_rbrace: tok_rbrace,
			},
		}
	case LBra:
		p.tokens.Backward()
		block := p.ParseBlock()
		return AstExpr{
			Etype: ANEBlock,
			Item: &block,
		}
	case OBOr:
		p.tokens.Backward()
		closure := p.ParseClosure()
		return AstExpr{
			Etype: ANEClosure,
			Item: &closure,
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

func (p *Parser) ParseBlock() AstBlock {
	tok_lbra := p.tokens.AssertNextOrReport(LBra, p.file)
	items := make([]AstStmt, 0)
	var tok_rbra *Token
	for ; !p.tokens.IsEof(); {
		tok := p.tokens.Peek()
		if tok.ttype == RBra {
			tok_rbra = tok
			p.tokens.Forward()
			break
		}
		stmt := p.ParseStmt()
		items = append(items, stmt)
	}
	return AstBlock{
		Tok_lbra: tok_lbra,
		Items: items,
		Tok_rbra: tok_rbra,
	}
}

func (p *Parser) ParseClosure() AstClosure {
	tok_lbor := p.tokens.Next()
	args := make([]AstFnArg, 0)
	var tok_rbor *Token
	for ; !p.tokens.IsEof(); {
		tok := p.tokens.Peek()
		if tok.ttype == OBOr {
			tok_rbor = tok
			p.tokens.Forward()
			break
		}
		aname := p.tokens.AssertNextOrReport(Ident, p.file)
		tok_colon := p.tokens.AssertNext(Colon)
		var atype AstValType
		if tok_colon != nil {
			atype = p.ParseValType()
		} else {
			atype = ANTUnknownNew()
		}
		tok_comma := p.tokens.AssertNext(Comma)
		args = append(args, AstFnArg{
			Name: aname,
			Tok_colon: tok_colon,
			Atype: atype,
			Tok_comma: tok_comma,
		})
	}
	tok_thinarr := p.tokens.AssertNext(ThinArr)
	var ret_t AstValType
	if tok_thinarr != nil {
		ret_t = p.ParseValType()
	} else {
		ret_t = ANTUnknownNew()
	}
	body := p.ParseExpr()
	return AstClosure{
		Tok_lbor: tok_lbor,
		Args: args,
		Tok_rbor: tok_rbor,
		Tok_thinarr: tok_thinarr,
		Ret_t: ret_t,
		Body: body,
	}
}