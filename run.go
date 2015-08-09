package genregex

// 2 * the number of capture groups in a regex
const NumCaptures = 20
const NumGroups = NumCaptures / 2

type Thread struct {
	pc      int64
	saveReg *[NumCaptures]int
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
	cSet := make(map[int64]bool)
	nSet := make(map[int64]bool)
	cList := make([]Thread, 0, progLength)
	nList := make([]Thread, 0, progLength)
	cList = append(cList, Thread{0, &[NumCaptures]int{}})
	cSet[0] = true
	again := true

	res := make([]string, 0, NumGroups)

	for c, chr := range input {
	top:
		for i := 0; i < len(cList); i++ {
			thr := cList[i]
			inst := in[thr.pc]
			switch inst.Op {
			case Char:
				if inst.Char != chr {
					break
				}
				nList, nSet = insertKey(nList, nSet, Thread{thr.pc + 1, thr.saveReg})
			case Jump:
				cList, cSet = insertKey(cList, cSet, Thread{inst.Label1, thr.saveReg})
			case Split:
				cList, cSet = insertKey(cList, cSet, Thread{inst.Label1, thr.saveReg})
				cList, cSet = insertKey(cList, cSet, Thread{inst.Label2, thr.saveReg})
			case Save:
				thr.saveReg[inst.Label1] = c
				cList, cSet = insertKey(cList, cSet, Thread{thr.pc + 1, thr.saveReg})
			case Match:
				for j := 0; j < NumGroups; j++ {
					res = append(res, input[thr.saveReg[2*j]:thr.saveReg[2*j+1]])
				}
				return true, res
			case Nop:
				panic("Nop should be optimizied out")
			default:
				panic("invalid opcode")
			}
		}
		tmpL := cList
		cList, cSet = nList, nSet
		nList = tmpL[:0]
		nSet = make(map[int64]bool)

		if c == len(input)-1 && again {
			// we need to perform the last iteration twice, to clear any remaining
			// match operations we can still access
			again = false
			goto top
		}
	}

	return false, nil
}
