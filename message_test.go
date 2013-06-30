package irc

import (
	"bytes"
	"fmt"
	"testing"
)

func TestParseMessageSimple(t *testing.T) {
	in := []byte("cmd")
	msg, err := ParseMessage(in)
	if err != nil {
		t.Errorf("ParseMessage(%q) failed: %v", in, err)
	}
	if x := msg.Prefix(); x != nil {
		t.Errorf("ParseMessage(%q).Prefix() = %q, want nil", in, x)
	}
	if x := msg.Command(); !bytes.Equal(x, in) {
		t.Errorf("ParseMessage(%q).Command() = %q, want %q", in, x, in)
	}
	if x := msg.Params(); len(x) != 0 {
		t.Errorf("ParseMessage(%q).Params() = %q, want []", in, x)
	}
}

func TestParseMessageComplex(t *testing.T) {
	in := []byte(":server!user@host cmd param1 param2 :last param\r\n")
	prefix, command := []byte("server!user@host"), []byte("cmd")
	params := `["param1" "param2" "last param"]`
	
	msg, err := ParseMessage(in)
	if err != nil {
		t.Errorf("ParseMessage(%q) failed: %v", in, err)
	}
	if x := msg.Prefix(); !bytes.Equal(x, prefix) {
		t.Errorf("ParseMessage(%q).Prefix() = %q, want %q", in, x, prefix)
	}
	if x := msg.Command(); !bytes.Equal(x, command) {
		t.Errorf("ParseMessage(%q).Command() = %q, want %q", in, x, command)
	}
	if x := fmt.Sprintf("%q", msg.Params()); x != params {
		t.Errorf("ParseMessage(%q).Command() = %s, want %s", in, x, params)
	}
}

func TestMessageBytes(t *testing.T) {
	msg, _ := NewMessage(
		[]byte("prefix"), []byte("cmd"), []byte("param1"), []byte("last param"))
	out := []byte(":prefix cmd param1 :last param")
	if x := msg.Bytes(); !bytes.Equal(x, out) {
		t.Errorf("Message.Bytes() = %q, want %q", x, out)
	}
}