package frontend

import (
	"belt/reporter"
	"belt/utils"
	"fmt"
	"strconv"
	"strings"
)

type TokenType int

type Token struct {
	ttype TokenType
	value string
	where reporter.Where
}

const (
	EoF TokenType = iota

	Ident     // 🎀 ([_a-zA-Z][_a-zA-Z0-9]*)

	LlInt     // 🎀 ([0-9]+)
    LlFloat   // 🎀 ([0-9]+\.[0-9]*)
	LlString  // 🎀 "(.*)"
	LlBool    // 🎀 (true)|(false)
	LlNil     // 🎀 nil

	KFn       // 🎀 fn
	KIf       // 🎀 if
	KElse     // 🎀 else
	KWhile    // 🎀 while
	KFor      // 🎀 for
	KIn       // 🎀 in
	KLet      // 🎀 let
	KBreak    // 🎀 break
	KContinue // 🎀 continue
	KReturn // 🎀 return

	KTInt     // 🎀 int
	KTFloat   // 🎀 float
	KTString  // 🎀 string
	KTBool    // 🎀 bool
	KTVar     // 🎀 '<Ident>

	OAdd      // 🎀 +
	OAddf     // 🎀 +.
	OConnect  // 🎀 ++
	OSub      // 🎀 -
	OSubf     // 🎀 -.
	OMul      // 🎀 *
	OMulf     // 🎀 *.
	ODiv      // 🎀 /
	ODivf     // 🎀 /.
	OEq       // 🎀 ==
	ONeq	  // 🎀 !=
	OGrt      // 🎀 >
	OGeq      // 🎀 >=
	OLes      // 🎀 <
	OLeq      // 🎀 <=
	OAnd      // 🎀 &&
	OOr       // 🎀 ||
	OBXor     // 🎀 ^
	OBAnd     // 🎀 &
	OBOr      // 🎀 |
	ONot      // 🎀 !
	OBNot     // 🎀 ~
	OMovl     // 🎀 <<
	OMovr     // 🎀 >>
	OMember   // 🎀 .
	OLookup   // 🎀 ::
	OAssign   // 🎀 =

	Colon     // 🎀 :
	Comma     // 🎀 ,
	Semi      // 🎀 ;
	ThinArr   // 🎀 ->
	FatArr    // 🎀 =>

	LBrace    // 🎀 (
	RBrace    // 🎀 )
	LBracket  // 🎀 [
	RBracket  // 🎀 ]
	LBra      // 🎀 {
	RBra      // 🎀 }
)

var TokenStringMap = map[TokenType]string{
	EoF:      	"eof",
	Ident:    	"identifier",
	LlInt:    	"literal integer",
	LlFloat:  	"literal float",
	LlString: 	"literal string",
	LlBool:   	"literal boolean",
	LlNil:    	"literal nil",
	KFn:      	"keyword fn",
	KIf:      	"keyword if",
	KElse:    	"keyword else",
	KWhile:   	"keyword while",
	KFor:     	"keyword for",
	KIn:      	"keyword in",
	KLet:     	"keyword let",
	KBreak:   	"keyword break",
	KContinue: 	"keyword continue",
	KReturn: 	"keyword return",
	KTInt: 		"type `int`",
	KTFloat: 	"type `float`",
	KTString: 	"type `string`",
	KTBool: 	"type `bool`",
	KTVar: 		"type variable",
	OAdd: 		"operator +",
	OAddf: 		"operator +.",
	OConnect: 	"operator ++",
	OSub: 		"operator -",
	OSubf: 		"operator -.",
	OMul: 		"operator *",
	OMulf: 		"operator *.",
	ODiv: 		"operator /",
	ODivf: 		"operator /.",
	OEq: 		"operator ==",
	ONeq: 		"operator !=",
	OGrt: 		"operator >",
	OGeq: 		"operator >=",
	OLes: 		"operator <",
	OLeq: 		"operator <=",
	OAnd: 		"operator &&",
	OOr: 		"operator ||",
	OBXor: 		"operator ^",
	OBAnd: 		"operator &",
	OBOr: 		"operator |",
	ONot: 		"operator !",
	OBNot: 		"operator ~",
	OMovl: 		"operator <<",
	OMovr: 		"operator >>",
	OMember: 	"operator .",
	OLookup: 	"operator ::",
	OAssign: 	"operator =",
	Colon: 		"colon",
	Comma: 		"comma",
	Semi: 		"semi",
	ThinArr: 	"thin arrow",
	FatArr: 	"fat arrow",
	LBrace: 	"left brace",
	RBrace: 	"right brace",
	LBracket: 	"left square bracket",
	RBracket: 	"right square bracket",
	LBra: 		"left bracket",
	RBra: 		"right bracket",
}

func (tt *TokenType) ToString() string {
	return TokenStringMap[*tt]
}

type TokenCastError struct {
	token *Token
	message string
}

func (e TokenCastError) Error() string {
	return fmt.Sprintf("parse `%v` error: %v %v", e.token.value, e.message, e.token.where.ToString())
}

func (t *Token) AsInt() (int, error) {
	i, err := strconv.ParseInt(t.value, 10, 0)
	if err != nil {
		return 0, TokenCastError {
			token: t,
			message: err.Error(),
		}
	}
	return int(i), nil
}

func (t *Token) AsFloat() (float64, error) {
	i, err := strconv.ParseFloat(t.value, 64)
	if err != nil {
		return 0, TokenCastError {
			token: t,
			message: err.Error(),
		}
	}
	return float64(i), nil
}

func (t *Token) AsBool() (bool, error) {
	if t.value == "true" {
		return true, nil
	} else if t.value == "false" {
		return false, nil
	}
	return false, TokenCastError {
		token: t,
		message: fmt.Sprintf("cannot convert `%v` into boolean", t.value),
	}
}

func (t *Token) AssertType(tt TokenType) bool {
	return t.ttype == tt
}

func (t *Token) Where() reporter.Where {
	return t.where
}

type TokenStream struct {
	tokens []Token
	curr uint
}

func (ts *TokenStream) Next() *Token {
	res := &ts.tokens[ts.curr]
	ts.curr += 1
	return res
}

func (ts *TokenStream) Peek() *Token {
	res := &ts.tokens[ts.curr]
	return res
}

func (ts *TokenStream) Forward() {
	ts.curr += 1;
}

func (ts *TokenStream) Backward() {
	ts.curr -= 1;
}

func (ts *TokenStream) IsEof() bool {
	return ts.Peek().ttype == EoF
}

func (ts *TokenStream) AssertNext(tt TokenType) *Token {
	tok := ts.Peek()
	if tok.ttype == tt {
		ts.Forward()
		return tok
	}
	return nil
}

func (ts *TokenStream) AssertNextOrReport(tt TokenType, file *utils.File) *Token {
	tok := ts.Peek()
	if tok.ttype == tt {
		ts.Forward()
		return tok
	}
	err := reporter.Error(
		tok.where,
		fmt.Sprintf("expected %v, found %v", tt.ToString(), tok.ttype.ToString()),
	)
	reporter.Report(&err, file)
	utils.Exit(1)
	panic("reaching an unreachable code! something went wrong") // unreachable
}

func (ts *TokenStream) ToString() string {
	res := make([]string, 0)
	for i := range(ts.tokens) {
		tok := ts.tokens[i]
		switch tok.ttype {
		case LlString:
			res = append(res, fmt.Sprintf("\"%v\"", tok.value))
		default:
			res = append(res, tok.value)
		}
	}
	return strings.Join(res, " ")
}