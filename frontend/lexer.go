package frontend

import (
	"belt/compiler"
	"belt/reporter"
	"fmt"
)

type Lexer struct {
	src   *compiler.File
	line  uint
	curr  uint
	cmark uint
	lmark uint
}

func (l *Lexer) next() rune {
	res := l.src.At(l.curr)
	l.curr += 1
	return res
}

func (l *Lexer) peek() rune {
	res := l.src.At(l.curr)
	return res
}

func (l *Lexer) forward() {
	l.curr += 1
}

func (l *Lexer) backward() {
	l.curr -= 1
}

func (l *Lexer) begin() {
	l.cmark = l.curr - 1
	l.lmark = l.line
}

func (l *Lexer) end() reporter.Where {
	return reporter.WhereNew(l.lmark, l.line, l.cmark, l.curr)
}

func (l *Lexer) here() reporter.Where {
	return reporter.WhereNew(l.line, l.line, l.curr-1, l.curr)
}

func (l *Lexer) is_eof() bool {
	return l.curr >= uint(len(l.src.Src()))
}

func LexerFromFile(src *compiler.File) Lexer {
	return Lexer{
		src: src,
	}
}

func (l *Lexer) Tokenize() TokenStream {
	tokens := make([]Token, 0)
	for !l.is_eof() {
		tok := l.next()
		switch {
		case tok == ' ' || tok == '\r':
			{
			}
		case tok == '\n':
			l.line += 1
		case tok >= '0' && tok <= '9':
			l.begin()
			var is_float bool
			var ttype TokenType
		loop_n:
			for !l.is_eof() {
				tok := l.peek()
				switch {
				case tok >= '0' && tok <= '9':
					l.forward()
				case tok == '.':
					if is_float {
						break loop_n
					} else {
						l.forward()
						is_float = true
					}
				default:
					break loop_n
				}
			}
			if is_float {
				ttype = LlFloat
			} else {
				ttype = LlInt
			}
			value := string(l.src.Slice(l.cmark, l.curr))
			tokens = append(tokens, Token{
				ttype: ttype,
				value: value,
				where: l.end(),
			})
		case tok == '_' || (tok >= 'A' && tok <= 'Z') || (tok >= 'a' && tok <= 'z'):
			l.begin()
			var ttype TokenType
		loop_i:
			for !l.is_eof() {
				tok := l.peek()
				switch {
				case tok == '_' || (tok >= 'A' && tok <= 'Z') || (tok >= 'a' && tok <= 'z') || (tok >= '0' && tok <= '9'):
					l.forward()
				default:
					break loop_i
				}
			}
			value := string(l.src.Slice(l.cmark, l.curr))
			switch value {
			case "fn":
				ttype = KFn
			case "let":
				ttype = KLet
			default:
				ttype = Ident
			}
			tokens = append(tokens, Token{
				ttype: ttype,
				value: value,
				where: l.end(),
			})
		default:
			err := reporter.Error(
				l.here(),
				fmt.Sprintf("unexpected char `%v`", string(tok)),
			)
			reporter.Report(&err, l.src)
			compiler.Exit(1)
		}
	}
	return TokenStream{
		tokens: tokens, curr: 0,
	}
}
