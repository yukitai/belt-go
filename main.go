package main

import (
	"belt/compiler"
	"belt/frontend"
	"fmt"
)

func main() {
	file := compiler.FileFromString(
		`
let a = 114514
`, "test.bl")
	lexer := frontend.LexerFromFile(&file)
	tokens := lexer.Tokenize()
	fmt.Printf("%v", tokens)
}
