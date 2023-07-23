package frontend

import (
	"belt/reporter"
	"belt/compiler"
)

type Lexer struct {
	src compiler.File
	line uint
	curr uint
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
	l.cmark = l.curr
	l.lmark = l.line
}

func (l *Lexer) end() reporter.Where {
	return reporter.WhereNew(l.lmark, l.line, l.cmark, l.curr)
}

func (l *Lexer) here() reporter.Where {
	return reporter.WhereNew(l.line, l.line, l.curr - 1, l.curr)
}

func (l *Lexer) is_eof() bool {
	return l.curr >= uint(len(l.src.Src()))
}

func LexerFromFile(src compiler.File) Lexer {
	return Lexer {
		src: src,
	}
}

func (l *Lexer) Tokenize() TokenStream {
	tokens := make([]Token, 0)
	for ; !l.is_eof(); {
		tok := l.next()
		switch tok {
		case ' ', '\r': {}
		default:
			// todo: report here
		}
	}
	return TokenStream {
		tokens: tokens, curr: 0,
	}
}