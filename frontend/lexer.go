package frontend

import (
	"belt/utils"
	"belt/reporter"
	"fmt"
)

type Lexer struct {
	src   *utils.File
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

func (l *Lexer) begin(t int) {
	l.cmark = uint(int(l.curr) + t)
	l.lmark = l.line
}

func (l *Lexer) end() reporter.Where {
	return reporter.WhereNew(l.lmark, l.line, l.cmark, l.curr)
}

func (l *Lexer) here() reporter.Where {
	return reporter.WhereNew(l.line, l.line, l.curr-1, l.curr)
}

func (l *Lexer) last(t uint) reporter.Where {
	return reporter.WhereNew(l.line, l.line, l.curr-t, l.curr)
}

func (l *Lexer) is_eof() bool {
	return l.curr >= uint(len(l.src.Src()))
}

func LexerFromFile(src *utils.File) Lexer {
	return Lexer{
		src: src,
		line: 1,
	}
}

func (l *Lexer) Tokenize() TokenStream {
	tokens := make([]Token, 0)
	for !l.is_eof() {
		tok := l.next()
		switch {
		case tok == ' ' || tok == '\r' || tok == '\t':
			{}
		case tok == '\n':
			l.line += 1
		case tok == '(':
			tokens = append(tokens, Token{
				ttype: LBrace,
				value: "(",
				where: l.here(),
			})
		case tok == ')':
			tokens = append(tokens, Token{
				ttype: RBrace,
				value: ")",
				where: l.here(),
			})
		case tok == '[':
			tokens = append(tokens, Token{
				ttype: LBracket,
				value: "[",
				where: l.here(),
			})
		case tok == ']':
			tokens = append(tokens, Token{
				ttype: RBracket,
				value: "]",
				where: l.here(),
			})
		case tok == '{':
			tokens = append(tokens, Token{
				ttype: LBra,
				value: "{",
				where: l.here(),
			})
		case tok == '}':
			tokens = append(tokens, Token{
				ttype: RBra,
				value: "}",
				where: l.here(),
			})
		case tok == '.':
			tokens = append(tokens, Token{
				ttype: OMember,
				value: ".",
				where: l.here(),
			})
		case tok == ',':
			tokens = append(tokens, Token{
				ttype: Comma,
				value: ",",
				where: l.here(),
			})
		case tok == ';':
			tokens = append(tokens, Token{
				ttype: Semi,
				value: ";",
				where: l.here(),
			})
		case tok == '^':
			tokens = append(tokens, Token{
				ttype: OBXor,
				value: "^",
				where: l.here(),
			})
		case tok == '~':
			tokens = append(tokens, Token{
				ttype: OBNot,
				value: "~",
				where: l.here(),
			})
		case tok == ':':
			ntok := l.peek()
			if ntok == ':' {
				l.forward()
				tokens = append(tokens, Token{
					ttype: OLookup,
					value: "::",
					where: l.last(2),
				})
			} else {
				tokens = append(tokens, Token{
					ttype: Colon,
					value: ":",
					where: l.here(),
				})
			}
		case tok == '+':
			ntok := l.peek()
			if ntok == '+' {
				l.forward()
				tokens = append(tokens, Token{
					ttype: OConnect,
					value: "++",
					where: l.last(2),
				})
			} else if ntok == '.' {
				l.forward()
				tokens = append(tokens, Token{
					ttype: OAddf,
					value: "+.",
					where: l.last(2),
				})
			} else {
				tokens = append(tokens, Token{
					ttype: OAdd,
					value: "+",
					where: l.here(),
				})
			}
		case tok == '-':
			ntok := l.peek()
			if ntok == '>' {
				l.forward()
				tokens = append(tokens, Token{
					ttype: ThinArr,
					value: "->",
					where: l.last(2),
				})
			} else if ntok == '.' {
				l.forward()
				tokens = append(tokens, Token{
					ttype: OSubf,
					value: "-.",
					where: l.last(2),
				})
			} else {
				tokens = append(tokens, Token{
					ttype: OSub,
					value: "-",
					where: l.here(),
				})
			}
		case tok == '*':
			ntok := l.peek()
			if ntok == '.' {
				l.forward()
				tokens = append(tokens, Token{
					ttype: OMulf,
					value: "*.",
					where: l.last(2),
				})
			} else {
				tokens = append(tokens, Token{
					ttype: OMul,
					value: "*",
					where: l.here(),
				})
			}
		case tok == '/':
			ntok := l.peek()
			if ntok == '.' {
				l.forward()
				tokens = append(tokens, Token{
					ttype: ODivf,
					value: "/.",
					where: l.last(2),
				})
			} else {
				tokens = append(tokens, Token{
					ttype: ODiv,
					value: "/",
					where: l.here(),
				})
			}
		case tok == '!':
			ntok := l.peek()
			if ntok == '=' {
				l.forward()
				tokens = append(tokens, Token{
					ttype: ONeq,
					value: "!=",
					where: l.last(2),
				})
			} else {
				tokens = append(tokens, Token{
					ttype: ONot,
					value: "!",
					where: l.here(),
				})
			}
		case tok == '&':
			ntok := l.peek()
			if ntok == '&' {
				l.forward()
				tokens = append(tokens, Token{
					ttype: OAnd,
					value: "&&",
					where: l.last(2),
				})
			} else {
				tokens = append(tokens, Token{
					ttype: OBAnd,
					value: "&",
					where: l.here(),
				})
			}
		case tok == '|':
			ntok := l.peek()
			if ntok == '|' {
				l.forward()
				tokens = append(tokens, Token{
					ttype: OOr,
					value: "||",
					where: l.last(2),
				})
			} else {
				tokens = append(tokens, Token{
					ttype: OBOr,
					value: "|",
					where: l.here(),
				})
			}
		case tok == '=':
			ntok := l.peek()
			if ntok == '>' {
				l.forward()
				tokens = append(tokens, Token{
					ttype: FatArr,
					value: "=>",
					where: l.last(2),
				})
			} else if ntok == '=' {
				l.forward()
				tokens = append(tokens, Token{
					ttype: OEq,
					value: "==",
					where: l.last(2),
				})
			} else {
				tokens = append(tokens, Token{
					ttype: OAssign,
					value: "=",
					where: l.here(),
				})
			}
		case tok == '>':
			ntok := l.peek()
			if ntok == '>' {
				l.forward()
				tokens = append(tokens, Token{
					ttype: OMovr,
					value: ">>",
					where: l.last(2),
				})
			} else if ntok == '=' {
				l.forward()
				tokens = append(tokens, Token{
					ttype: OGeq,
					value: ">=",
					where: l.last(2),
				})
			} else {
				tokens = append(tokens, Token{
					ttype: OGrt,
					value: ">",
					where: l.here(),
				})
			}
		case tok == '<':
			ntok := l.peek()
			if ntok == '<' {
				l.forward()
				tokens = append(tokens, Token{
					ttype: OMovl,
					value: "<<",
					where: l.last(2),
				})
			} else if ntok == '=' {
				l.forward()
				tokens = append(tokens, Token{
					ttype: OLeq,
					value: "<=",
					where: l.last(2),
				})
			} else {
				tokens = append(tokens, Token{
					ttype: OLes,
					value: "<",
					where: l.here(),
				})
			}
		case tok == '"':
			l.begin(-1)
			var nchr bool
			var value []rune
		loop_s:
			for !l.is_eof() {
				tok := l.peek()
				switch tok {
				case '"':
					l.forward()
					break loop_s
				case '\\':
					l.forward()
					if nchr {
						value = append(value, '\\')
					} else {
						nchr = true
					}
				default:
					l.forward()
					if nchr {
						switch tok {
						case '\'':
							value = append(value, '\'')
						case '"':
							value = append(value, '"')
						case 'n':
							value = append(value, '\n')
						case 'r':
							value = append(value, '\r')
						case 't':
							value = append(value, '\t')
						default:
							err := reporter.Error(
								l.last(2),
								fmt.Sprintf("unknown escape char `\\%v`", string(tok)),
							)
							reporter.Report(&err, l.src)
						}
					} else {
						value = append(value, tok)
					}
				}
			}
			tokens = append(tokens, Token{
				ttype: LlString,
				value: string(value),
				where: l.end(),
			})
		case tok >= '0' && tok <= '9':
			l.begin(-1)
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
		case tok == '\'':
			fallthrough
		case tok == '_' || (tok >= 'A' && tok <= 'Z') || (tok >= 'a' && tok <= 'z'):
			l.begin(-1)
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
			case "nil":
				ttype = LlNil
			case "true":
				ttype = LlBool
			case "false":
				ttype = LlBool
			case "fn":
				ttype = KFn
			case "if":
				ttype = KIf
			case "else":
				ttype = KElse
			case "while":
				ttype = KWhile
			case "for":
				ttype = KFor
			case "in":
				ttype = KIn
			case "let":
				ttype = KLet
			case "break":
				ttype = KBreak
			case "continue":
				ttype = KContinue
			case "int":
				ttype = KTInt
			case "float":
				ttype = KTFloat
			case "string":
				ttype = KTString
			case "bool":
				ttype = KTBool
			default:
				if value[0] == '\'' {
					ttype = KTVar
				} else {
					ttype = Ident
				}
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
			utils.Exit(1)
		}
	}
	return TokenStream{
		tokens: tokens, curr: 0,
	}
}
