package genregex

// 2 * the number of capture groups in a regex
const NumCaptures = 20
const NumGroups = NumCaptures / 2

type Thread struct {
	pc      int64
	saveReg [NumCaptures]int
}

func insertKey(list []Thread, set map[int64]bool, t Thread) ([]Thread, map[int64]bool) {
	if set[t.pc] {
		return list, set
	} else {
		set[t.pc] = true
		return append(list, t), set
	}
}

func ThompsonVM(in []Inst, input string) (bool, []string) {
	progLength := len(in)
	//TODO(ezrosent) need to make adding to cList and nList not insert duplicates
	//can probably just do this with a map[[]Inst][bool] or something similar
	cSet := make(map[int64]bool)
	nSet := make(map[int64]bool)
	cList := make([]Thread, 0, progLength)
	nList := make([]Thread, 0, progLength)
	cList = append(cList, Thread{0, [NumCaptures]int{}})

	res := make([]string, 0, NumGroups)

	for p, c := range input {
		for i := 0; i < len(cList); i++ {
			thr := cList[i]
			inst := in[thr.pc]
			switch inst.Op {
			case Char:
				if c != inst.Char {
					break
				}
				nList, nSet = insertKey(nList, nSet, Thread{thr.pc + 1, [NumCaptures]int{}})
			case Match:
				for j := 0; j < len(input); j += 2 {
					res = append(res, input[thr.saveReg[j]:thr.saveReg[j+1]])
				}
				return true, res
			case Jump:
				cList, cSet = insertKey(cList, cSet, Thread{inst.Label1, [NumCaptures]int{}})
			case Split:
				cList, cSet = insertKey(cList, cSet, Thread{inst.Label1, [NumCaptures]int{}})
				cList, cSet = insertKey(cList, cSet, Thread{inst.Label2, [NumCaptures]int{}})
			case Save:
				cList[i].saveReg[inst.Label1] = p
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
			for j := 0; j < len(input); j += 2 {
				res = append(res, input[cList[i].saveReg[j]:cList[i].saveReg[j+1]])
			}
			return true, res
		case Jump:
			cList, cSet = insertKey(cList, cSet, Thread{inst.Label1, [NumCaptures]int{}})
		case Split:
			cList, cSet = insertKey(cList, cSet, Thread{inst.Label1, [NumCaptures]int{}})
			cList, cSet = insertKey(cList, cSet, Thread{inst.Label2, [NumCaptures]int{}})
		case Save:
			cList[i].saveReg[inst.Label1] = len(input) - 1
		case Nop:
			panic("nop in final instruction stream")
		default:
			panic("invalid opcode")
		}
	}
	return false, nil
}
