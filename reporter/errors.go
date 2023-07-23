package reporter

import (
	"belt/compiler"
	"fmt"

	"github.com/fatih/color"
)

type ErrorReporter struct {
	where Where
	message string
}

func Error(where Where, message string) ErrorReporter {
	return ErrorReporter {
		where, message,
	}
}

func (r *ErrorReporter) Level() ReportLevel {
	return LError
}

var color_error = color.New(color.FgHiRed, color.Bold)
var color_where = color.New(color.FgHiCyan)

func (r *ErrorReporter) Report(f *compiler.File) {
	color_error.Printf("error")
	fmt.Printf(": %v\n", r.message)
	color_where.Printf("   --> in %v, %v", f.Name(), r.where.ToString())
	fmt.Printf("\n")
	SessionNew(r.where, f).Println()
}