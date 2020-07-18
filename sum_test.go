package main

import (
	"testing"
)

func TestAdd(t *testing.T) {
	test(t, "ap ap add 1 2", "3")
	test(t, "ap ap add 2 1", "3")
	test(t, "ap ap add 0 1", "1")
	test(t, "ap ap add 2 3", "5")
	test(t, "ap ap add 3 5", "8")
}
