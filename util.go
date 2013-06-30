package irc

import "bytes"

// ToLower returns a copy of the byte slice s converted to IRC lower case as
// described in RFC 1459 [2.2].
func ToLower(s []byte) []byte {
	out := make([]byte, len(s))
	for i, b := range s {
		if b >= 65 && b <= 93 { // 'A' .. ']'
			b += 32             // 'a' .. '}'
		}
		out[i] = b 
	}
	return out
}

// EqualFold reports whether s and t are equal under IRC case equivalence.
func EqualFold(s, t []byte) bool {
	return bytes.Equal(ToLower(s), ToLower(t))
}