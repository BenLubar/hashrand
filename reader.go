package hashrand

import "io"

// Read implements io.Reader. Unlike rand.Rand.Read, this method does not skip
// every eighth byte.
func (s *Source) Read(b []byte) (int, error) {
	for len(s.buf) < len(b) {
		s.fill()
	}

	n := copy(b, s.buf)
	s.buf = s.buf[n:]
	return n, nil
}

var _ io.Reader = (*Source)(nil)
