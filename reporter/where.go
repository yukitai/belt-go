package reporter

import "fmt"

type Where struct {
	line1 uint
	line2 uint
	begin uint
	end   uint
	fake  bool
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
		line1: line1, 
		line2: line2, 
		begin: begin, 
		end:   end,
	}
}

func (w *Where) Clone() Where {
	return Where {
		line1: w.line1, 
		line2: w.line2, 
		begin: w.begin, 
		end: w.end,
	}
}

func (w1 *Where) Merge(w2 *Where) Where {
	if w2 == nil {
		return w2.Clone()
	}
	line1 := min(w1.line1, w2.line1)
	line2 := max(w1.line2, w2.line2)
	begin := min(w1.begin, w2.begin)
	end := max(w1.begin, w2.end)
	fake := false
	return Where {
		line1, line2, begin, end, fake,
	}
}

func (w *Where) ToString() string {
	return fmt.Sprintf("line: [%v:%v], char: [%v:%v]", w.line1, w.line2, w.begin, w.end)
}

func FakeWhere() Where {
	return Where{
		fake: true,
	}
}

func (w *Where) IsFake() bool {
	return w.fake
}