package irc

import (
	"bytes"
	"testing"
)

func TestToLower(t *testing.T) {
	in, out := []byte("Test{[\\]}@~"), []byte("test{{|}}@~")
	if x := ToLower(in); !bytes.Equal(x, out) {
		t.Errorf("ToLower(%q) = %q, want %q", in, x, out)
	}
}

func TestEqualFold(t *testing.T) {
	a, b := []byte("Test{[\\]}@~"), []byte("test{{|}}@~")
	if !EqualFold(a, b) {
		t.Errorf("EqualFold(%q, %q) should be true", a, b)
	}

	b = []byte("other")
	if EqualFold(a, b) {
		t.Errorf("EqualFold(%q, %q) should be false", a, b)
	}
}