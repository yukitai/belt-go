package compiler

import (
	"belt/frontend"
	"belt/utils"
	"fmt"
)

/*
(
	1. Lexer
	2. Parser
	3. Analyzer
	4. write LLVM-IR to `{name}.ll`
	5. run `llc -filetype=obj {name}.ll`
) -- CompileFile
(
	6. run `clang {name}.o -o {name}`
	7. remove `{name}.ll` & `{name}.o`
) -- Compile
*/

func CompileFile(path string) {
	file := utils.FileOpen(path)
	lexer := frontend.LexerFromFile(&file)
	tokens := lexer.Tokenize()
	// fmt.Printf("TokenStream {\n\t%v\n}\n", tokens.ToString())
	parser := frontend.ParserNew(&file, tokens)
	ast := parser.ParseFile()
	// ast.Debug(0)
	analyzer := frontend.AnalyzerNew(&ast, &file)
	analyzer.Analyze()
	builder := frontend.AstLLVMBuilderNew(ast, &file)
	ir := builder.Build()
	fmt.Printf("%v\n", ir)
}