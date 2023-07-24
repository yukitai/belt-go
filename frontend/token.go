package frontend

import (
	"belt/reporter"
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

func (ts *TokenStream) Next() Token {
	res := ts.tokens[ts.curr]
	ts.curr += 1
	return res
}

func (ts *TokenStream) Peek() Token {
	res := ts.tokens[ts.curr]
	return res
}

func (ts *TokenStream) Forward() {
	ts.curr += 1;
}

func (ts *TokenStream) Backward() {
	ts.curr -= 1;
}

func (ts *TokenStream) IsEoF() bool {
	return ts.Peek().ttype == EoF
}

func (ts *TokenStream) AssertNext(tt TokenType) bool {
	tok := ts.Peek()
	if tok.ttype == tt {
		ts.Forward()
		return true
	}
	return false
}

func (ts *TokenStream) ToString() string {
	res := make([]string, 0)
	for i := range(ts.tokens) {
		tok := ts.tokens[i]
		switch tok.ttype {
		case LlString:
			res = append(res, "\"", tok.value, "\"")
		default:
			res = append(res, tok.value)
		}
	}
	return strings.Join(res, " ")
}