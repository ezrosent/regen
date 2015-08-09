# regen
Compile-time regular expression parsing for go-lang.

regen generates a custom regular expression matching function for a given
pattern. It can be used with `go generate`. This is currently a proof of
concept: it does not support all of the common regular expression
meta-characters and the matching implementation can be made much faster.

The current implementation is based on [Regular Expression Matching: the Virtual
Machine Approach](https://swtch.com/~rsc/regexp/regexp2.html) by Russ Cox.

#TODOs
* Implement character classes and wildcards.
* Implement different backends (like the D-lang regex library in phobos)
