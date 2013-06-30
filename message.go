package irc

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
)

const (
	MaxMessageLength = 510

	cr = 13
	lf = 10
	space = 32
	colon = 58

	reNick = `[[:alpha:]][[:alnum:]\[\]\\^{}` + "`" + `-]*`
	reChan = `[#&][^ \a\0\r\n]+`
	reUser = `[^ \0\r\n]+`
	reHost = `[[:alnum:]](?:[[:alnum:].-]*[[:alnum:]])?`
	rePrefix =
		reHost + `|` + reNick + `(?:!` + reUser + `)?(?:@` + reHost + `)?`
	reCommand = `[[:alpha:]]+|[[:digit:]]{3}`
	reParam = `[^ \0\r\n:][^ \0\r\n]*`
	reParams = `((?: +` + reParam + `)*)(?: +:([^\0\r\n]*))?`
	reMessage = `^ *` +
		`(?::(` + rePrefix + `) +)?` +
		`(` + reCommand + `)` + reParams +
		`[\r\n]*$`
)

var (
	paramRegexp = regexp.MustCompile(reParam)
	
	NickRegexp = regexp.MustCompile(reNick)
	ChanRegexp = regexp.MustCompile(reChan)
	MessageRegexp = regexp.MustCompile(reMessage)
)

// Message represents an IRC message. See RFC 1459 [2.3] for a description of
// its format and contents.
type Message struct {
	data []byte
	prefix []byte
	command []byte
	params [][]byte
}

// NewMessage creates a Message for the given prefix, command, and params.
// Prefix may be nil or empty (which will be converted to nil).
func NewMessage(prefix, command []byte, params ...[]byte) (*Message, error) {
	if command == nil {
		return nil, fmt.Errorf("command cannot be nil")
	}
	if len(prefix) == 0 {
		prefix = nil
	}
	return &Message{prefix: prefix, command: command, params: params}, nil
}

// ParseMessage creates a Message by parsing IRC protocol data. Parsing is done
// with a complex Regexp; if this Regexp doesn't match ParseMessage will return
// an error.
func ParseMessage(data []byte) (*Message, error) {
	groups := MessageRegexp.FindSubmatch(data)
	if groups == nil {
		return nil, fmt.Errorf("could not parse message %q", data)
	}

	prefix := groups[1]
	if len(prefix) == 0 {
		prefix = nil
	}

	params := paramRegexp.FindAll(groups[3], -1)
	if params == nil {
		params = [][]byte{}
	}
	if len(groups[4]) != 0 {
		params = append(params, groups[4])
	}
	
	return &Message{
		data: data,
		prefix: prefix,
		command: groups[2],
		params: params,
	}, nil
}

// Bytes returns a byte array containing the encoded Message.
func (m *Message) Bytes() []byte {
	if m.data == nil {
		buf := bytes.NewBuffer(make([]byte, 0, MaxMessageLength))
		if m.prefix != nil {
			buf.WriteByte(colon)
			buf.Write(m.prefix)
			buf.WriteByte(space)
		}
		buf.Write(m.command)
		for i, param := range m.params {
			buf.WriteByte(space)
			if i == len(m.params) - 1 {
				buf.WriteByte(colon)
			}
			buf.Write(param)
		}
		m.data = buf.Bytes()
	}
	return m.data
}

// WriteTo implements the io.WriterTo interface. It writes the encoded Message
// and CRLF ("\r\n") to w.
func (m *Message) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(append(m.Bytes(), cr, lf))
	return int64(n), err
}

func (m *Message) Prefix() []byte { return m.prefix }
func (m *Message) Command() []byte { return m.command }
func (m *Message) Params() [][]byte { return m.params }