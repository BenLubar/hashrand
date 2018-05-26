package hashrand

import (
	"encoding"
	"encoding/binary"
	"errors"
)

var _ encoding.BinaryMarshaler = (*Source)(nil)
var _ encoding.BinaryUnmarshaler = (*Source)(nil)
var errInvalidData = errors.New("hashrand: invalid data")

// MarshalBinary encodes the current state of s into a []byte. The result does
// not contain any information about the hash function.
func (s *Source) MarshalBinary() ([]byte, error) {
	buf := make([]byte, 1+binary.MaxVarintLen64*3)
	buf[0] = 0 // format version number
	i := 1
	i += binary.PutUvarint(buf[i:], uint64(len(s.seed)))
	i += binary.PutUvarint(buf[i:], uint64(len(s.buf)))
	i += binary.PutUvarint(buf[i:], uint64(len(s.last)))
	buf = buf[:i]
	buf = append(buf, s.seed...)
	buf = append(buf, s.buf...)
	buf = append(buf, s.last...)
	return buf, nil
}

// UnmarshalBinary resets the state of s to the state encoded in b. It does
// not assign a value to Hash, so Hash should be assigned manually before
// using this Source to generate any random numbers.
func (s *Source) UnmarshalBinary(b []byte) error {
	if len(b) < 1 || b[0] != 0 {
		return errInvalidData
	}
	b = b[1:]

	seedLen, n := binary.Uvarint(b)
	if n <= 0 {
		return errInvalidData
	}
	b = b[n:]

	bufLen, n := binary.Uvarint(b)
	if n <= 0 {
		return errInvalidData
	}
	b = b[n:]

	lastLen, n := binary.Uvarint(b)
	if n <= 0 {
		return errInvalidData
	}
	b = b[n:]

	if uint64(len(b)) != seedLen+bufLen+lastLen {
		return errInvalidData
	}

	s.seed = make([]byte, seedLen)
	s.buf = make([]byte, bufLen)
	s.last = make([]byte, lastLen)

	b = b[copy(s.seed, b):]
	b = b[copy(s.buf, b):]
	b = b[copy(s.last, b):]

	if len(b) != 0 {
		panic("hashrand: internal error")
	}

	return nil
}
