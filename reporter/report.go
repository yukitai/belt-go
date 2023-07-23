package reporter

import "belt/compiler"

type ReportLevel int

const (
	LNever ReportLevel = iota
	LError
	LWarn
	LNote
	LHelp
)

var RecieveReportLevel = LHelp

type Reporter interface {
	Level() ReportLevel
	Report(*compiler.File)
}

func Report(reporter Reporter, f *compiler.File) {
	level := reporter.Level()
	if level > RecieveReportLevel {
		return
	}
	reporter.Report(f)
}