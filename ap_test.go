package main

import (
	"testing"
)

func TestAp(t *testing.T) {
	testString(t, "ap inc ap inc 0", "2")
	testString(t, "ap inc ap inc ap inc 0", "3")
	testString(t, "ap ap add ap ap add 2 3 4", "9")
	testString(t, "ap ap add 2 ap ap add 3 4", "9")
	testString(t, "ap ap add ap ap mul 2 3 4", "10")
	testString(t, "ap ap mul 2 ap ap add 3 4", "14")
}
