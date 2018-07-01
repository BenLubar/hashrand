package hashrand

import (
	"encoding/binary"
	"hash"
	"math/rand"
	"runtime"
)

// Source implements rand.Source and rand.Source64. Source is not safe for
// concurrent use by multiple goroutines.
type Source struct {
	// Hash is the hash function used by this source. It should be treated
	// as immutable once the source has been used to generate random
	// numbers.
	//
	// If Hash is not explicitly set, it defaults to sha1.New().
	Hash hash.Hash

	seed []byte
	buf  []byte
	last []byte
}

var _ rand.Source64 = (*Source)(nil)

// Uint64 implements rand.Source64.
func (s *Source) Uint64() uint64 {
	for len(s.buf) < 8 {
		s.fill()
	}

	var b []byte
	b, s.buf = s.buf[:8], s.buf[8:]

	return binary.LittleEndian.Uint64(b)
}

// Int63 implements rand.Source.
func (s *Source) Int63() int64 {
	pc, _, _, ok := runtime.Caller(3)
	if ok {
		f := runtime.FuncForPC(pc - 1)
		if f.Name() == "math/rand.read" {
			// HACK: only return 7 bytes of entropy as rand.Read drops the top byte.
			for len(s.buf) < 7 {
				s.fill()
			}

			var b []byte
			b, s.buf = s.buf[:7:7], s.buf[7:]
			b = append(b, 0)

			return int64(binary.LittleEndian.Uint64(b))
		}
	}

	n := s.Uint64()
	n &^= (1 << 63) // strip top bit
	return int64(n)
}

// Seed implements rand.Source. It is equivalent to calling s.Reset
// followed by s.AppendSeed with the little-endian two's complement
// encoding of seed.
func (s *Source) Seed(seed int64) {
	s.Reset()
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(seed))
	s.AppendSeed(b[:])
}
