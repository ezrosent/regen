package genregex

type regexParse interface {
	compile() []*inst
}

type regexParser struct {
	prev regexParse
	rest []regexParse
}

func (r *regexParser) Parse(b []byte) regexParse {
	for _, char := range b {
		switch char {
		case '?':
			r.prev = optional{r.prev}
		case '*':
			r.prev = many{r.prev}
		case '+':
			r.prev = concat{[]regexParse{r.prev, many{r.prev}}}
		default:
			r.rest = append(r.rest, r.prev)
			r.prev = constant{char}
		}
	}
	return concat{append(r.rest, r.prev)}
}

type OpCode uint

const (
	Char OpCode = iota
	Match
	Jump
	Split
	Nop
)

// intermediate version that has labels as pointers
type inst struct {
	op     OpCode
	char   byte  // used for char
	label1 *inst // used by jump and split
	label2 *inst // used by split
	index  int64
}

func nop() *inst {
	return &inst{op: Nop}
}

type optional struct {
	option regexParse
}

func (o optional) compile() []*inst {
	nop1 := nop() // before fldinst
	nop2 := nop() // after fldinst
	split := &inst{op: Split, label1: nop1, label2: nop2}
	fldinst := o.option.compile()
	return append(append([]*inst{split, nop1}, fldinst...), nop2)
}

type many struct {
	field regexParse
}

func (m many) compile() []*inst {
	nop1 := nop() // before fldinst
	nop2 := nop() // after fldinst
	split := &inst{op: Split, label1: nop1, label2: nop2}
	fldinst := m.field.compile()
	return append(append(append([]*inst{split, nop1}, fldinst...),
		&inst{op: Jump, label1: nop1}), nop2)
}

type constant struct {
	field byte
}

func (c constant) compile() []*inst {
	return []*inst{&inst{op: Char, char: c.field}}
}

type concat struct {
	sequence []regexParse
}

func (c concat) compile() []*inst {
	ret := []*inst{}
	for _, regex := range c.sequence {
		ret = append(ret, regex.compile()...)
	}
	return ret
}

// final version that has indices for labels and no Nops
type Inst struct {
	op     OpCode
	char   byte
	label1 int64
	label2 int64
}

func nextLabel(instructs []*inst, inst *inst) *inst {
	ret := instructs[inst.label1.index+1]

	// hopefully cannot get a nop cycle...
	if ret.op == Nop {
		return nextLabel(instructs, ret)
	} else {
		return ret
	}
}

func finalizeInst(instructs []*inst) []Inst {
	for i, inst := range instructs {
		inst.index = int64(i)
	}

	instructs = append(instructs, &inst{op: Match})

	//if a label points to a nop, point it to the following instruction
	for _, inst := range instructs {
		switch inst.op {
		case Jump:
			if inst.label1.op == Nop {
				inst.label1 = nextLabel(instructs, inst.label1)
			}
		case Split:
			if inst.label1.op == Nop {
				inst.label1 = nextLabel(instructs, inst.label1)
			}

			if inst.label2.op == Nop {
				inst.label2 = nextLabel(instructs, inst.label2)
			}
		}
	}

	//remove the nops
	instructsNew := []*inst{}
	for _, inst := range instructs {
		if inst.op != Nop {
			instructsNew = append(instructsNew, inst)
			inst.index = int64(len(instructsNew) - 1)
		}
	}

	ret := []Inst{}
	for _, inst := range instructs {
		l1 := int64(0)
		l2 := int64(0)
		if inst.label1 != nil {
			l1 = inst.label1.index
		}
		if inst.label2 != nil {
			l2 = inst.label2.index
		}
		ret = append(ret, Inst{
			op:     inst.op,
			char:   inst.char,
			label1: l1,
			label2: l2,
		})
	}
	return ret
}
