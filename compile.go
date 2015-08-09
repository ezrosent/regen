package genregex

import "fmt"

type regexParse interface {
	compile() []*inst
}

type RegexParser struct {
	prev regexParse
	rest []regexParse
}

func (r RegexParser) Parse(b string) (regexParse, error) {
	c := 0
	rp, _, err := r.parse(b, false, &c)
	return rp, err
}

//TODO(ezrosent) support escaping
func (r RegexParser) parse(b string, inCapture bool, matchCounter *int) (regexParse, int, error) {
	errorStr := func(r rune) error {
		return fmt.Errorf("must have previous regex before using metacharacter %c", r)
	}
	lastIndex := len(b) - 1
	for i := 0; i < len(b); i++ {
		char := rune(b[i])
		switch char {
		case '?':
			if r.prev == nil {
				return nil, 0, errorStr(char)
			}
			r.prev = optional{r.prev}
		case '*':
			if r.prev == nil {
				return nil, 0, errorStr(char)
			}
			r.prev = many{r.prev}
		case '+':
			if r.prev == nil {
				return nil, 0, errorStr(char)
			}
			r.prev = concat{[]regexParse{r.prev, many{r.prev}}}
		case '(':
			capture := capture{
				index: *matchCounter,
			}
			*matchCounter++
			if i == len(b)-1 {
				return nil, 0, fmt.Errorf("mismatched parens, (")
			}
			var rr RegexParser
			if capt, j, err := rr.parse(b[i+1:], true, matchCounter); err == nil {
				capture.field = capt
				r.rest = append(r.rest, r.prev)
				r.prev = capture
				i += j + 1
			} else {
				return nil, 0, err
			}
		case ')':
			if inCapture {
				// we set it to i because it will be incremented by the outer for loop
				lastIndex = i
				goto out
			} else {
				return nil, 0, fmt.Errorf("mismatched parens, )")
			}
		default:
			if r.prev != nil {
				r.rest = append(r.rest, r.prev)
			}
			r.prev = constant{char}
		}
	}

out:
	return concat{append(r.rest, r.prev)}, lastIndex, nil
}

//go:generate stringer -type=OpCode
type OpCode uint

const (
	Char OpCode = iota
	Match
	Jump
	Split
	Save
	Nop
)

// intermediate version that has labels as pointers
type inst struct {
	op        OpCode
	char      rune  // used for char
	label1    *inst // used by jump and split
	label2    *inst // used by split
	saveIndex int   // used by save
	index     int64
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

type capture struct {
	index int
	field regexParse
}

func (c capture) compile() []*inst {
	save1 := &inst{op: Save, saveIndex: 2 * c.index}
	save2 := &inst{op: Save, saveIndex: 2*c.index + 1}
	return append(append([]*inst{save1}, c.field.compile()...), save2)
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
		&inst{op: Jump, label1: split}), nop2)
}

type constant struct {
	field rune
}

func (c constant) compile() []*inst {
	return []*inst{&inst{op: Char, char: c.field}}
}

type concat struct {
	sequence []regexParse
}

func (c concat) compile() []*inst {
	ret := make([]*inst, 0, 10)
	for _, regex := range c.sequence {
		ret = append(ret, regex.compile()...)
	}
	return ret
}

// final version that has indices for labels and no Nops
type Inst struct {
	Op     OpCode
	Char   rune  // used by char
	Label1 int64 // used by jmp, split, save
	Label2 int64 // used by split
}

func (i Inst) String() string {
	return fmt.Sprintf("Inst{%s, %c, %d, %d}", i.Op.String(), i.Char, i.Label1, i.Label2)
}

func nextLabel(instructs []*inst, inst *inst) *inst {
	ret := instructs[inst.index+1]

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
	for _, inst := range instructsNew {
		l1 := int64(0)
		l2 := int64(0)
		if inst.label1 != nil {
			l1 = inst.label1.index
		}
		if inst.label2 != nil {
			l2 = inst.label2.index
		}
		if inst.op == Save {
			l1 = int64(inst.saveIndex)
		}
		ret = append(ret, Inst{
			Op:     inst.op,
			Char:   inst.char,
			Label1: l1,
			Label2: l2,
		})
	}
	return ret
}
