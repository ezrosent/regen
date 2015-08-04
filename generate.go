package genregex

import (
	"os"
	"text/template"
)

//TODO(ezrosent): gobindata for source code
//either: capture groups, or actual go generate
type parseGenData struct {
	PackageName    string
	InterpFunc     string
	InterpFuncName string
	OpCodeT        string
	InstructT      string
	FuncName       string
	Instructions   []Inst
}

const instruct = `
type Inst struct {
	Op     OpCode
	Char   rune
	Label1 int64
	Label2 int64
}
`
const opcode = `
type OpCode uint
const (
	Char OpCode = iota
	Match
	Jump
	Split
	Nop
)

`

const tompsonvm = `
type Thread struct {
	pc int64
}

func insertKey(list []Thread, set map[int64]bool, t Thread) ([]Thread, map[int64]bool) {
	if set[t.pc] {
		return list, set
	} else {
		set[t.pc] = true
		return append(list, t), set
	}
}


func __ThompsonVM(in []Inst, input string) bool {
	progLength := len(in)
	//TODO(ezrosent) need to make adding to cList and nList not insert duplicates
	//can probably just do this with a map[[]Inst][bool] or something similar
	cSet := make(map[int64]bool)
	nSet := make(map[int64]bool)
	cList := make([]Thread, 0, progLength)
	nList := make([]Thread, 0, progLength)
	cList = append(cList, Thread{0})

	for _, c := range input {
		for i := 0; i < len(cList); i++ {
			inst := in[cList[i].pc]
			switch inst.Op {
			case Char:
				if c != inst.Char {
					break
				}
				nList, nSet = insertKey(nList, nSet, Thread{cList[i].pc + 1})
			case Match:
				return true
			case Jump:
				cList, cSet = insertKey(cList, cSet, Thread{inst.Label1})
			case Split:
				cList, cSet = insertKey(cList, cSet, Thread{inst.Label1})
				cList, cSet = insertKey(cList, cSet, Thread{inst.Label2})
			case Nop:
				panic("nop in final instruction stream")
			default:
				panic("invalid opcode")
			}
		}
		tmp := cList
		cList = nList
		nList = tmp[:0]
	}

	// No more characters to consume, but we still need one more pass
	// to determine if we would have terminated appropriately
	for i := 0; i < len(cList); i++ {
		inst := in[cList[i].pc]
		switch inst.Op {
		case Char:
		case Match:
			return true
		case Jump:
			cList, cSet = insertKey(cList, cSet, Thread{inst.Label1})
		case Split:
			cList, cSet = insertKey(cList, cSet, Thread{inst.Label1})
			cList, cSet = insertKey(cList, cSet, Thread{inst.Label2})
		case Nop:
			panic("nop in final instruction stream")
		default:
			panic("invalid opcode")
		}
	}

	return false
}
`

var genTemplate = template.Must(template.New("matcher").Parse(`
package {{.PackageName}}
{{.OpCodeT}}

{{.InstructT}}

{{.InterpFunc}}

func {{.FuncName}}(input string) bool {
	in := []Inst{ {{range .Instructions}}
	    Inst{ {{.Op}}, {{.Char}}, {{.Label1}}, {{.Label2}}},
	{{end}} }

	return {{.InterpFuncName}}(in, input)
}
`))

func GenMatcher(regex string, name string) error {
	parser := new(RegexParser)
	if inst, err := parser.Parse(regex); err != nil {
		return err
	} else {
		genData := parseGenData{
			PackageName:    "test",
			InterpFunc:     tompsonvm,
			InterpFuncName: "__ThompsonVM",
			OpCodeT:        opcode,
			InstructT:      instruct,
			FuncName:       name,
			Instructions:   finalizeInst(inst.compile()),
		}
		return genTemplate.Execute(os.Stdout, genData)
	}
}
