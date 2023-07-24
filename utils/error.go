package utils

import (
	"fmt"

	"github.com/fatih/color"
)

var color_error = color.New(color.FgHiRed, color.Bold)

func CompilerError(message string) {
	color_error.Printf("error")
	fmt.Printf(": %v\n\n", message)
}