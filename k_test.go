package main

import (
	"testing"
)

func TestK(t *testing.T) {
	testString(t, "ap ap t 1 5", "1")
	testString(t, "ap ap t t i", "t/2( )")
	testString(t, "ap ap t t ap inc 5", "t/2( )")
	testString(t, "ap ap t ap inc 5 t", "6")
}
