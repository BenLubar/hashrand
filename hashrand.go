// Package hashrand provides a random number source using hash functions.
package hashrand // import "github.com/BenLubar/hashrand"

import (
	"crypto/sha1"
	"io"
)

func (s *Source) fill() {
	if s.Hash == nil {
		s.Hash = sha1.New()
	}

	// Don't output the raw hash of the seed first.
	if s.last == nil {
		_, _ = s.Hash.Write(s.seed)
		s.last = s.Hash.Sum(nil)
	}

	s.Hash.Reset()
	_, _ = s.Hash.Write(s.last)
	_, _ = s.Hash.Write(s.seed)
	s.last = s.Hash.Sum(s.last[:0])
	if len(s.last) == 0 {
		panic(io.ErrNoProgress)
	}
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
