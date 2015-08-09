package main

import "testing"

const NumCaptures = 20

func TestRegexBasicMatch(t *testing.T) {
	testRegex := "aab*"
	regexParse := regexParser{}
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

func TestCaptureGroups(t *testing.T) {
	testRegex := "ab(a*)b"
	regexParse := regexParser{}
	if parseRegex, err := regexParse.Parse(testRegex); err == nil {
		inst := finalizeInst(parseRegex.compile())
		t.Log(testRegex, " -> ", inst)
		match := "abaaaaaaaaab"
		if ok, capt := ThompsonVM(inst, match); ok {
			if capt[0] != "aaaaaaaaa" {
				t.Logf("pattern %s correctly matched string %s but instead of capturing '%s' captured '%s'",
					testRegex, match, "aaaaaaaaa", capt[0])
				t.Fail()
			} else {
				t.Log("Success")
			}
		} else {
			t.Logf("pattern [%s] failed to match string: '%s'", testRegex, match)
			t.Fail()
		}
	} else {
		t.Logf("Failed to parse regex %s: %v", testRegex, err)
		t.Fail()
	}
}
