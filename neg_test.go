package main

import (
	"testing"
)

func TestNeg(t *testing.T) {
	testString(t, "ap neg 0", "0")
	testString(t, "ap neg 1", "-1")
	testString(t, "ap neg -1", "1")
	testString(t, "ap neg 2", "-2")
	testString(t, "ap neg -2", "2")
}
