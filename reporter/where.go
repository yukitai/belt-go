package reporter

import "fmt"

type Where struct {
	line1 uint
	line2 uint
	begin uint
	end uint
}

func min(x, y uint) uint {
	if x < y {
		return x
	}
	return y
}

func max(x, y uint) uint {
	if x > y {
		return x
	}
	return y
}

func WhereNew(line1, line2, begin, end uint) Where {
	return Where {
		line1, line2, begin, end,
	}
}

func (w1 *Where) Merge(w2 *Where) Where {
	line1 := min(w1.line1, w2.line1)
	line2 := max(w1.line2, w2.line2)
	begin := min(w1.begin, w2.begin)
	end := max(w1.begin, w2.end)
	return Where {
		line1, line2, begin, end,
	}
}

func (w *Where) ToString() string {
	return fmt.Sprintf("line: [%v:%v], char: [%v:%v]", w.line1, w.line2, w.begin, w.end)
}