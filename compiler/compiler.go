package compiler

import (
	"belt/frontend"
	"belt/utils"
)

func CompileFile(path string) {
	file := utils.FileOpen(path)
	lexer := frontend.LexerFromFile(&file)
	tokens := lexer.Tokenize()
	// fmt.Printf("TokenStream {\n\t%v\n}\n", tokens.ToString())
	parser := frontend.ParserNew(&file, tokens)
	ast := parser.ParseFile()
	ast.Debug(0)
}