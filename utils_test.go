package main

import (
	"testing"
)

func TestSplit2(t *testing.T) {
	subject := "a.b.c"
	separator := "."

	a, bc := Split2(subject, separator)
	if a != "a" || bc != "b.c" {
		t.Fail()
	}

	b, c := Split2(bc, separator)
	if b != "b" || c != "c" {
		t.Fail()
	}
}
