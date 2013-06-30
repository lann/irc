package irc

import  (
	"bufio"
	"io"
)

// Scanner provides an interface for reading IRC messages from a Reader.
// Messages longer than 510 bytes will be truncated. Truncation can be detected
// by checking the Truncating field after a successful Scan.
type Scanner struct {
	*bufio.Scanner
	Truncating bool
}

func NewScanner(r io.Reader) *Scanner {
	scanner := Scanner{Scanner: bufio.NewScanner(r)}
	scanner.Split(scanner.split)
	return &scanner
}

func (s *Scanner) split(data []byte, eof bool) (int, []byte, error) {
	if eof && len(data) == 0 {
		return 0, nil, nil
	}

	offset := 0

	// Finish truncating long message
	if s.Truncating {
		for ; offset < len(data); offset++ {
			if c := data[offset]; c == lf || c == cr {
				s.Truncating = false
				break
			}
		}
	}
	
	// Strip leading whitespace, including empty messages
	for ; offset < len(data); offset++ {
		if c := data[offset]; c != lf && c != cr && c != space {
			break
		}
	}

	// Scan for a message delimiter
	for i := offset; i < len(data); i++ {
		if c := data[i]; c == lf || c == cr {
			// Message complete
			return i + 1, data[offset:i], nil
		} else if i - offset >= MaxMessageLength {
			// Message too long; truncate
			s.Truncating = true
			return i + 1, data[offset:i], nil
		}
	}

	return offset, nil, nil
}