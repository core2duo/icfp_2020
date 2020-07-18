package main

import (
	"testing"
)

func TestS(t *testing.T) {
	testString(t, "ap ap ap s add inc 1", "3")
	testString(t, "ap ap ap s mul ap add 1 6", "42")
}
