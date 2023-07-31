package frontend

import "belt/utils"

type Analyzer struct {
	ast     *AstFile
	globals map[string]*AstType
	file    *utils.File
	has_err bool
}

func AnalyzerNew(ast *AstFile, file *utils.File) Analyzer {
	return Analyzer{
		ast:     ast,
		globals: make(map[string]*AstType),
		file:    file,
		has_err: false,
	}
}

func (a *Analyzer) Analyze() {
	tyinfer := TyInferNew(a)
	tyinfer.InitGlobal()
	tyinfer.InferAll()

	if a.has_err {
		utils.Exit(1)
	}
}