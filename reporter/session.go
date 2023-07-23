package reporter

import (
	"belt/compiler"
	"fmt"
	"strings"
	"github.com/fatih/color"
)

type Session struct {
	where Where
	file *compiler.File
}

func SessionNew(where Where, file *compiler.File) Session {
	return Session {
		where, file,
	}
}

func (s *Session) GetLines() []string {
	return strings.Split(string(s.file.Src()), "\n")
}

func (s *Session) Mark(lines []string) []string {
	marked := make([]string, len(lines))
	for line := range(lines) {
		line2 := uint(line) + 1
		if line2 == s.where.line1 && line2 == s.where.line2 {
			marked[line] = color_mark.Sprintf("%v", lines[line])
		} else if line2 == s.where.line1 {
			marked[line] = color_mark.Sprintf("%v", lines[line])
		} else if line2 == s.where.line2 {
			marked[line] = color_mark.Sprintf("%v", lines[line])
		} else if line2 > s.where.line1 && line2 < s.where.line2 {
			marked[line] = color_mark.Sprintf("%v", lines[line])
		} else {
			marked[line] = lines[line]
		}
	}
	return marked
}

var color_leader = color.New(color.FgHiBlack)
var color_mark = color.New(color.Underline, color.Bold)

func (s Session) Print() {
	lines := s.Mark(s.GetLines())
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