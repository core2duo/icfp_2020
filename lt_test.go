package main

import (
	"testing"
)

func TestLt(t *testing.T) {
	testString(t, "ap ap lt 0 -1", "f")
	testString(t, "ap ap lt 0 0", "f")
	testString(t, "ap ap lt 0 1", "t")
	testString(t, "ap ap lt 0 2", "t")

	testString(t, "ap ap lt 1 0", "f")
	testString(t, "ap ap lt 1 1", "f")
	testString(t, "ap ap lt 1 2", "t")
	testString(t, "ap ap lt 1 3", "t")

	testString(t, "ap ap lt 2 1", "f")
	testString(t, "ap ap lt 2 2", "f")
	testString(t, "ap ap lt 2 3", "t")
	testString(t, "ap ap lt 2 4", "t")

	testString(t, "ap ap lt 19 20", "t")
	testString(t, "ap ap lt 20 20", "f")
	testString(t, "ap ap lt 21 20", "f")

	testString(t, "ap ap lt -19 -20", "f")
	testString(t, "ap ap lt -20 -20", "f")
	testString(t, "ap ap lt -21 -20", "t")
}
