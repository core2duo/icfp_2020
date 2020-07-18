package main

import (
	"testing"
)

func assertModInt(t *testing.T, num int64, v string) {
	if modulate(&Number{Value: num}) != v {
		t.Errorf("ap mod %d != %s", num, v)
	}
}

func TestMod(t *testing.T) {
	assertModInt(t, 0, "010")
	assertModInt(t, 1, "01100001")
	assertModInt(t, -1, "10100001")
	assertModInt(t, 2, "01100010")
	assertModInt(t, -2, "10100010")

	assertModInt(t, 16, "0111000010000")
	assertModInt(t, -16, "1011000010000")

	assertModInt(t, 255, "0111011111111")
	assertModInt(t, -255, "1011011111111")
	assertModInt(t, 256, "011110000100000000")
	assertModInt(t, -256, "101110000100000000")
}
