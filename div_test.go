package main

import (
	"testing"
)

func TestDiv(t *testing.T) {
	test(t, "ap ap div 4 2", "2")
	test(t, "ap ap div 4 3", "1")
	test(t, "ap ap div 4 4", "1")
	test(t, "ap ap div 4 5", "0")
	test(t, "ap ap div 5 2", "2")
	test(t, "ap ap div 6 -2", "-3")
	test(t, "ap ap div 5 -3", "-1")
	test(t, "ap ap div -5 3", "-1")
	test(t, "ap ap div -5 -3", "1")
}
