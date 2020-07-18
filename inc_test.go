package main

import (
	"testing"
)

func TestInc(t *testing.T) {
	env = make(map[string]Atom)
	expr := parse(map[string][]string{}, []string{"ap", "inc", "0"})
	if len(expr) != 1 {
		t.Errorf("unexpected parse result: %#v", expr)
	}
	s := Stack{}
	s = expr[0].Evaluate(s)
	if len(s) != 1 {
		t.Errorf("unexpected stack: %#v", expr)
	}
	if (Number{Value: 1}) != *s[0].(*Number) {
		t.Errorf("assertion failed: %#v != 1", s[0])
	}
}
