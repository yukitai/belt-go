package frontend

import (
	"belt/reporter"
	"fmt"
)

type UnionFindSet struct {
	parents  []uint
	ranks    []uint
	values   []*AstValType
	typevars map[string]uint
}

func UFSNew() UnionFindSet {
	return UnionFindSet{
		parents:  make([]uint, 0),
		ranks:    make([]uint, 0),
		values:   make([]*AstValType, 0),
		typevars: make(map[string]uint),
	}
}

func (ufs *UnionFindSet) Find(x uint) uint {
	if(ufs.parents[x] == x) {
		return x
	}
	ufs.parents[x] = ufs.Find(ufs.parents[x])     
    return ufs.parents[x]
}

func (ufs *UnionFindSet) Merge(a uint, b uint) {
	a_root := ufs.values[a].IsLlType()
	b_root := ufs.values[b].IsLlType()
    x := ufs.Find(a)
    y := ufs.Find(b)
    if x == y {
		return
	}
    if (ufs.ranks[x] > ufs.ranks[y] || a_root) && !b_root {
		ufs.parents[y]=x
	} else {
        if ufs.ranks[x] == ufs.ranks[y] {
			ufs.ranks[y] += 1
		}
        ufs.parents[x]=y
    }
}

func (ufs *UnionFindSet) Extend(value *AstValType) uint {
	if value.Vttype == ANTVar {
		tvarname := value.Item.(*AstValTypeVar).Ident.value
		tvar, ok := ufs.typevars[tvarname]
		if ok {
			return tvar
		}
		size := uint(len(ufs.values))
		ufs.parents = append(ufs.parents, size)
		ufs.ranks = append(ufs.ranks, 0)
		ufs.values = append(ufs.values, value)
		ufs.typevars[tvarname] = size
		return size
	}
	size := uint(len(ufs.values))
	ufs.parents = append(ufs.parents, size)
	ufs.ranks = append(ufs.ranks, 0)
	ufs.values = append(ufs.values, value)
	return size
}

func (ufs *UnionFindSet) ExtendTVar() uint {
	return ufs.Extend(&AstValType{
		Vttype: ANTVar,
		Item: &AstValTypeVar{
			Real: nil,
		},
	})
}

func (ufs *UnionFindSet) MakeEffect(a *Analyzer) {
	for i, item := range ufs.values {
		i := uint(i)
		t := ufs.Find(i)
		ty := *ufs.values[t]
		if ty.IsLlType() {
			if item.IsLlType() {
				where := item.Where()
				if !where.IsFake() {
					lt := item.ToString()
					rt := ty.ToString()
					if lt != rt {
						err := reporter.Error(
							where,
							fmt.Sprintf("mismatched type %v and type %v", lt, rt),
						)
						reporter.Report(&err, a.file)
						a.has_err = true
					}
				}
			} else {
				*item = ty
			}
		} else {
			where := item.Where()
			if !where.IsFake() {
				err := reporter.Error(
					where,
					fmt.Sprintf("type %v cannot be known at the compile time", item.ToString()),
				)
				reporter.Report(&err, a.file)
				a.has_err = true
			}
		}
	}
}