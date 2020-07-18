package main

import (
	"testing"
)

func assertModemInt(t *testing.T, num int64, s string) {
	atom, rest := demodulate(modulate(&Number{Value: num}))
	if rest != "" {
		t.Errorf("rest: %q", rest)
	}
	if atom.(*Number).Value != num {
		t.Errorf("ap dem ap mod %d != %d", atom.(*Number).Value, num)
	}

	atom, rest = demodulate(s)
	if rest != "" {
		t.Errorf("rest: %q", rest)
	}
	if modulate(atom) != s {
		t.Errorf("ap mod ap dem %d test has failed", num)
	}
}

func TestModem(t *testing.T) {
	assertModemInt(t, 0, "010")
	assertModemInt(t, 1, "01100001")
	assertModemInt(t, -1, "10100001")
	assertModemInt(t, 2, "01100010")
	assertModemInt(t, -2, "10100010")

	assertModemInt(t, 16, "0111000010000")
	assertModemInt(t, -16, "1011000010000")

	assertModemInt(t, 255, "0111011111111")
	assertModemInt(t, -255, "1011011111111")
	assertModemInt(t, 256, "011110000100000000")
	assertModemInt(t, -256, "101110000100000000")
}
