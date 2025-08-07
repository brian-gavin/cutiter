package cutiter_test

import (
	"iter"
	"strings"
	"testing"

	"github.com/brian-gavin/cutiter"
	"github.com/stretchr/testify/assert"
)

type splittestcase struct {
	name string
	s    string
	sep  string
}

func (tc *splittestcase) run(t *testing.T) {
	var (
		a       = assert.New(t)
		split   = strings.Split(tc.s, tc.sep)
		collect = make([]string, 0, len(split))
		it      cutiter.Iter
	)
	for k, ok := it.Start(tc.s, tc.sep); ok; k, ok = it.Advance() {
		collect = append(collect, k)
	}
	a.Equal(split, collect, "do not have the same result as strings.Split: %q", split)
}

func TestSplitEquivalence(t *testing.T) {
	tcs := []splittestcase{
		{"empty", "", "."},
		{"emptySepNonemptyString", "a", ""},
		{"singleLetter", "a", "."},
		{"normal", "a.b.c.d.e.f.g", "."},
		{"leadingSep", ".a", "."},
		{"trailingSep", "a.", "."},
		{"onlySep", ".", "."},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.run(t)
		})
	}
}

func TestAdvanceMidIteration(t *testing.T) {
	var (
		a     = assert.New(t)
		sep   = "."
		order = []string{"a", "b", "c", "d", "e", "f", "g"}
		s     = "a.0.b.c.d.e.f.g"
		i     int
		it    cutiter.Iter
	)
	for k, ok := it.Start(s, sep); ok; k, ok = it.Advance() {
		a.Equal(order[i], k, "expected %d element of order: %v", i, order)
		i++
		// skip "0"
		if k == "a" {
			k, ok = it.Advance()
			a.True(ok, "advancement should be OK")
			a.Equal("0", k, "advancing when at 'a' should result in 'b'")
		}
	}
}

type benchcase struct {
	name       string
	elements   int
	elementLen int
	sep        string
}

func (bc benchcase) genInput() string {
	ss := make([]string, bc.elements)
	for i := range bc.elements {
		ss[i] = strings.Repeat(string(rune('a'+(i%26))), bc.elementLen)
	}
	return strings.Join(ss[:], bc.sep)
}

var benchCases = []benchcase{
	{"empty", 0, 0, "."},
	{"short", 52, 1, "."},
	{"long", 1024, 256, "."},
	{"twoLongElements", 2, 1 << 16, "."},
	{"superLong", 2048, 1 << 16, "........"},
	{"shortEmptySeparator", 52, 1, ""},
	{"emptyStringEmptySeparator", 0, 0, ""},
}

func BenchmarkIter(b *testing.B) {
	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			input := bc.genInput()
			b.Logf("elements: %d | elementLen: %d | sep: %q | len(input): %d", bc.elements, bc.elementLen, bc.sep, len(input))
			b.Run("cutiter.Iter", func(b *testing.B) {
				for b.Loop() {
					var it cutiter.Iter
					for k, ok := it.Start(input, bc.sep); ok; k, ok = it.Advance() {
						_ = k
					}
				}
			})
			b.Run("strings.Split", func(b *testing.B) {
				for b.Loop() {
					s := strings.Split(input, bc.sep)
					for _, k := range s {
						_ = k
					}
				}
			})
			b.Run("strings.SplitSeq", func(b *testing.B) {
				for b.Loop() {
					for k := range strings.SplitSeq(input, bc.sep) {
						_ = k
					}
				}
			})
			b.Run("iter.Pull(strings.SplitSeq)", func(b *testing.B) {
				for b.Loop() {
					next, close := iter.Pull(strings.SplitSeq(input, bc.sep))
					for k, ok := next(); ok; k, ok = next() {
						_ = k
					}
					close()
				}
			})
		})
	}
}

// TestBenchmarkCases tests the benchmark cases.
func TestBenchmarkCases(t *testing.T) {
	for _, bc := range benchCases {
		tc := splittestcase{
			name: bc.name,
			s:    bc.genInput(),
			sep:  bc.sep,
		}
		t.Run(tc.name, func(t *testing.T) {
			tc.run(t)
		})
	}
}
