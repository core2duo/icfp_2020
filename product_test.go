package main

import (
	"testing"
)

func TestMul(t *testing.T) {
	test(t, "ap ap mul 4 2", "8")
	test(t, "ap ap mul 3 4", "12")
	test(t, "ap ap mul 3 -2", "-6")
}
