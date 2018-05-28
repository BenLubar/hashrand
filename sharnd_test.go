package hashrand_test

import (
	"bytes"
	"crypto/sha1"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/BenLubar/hashrand"
)

// With no hash function explicitly assigned, hashrand.Source functions identically to
// http://mattmahoney.net/dc/#sharnd
func TestSHARND(t *testing.T) {
	input := []struct {
		URL  string
		Name string
		Hash [20]byte
	}{
		{
			URL:  "http://mattmahoney.net/dc/sharnd.c",
			Name: "sharnd.c",
			Hash: [...]byte{0xf9, 0xef, 0x09, 0xac, 0x82, 0x35, 0xc4, 0x5f, 0x87, 0xea, 0xde, 0xd5, 0x36, 0x54, 0xa4, 0xf4, 0x25, 0xfb, 0x05, 0x92},
		},
		{
			URL:  "http://mattmahoney.net/dc/sha1.c",
			Name: "sha1.c",
			Hash: [...]byte{0x6b, 0xe5, 0x31, 0xd6, 0xf6, 0x7e, 0xc2, 0x54, 0x53, 0x1c, 0xaa, 0x3f, 0xb9, 0xcc, 0x45, 0x7d, 0x26, 0x8d, 0x56, 0x8f},
		},
		{
			URL:  "http://mattmahoney.net/dc/sha1.h",
			Name: "sha1.h",
			Hash: [...]byte{0x44, 0x95, 0x6e, 0x91, 0x6d, 0xaa, 0xb8, 0x01, 0x79, 0xd4, 0x32, 0x40, 0x6f, 0xfc, 0xdd, 0xa4, 0x64, 0x8c, 0x9e, 0x07},
		},
	}

	tempDir, err := ioutil.TempDir("", "hashrand")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if e := os.RemoveAll(tempDir); e != nil {
			t.Error(e)
		}
	}()

	for _, f := range input {
		downloadFile(f.URL, filepath.Join(tempDir, f.Name), f.Hash)
	}

	cc := os.Getenv("CC")
	if cc == "" {
		cc = "cc"
	}

	cmd := exec.Command(cc, "-o", "sharnd", "-O3", "sharnd.c", "sha1.c")
	cmd.Dir = tempDir

	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

	err = quick.CheckEqual(func(key string) []byte {
		cmd = exec.Command(filepath.Join(tempDir, "sharnd"), key, "1000")
		cmd.Dir = tempDir
		cmd.Stdout = ioutil.Discard
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			panic(err)
		}

		b, e := ioutil.ReadFile(filepath.Join(tempDir, "sharnd.out"))
		if e != nil {
			panic(err)
		}

		return b
	}, func(key string) []byte {
		s := &hashrand.Source{}
		s.AppendSeed([]byte(key))

		b := make([]byte, 1000)

		_, err = s.Read(b)
		if err != nil {
			panic(err)
		}

		return b
	}, &quick.Config{
		Values: func(v []reflect.Value, r *rand.Rand) {
			const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
			buf := make([]byte, r.Intn(20))
			for i := range buf {
				buf[i] = characters[r.Intn(len(characters))]
			}
			v[0] = reflect.ValueOf(string(buf))
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func downloadFile(url, path string, hash [20]byte) {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := res.Body.Close(); e != nil {
			panic(err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		panic("failed request to " + url + ": " + res.Status)
	}

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()

	digest := sha1.New()
	_, err = io.Copy(io.MultiWriter(digest, f), res.Body)
	if err != nil {
		panic(err)
	}

	if !bytes.Equal(digest.Sum(nil), hash[:]) {
		panic("hash mismatch for " + url)
	}
}
