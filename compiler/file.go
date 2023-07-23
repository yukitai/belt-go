package compiler

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