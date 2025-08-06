package cutiter_test

import (
	"strconv"
	"testing"
	"unicode"

	"github.com/brian.gavin/cutiter"
	"github.com/stretchr/testify/assert"
)

// checkValidFormat checks that s is dotted notation with each element being <letter>+ OR a.<number>
// this is in the README so we should test it's correct.
func checkValidFormat(s string) bool {
	var it cutiter.Iter
	for k, ok := it.Start(s, "."); ok; k, ok = it.Advance() {
		// if a, next must be a number
		if k == "a" {
			k, ok = it.Advance()
			if !ok {
				return false
			}
			if _, err := strconv.Atoi(k); err != nil {
				return false
			}
			continue
		}
		// otherwise, all runes of k must be letters
		for _, r := range k {
			if !unicode.IsLetter(r) {
				return false
			}
		}
	}
	return true
}

type readmetestcase struct {
	input string
	valid bool
}

var readmetcs = []readmetestcase{
	{"", true},
	{"b", true},
	{"bb", true},
	{"a.0", true},
	{"a", false},
	{"a.0.bbb", true},
	{"a.bb", false},
	{"bb.a.0.cc", true},
	{"bb.a.cc", false},
}

func TestREADMESample(t *testing.T) {
	for _, tc := range readmetcs {
		t.Run(tc.input, func(t *testing.T) {
			assert.Equal(t, tc.valid, checkValidFormat(tc.input))
		})
	}
}

func BenchmarkREADMESample(b *testing.B) {
	for _, tc := range readmetcs {
		b.Run(tc.input, func(b *testing.B) {
			for b.Loop() {
				_ = checkValidFormat(tc.input)
			}
		})
	}
}
