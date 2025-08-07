package cutiter_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"unicode"

	"github.com/brian-gavin/cutiter"
	"github.com/stretchr/testify/assert"
)

var (
	errMissingNumber = errors.New("missing number after 'a'")
	errNotANumber    = errors.New("not a number")
	errNotALetter    = errors.New("not a letter")
)

// checkValidFormat checks that s is dotted notation with each element being <letter>+ OR a.<number>
// this is in the README so we should test it's correct.
func checkValidFormat(s string) error {
	var it cutiter.Iter
	for k, ok := it.Start(s, "."); ok; k, ok = it.Advance() {
		// if a, next must be a number
		if k == "a" {
			k, ok = it.Advance()
			if !ok {
				return errMissingNumber
			}
			if _, err := strconv.Atoi(k); err != nil {
				return fmt.Errorf("%s: %w", k, errors.Join(errNotANumber, err))
			}
			continue
		}
		// otherwise, all runes of k must be letters
		for _, r := range k {
			if !unicode.IsLetter(r) {
				return fmt.Errorf("%c: %w", r, errNotALetter)
			}
		}
	}
	return nil
}

type readmetestcase struct {
	input  string
	expErr error
}

var readmetcs = []readmetestcase{
	{"", nil},
	{"b", nil},
	{"bb", nil},
	{"a.0", nil},
	{"a", errMissingNumber},
	{"a.0.bbb", nil},
	{"a.bb", errNotANumber},
	{"bb.a.0.cc", nil},
	{"bb.a.cc", errNotANumber},
	{"b.0", errNotALetter},
}

func TestREADMESample(t *testing.T) {
	for _, tc := range readmetcs {
		t.Run(tc.input, func(t *testing.T) {
			assert.ErrorIs(t, checkValidFormat(tc.input), tc.expErr)
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
