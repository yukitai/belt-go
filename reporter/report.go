package reporter

import "belt/utils"

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
	Report(*utils.File)
}

func Report(reporter Reporter, f *utils.File) {
	level := reporter.Level()
	if level > RecieveReportLevel {
		return
	}
	reporter.Report(f)
}