package genregex

import "testing"

func TestRegexBasicMatch(t *testing.T) {
	testRegex := "aab*"
	regexParse := RegexParser{}
	if parseRegex, err := regexParse.Parse(testRegex); err == nil {
		inst := finalizeInst(parseRegex.compile())
		t.Log(testRegex, " -> ", inst)
		match := "aabbbbbbbbbbbbbbbb"
		if b, _ := ThompsonVM(inst, match); !b {
			t.Logf("pattern [%s] failed to match string: '%s'", testRegex, match)
			t.Fail()
		} else {
			t.Log("Success")
		}
		match = "aa"
		if b, _ := ThompsonVM(inst, match); !b {
			t.Logf("pattern [%s] failed to match string: '%s'", testRegex, match)
			t.Fail()
		} else {
			t.Log("Success")
		}
		match = "aba"
		if b, _ := ThompsonVM(inst, match); b {
			t.Logf("pattern [%s] incorrectly matched string: '%s'", testRegex, match)
			t.Fail()
		} else {
			t.Log("Success")
		}

	} else {
		t.Logf("Failed to parse regex %s", testRegex)
		t.Fail()
	}
}

//func TestRegexGen(t *testing.T) {
//testRegex := "aab*c+ddd"
//GenMatcher(testRegex, "TestRegexMatch")
//}
