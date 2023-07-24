package utils

import (
	"fmt"
	"os"
)

type File struct {
	src []rune
	name string
}

func FileFromString(src string, name string) File {
	return File {
		src: ([]rune)(src),
		name: name,
	}
}

func FileOpen(path string) File {
	ofileinfo, err := os.Stat(path)
	if err != nil {
		CompilerError(err.Error())
		Exit(1)
	}
	if ofileinfo.IsDir() {
		CompilerError(fmt.Sprintf("`%v` is a directory", path))
		Exit(1)
	}
	octt, _ := os.ReadFile(path)
	return FileFromString(string(octt), path)
}

func (f *File) Src() []rune {
	return f.src
}

func (f *File) Name() string {
	return f.name
}

func (f *File) At(i uint) rune {
	return f.src[i]
}

func (f *File) Slice(i, j uint) []rune {
	return f.src[i:j]
}