package hashrand_test

import (
	"io"
	"math/rand"
	"testing"
	"testing/quick"

	"github.com/BenLubar/hashrand"
)

func TestRead(t *testing.T) {
	if err := quick.CheckEqual(func(seed int64) ([]byte, error) {
		var s hashrand.Source
		b := make([]byte, 1000)
		s.Seed(seed)
		r := rand.New(&s)
		_, err := io.ReadFull(r, b)
		return b, err
	}, func(seed int64) ([]byte, error) {
		var s hashrand.Source
		b := make([]byte, 1000)
		s.Seed(seed)
		_, err := io.ReadFull(&s, b)
		return b, err
	}, nil); err != nil {
		t.Error(err)
	}
}
