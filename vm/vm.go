package vm

type VM struct {
	pc uint
	instr []Instruction
	ctx []VMContext
}

type VMScope struct {
	items map[uint]any
}

type VMContext struct {
	name string
	scopes []VMScope
}

func VMNew(instr []Instruction) VM {
	return VM{
		pc: 0,
		instr: instr,
		ctx: make([]VMContext, 0),
	}
}

func (vm *VM) Entry(entry string) any {
	return 0
}