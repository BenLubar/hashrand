package hashrand_test

import (
	"io"
	"testing"

	"github.com/BenLubar/hashrand"
)

func TestBadHash(t *testing.T) {
	defer func() {
		if r := recover(); r != io.ErrNoProgress {
			t.Errorf("expected io.ErrNoProgress but got %v", r)
		}
	}()

	s := &hashrand.Source{Hash: &badHash{}}
	s.Int63()
}

type badHash struct{}

func (*badHash) BlockSize() int                    { return 1 }
func (*badHash) Reset()                            {}
func (*badHash) Size() int                         { return 0 }
func (*badHash) Sum(b []byte) []byte               { return b }
func (*badHash) Write(b []byte) (n int, err error) { return len(b), nil }

func TestBadUnmarshal(t *testing.T) {
	cases := []struct {
		Name string
		Data []byte
	}{
		{
			Name: "Version",
			Data: []byte{1},
		},
		{
			Name: "VarInt",
			Data: []byte{0, 255},
		},
		{
			Name: "SliceLength",
			Data: []byte{0, 1, 0, 0},
		},
	}

	for _, c := range cases {
		c := c // shadow
		t.Run(c.Name, func(t *testing.T) {
			var s hashrand.Source
			if err := s.UnmarshalBinary(c.Data); err == nil {
				t.Error("unexpected success: ", c.Name)
			}
		})
	}
}
