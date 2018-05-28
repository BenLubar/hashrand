package hashrand_test

import (
	"math/rand"
	"testing"
	"testing/quick"

	"github.com/BenLubar/hashrand"
)

func TestMarshalEmpty(t *testing.T) {
	s1 := &hashrand.Source{}
	b, err := s1.MarshalBinary()
	if err != nil {
		t.Error(err)
	}

	s2 := &hashrand.Source{}
	err = s2.UnmarshalBinary(b)
	if err != nil {
		t.Error(err)
	}
}

func TestMarshal(t *testing.T) {
	err := quick.CheckEqual(func(seed int64) []int {
		s := &hashrand.Source{}
		s.Seed(seed)
		r := rand.New(s)
		r.Perm(100)
		return r.Perm(100)
	}, func(seed int64) []int {
		s1 := &hashrand.Source{}
		s1.Seed(seed)
		r1 := rand.New(s1)
		r1.Perm(100)
		b, err := s1.MarshalBinary()
		if err != nil {
			panic(err)
		}
		s2 := &hashrand.Source{}
		err = s2.UnmarshalBinary(b)
		if err != nil {
			panic(err)
		}
		r2 := rand.New(s2)
		return r2.Perm(100)
	}, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestReset(t *testing.T) {
	s1 := &hashrand.Source{}
	s2 := &hashrand.Source{}

	i1 := s1.Int63()
	i2 := s2.Int63()
	if i1 != i2 {
		t.Errorf("expected equal (1): %d, %d", i1, i2)
	}

	s1.Int63()

	i1 = s1.Int63()
	i2 = s2.Int63()
	if i1 == i2 {
		t.Errorf("expected not equal (2): %d, %d", i1, i2)
	}

	s1.Reset()
	s3 := &hashrand.Source{}

	i1 = s1.Int63()
	i3 := s3.Int63()

	if i1 != i3 {
		t.Errorf("expected equal (3): %d, %d", i1, i3)
	}
}
