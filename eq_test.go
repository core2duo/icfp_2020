package main

import (
	"testing"
)

func TestEq(t *testing.T) {
	testString(t, "ap ap eq 0 -2", "f")
	testString(t, "ap ap eq 0 -1", "f")
	testString(t, "ap ap eq 0 0", "t")
	testString(t, "ap ap eq 0 1", "f")
	testString(t, "ap ap eq 0 2", "f")

	testString(t, "ap ap eq 1 -1", "f")
	testString(t, "ap ap eq 1 0", "f")
	testString(t, "ap ap eq 1 1", "t")
	testString(t, "ap ap eq 1 2", "f")
	testString(t, "ap ap eq 1 3", "f")

	testString(t, "ap ap eq 2 0", "f")
	testString(t, "ap ap eq 2 1", "f")
	testString(t, "ap ap eq 2 2", "t")
	testString(t, "ap ap eq 2 3", "f")
	testString(t, "ap ap eq 2 4", "f")

	testString(t, "ap ap eq 19 20", "f")
	testString(t, "ap ap eq 20 20", "t")
	testString(t, "ap ap eq 21 20", "f")

	testString(t, "ap ap eq -19 -20", "f")
	testString(t, "ap ap eq -20 -20", "t")
	testString(t, "ap ap eq -21 -20", "f")
}
