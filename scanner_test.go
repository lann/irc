package irc

import (
	"bytes"
	"testing"
)

func TestScanner(t *testing.T) {
	in := append(
		bytes.Repeat([]byte("x"), MaxMessageLength * 3),
		[]byte("\nline one\r\nline two \n  line three\rline four\n")...)
	outs := bytes.Split(
		append(
			bytes.Repeat([]byte("x"), MaxMessageLength),
			[]byte(",line one,line two ,line three,line four")...),
		[]byte(","))

	scanner := NewScanner(bytes.NewReader(in))
	for _, out := range outs {
		if !scanner.Scan() {
			t.Errorf("Scan() returned false, Err() = %v", scanner.Err())
		}
		if x := scanner.Bytes(); !bytes.Equal(x, out) {
			t.Errorf("Bytes() = %q, want %q", x, out)
		}
	}
}