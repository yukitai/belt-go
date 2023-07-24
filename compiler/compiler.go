package compiler

import (
	"belt/frontend"
	"belt/utils"
	"fmt"
)

func CompileFile(path string) {
	file := utils.FileOpen(path)
	lexer := frontend.LexerFromFile(&file)
	tokens := lexer.Tokenize()
	fmt.Printf("TokenStream {\n\t%v\n}\n", tokens.ToString())
}