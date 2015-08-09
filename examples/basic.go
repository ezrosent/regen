package main

import "fmt"

//go:generate ../regen --pattern=a(a*)bcdef --func=MatchPat --out=pat1.go

func main() {
	ok, capt := MatchPat("aaabcdef")
	if ok {
		fmt.Println(capt)
	} else {
		fmt.Println("failure")
	}
}
