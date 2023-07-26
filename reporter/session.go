package reporter

import (
	"belt/utils"
	"fmt"
	"strings"
	"github.com/fatih/color"
)

type Session struct {
	where Where
	file *utils.File
}

func SessionNew(where Where, file *utils.File) Session {
	return Session {
		where, file,
	}
}

func (s *Session) GetLines() []string {
	return strings.Split(string(s.Mark(s.file.Src())), "\n")
}

var color_leader = color.New(color.FgHiBlack)
var color_mark = color.New(color.Underline, color.Bold)

func (s *Session) Mark(src []rune) []rune {
	marked := make([]rune, 0)
	for i, chr := range(src) {
		i := uint(i)
		if i >= s.where.begin && i < s.where.end && chr != '\n' {
			marked = append(
				marked,
				([]rune)(color_mark.Sprintf("%v", string(chr)))...,
			)
		} else {
			marked = append(marked, src[i])
		}
	}
	return marked
}

func (s Session) Print() {
	lines := s.GetLines()
	start := max(1, s.where.line1 - 1)
	end := min(uint(len(lines)), s.where.line2 + 1)
	for line := start; line <= end; line += 1 {
		color_leader.Printf("%7v ", line)
		fmt.Printf("%v\n", lines[line - 1])
	}
}

func (s Session) Println() {
	s.Print()
	fmt.Printf("\n")
}