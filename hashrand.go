// Package hashrand provides a random number source using hash functions.
package hashrand // import "github.com/BenLubar/hashrand"

func (s *Source) fill() {
	s.Hash.Reset()
	_, _ = s.Hash.Write(s.seed)
	_, _ = s.Hash.Write(s.last)
	s.last = s.Hash.Sum(s.last[:0])
	s.buf = append(s.buf, s.last...)
}

// Reset clears the state of s.
func (s *Source) Reset() {
	s.seed = nil
	s.buf = nil
	s.last = nil
}

// AppendSeed adds entropy to the current random seed. The new random seed
// will not take effect until the internal buffer is exhausted.
func (s *Source) AppendSeed(b []byte) {
	s.seed = append(s.seed, b...)
}
