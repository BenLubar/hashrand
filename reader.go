package hashrand

import "io"

// Read implements io.Reader. Works identically to rand.Rand.Read unless
// multiple instances of rand.Rand use the same Source.
func (s *Source) Read(b []byte) (int, error) {
	for len(s.buf) < len(b) {
		s.fill()
	}

	n := copy(b, s.buf)
	s.buf = s.buf[n:]
	return n, nil
}

var _ io.Reader = (*Source)(nil)
