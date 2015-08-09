package regex_benchmarks

import (
	"local/regen/benchmarks/pats/pat1"
	"local/regen/benchmarks/pats/pat2"
	"regexp"
	"testing"
)

//go:generate ../regen --pattern=a(a*)bcdef --func=MatchPat --package=pat1 --out=pats/pat1/pat1.go
func BenchmarkRegex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pat1.MatchPat("aaaaaaaaabcdef")
	}
}

func BenchmarkDefaultRegex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if re, err := regexp.Compile("a(a*)bcdef"); err == nil {
			re.FindAllStringSubmatch("aaaaaaaaabcdef", -1)
		}
	}
}

func BenchmarkDefaultRegexCompile(b *testing.B) {
	if re, err := regexp.Compile("a(a*)bcdef"); err == nil {
		for i := 0; i < b.N; i++ {
			re.FindAllStringSubmatch("aaaaaaaaabcdef", -1)
		}
	}
}

//go:generate ../regen --pattern=aa*bcdef --func=MatchPat2 --package=pat2 --out=pats/pat2/pat2.go
func BenchmarkRegexNoCapture(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pat2.MatchPat2("aaaaaaaaabcdef")
	}
}

func BenchmarkDefaultRegexNoCapture(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if re, err := regexp.Compile("a(a*)bcdef"); err == nil {
			re.MatchString("aaaaaaaaabcdef")
		}
	}
}

func BenchmarkDefaultRegexCompileNoCapture(b *testing.B) {
	if re, err := regexp.Compile("a(a*)bcdef"); err == nil {
		for i := 0; i < b.N; i++ {
			re.MatchString("aaaaaaaaabcdef")
		}
	}
}
