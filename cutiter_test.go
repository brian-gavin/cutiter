package cutiter_test

import (
	"strings"
	"testing"

	"github.com/brian.gavin/cutiter"
	"github.com/stretchr/testify/assert"
)

func TestIteration(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var it cutiter.Iter
		for k, ok := it.Start("", "."); ok; k, ok = it.Advance() {
			t.Fatalf("iteration should not occur: k: %q", k)
		}
	})
	t.Run("normal", func(t *testing.T) {
		var (
			a     = assert.New(t)
			sep   = "."
			order = []string{"a", "b", "c", "d", "e", "f", "g"}
			s     = strings.Join(order, sep)
			i     int
			it    cutiter.Iter
		)
		for k, ok := it.Start(s, sep); ok; k, ok = it.Advance() {
			a.Equal(order[i], k, "expected %d element of order: %v", i, order)
			i++
		}
		a.Equal(len(order), i, "did not iterate over every element.")
	})
	t.Run("advanceMidIteration", func(t *testing.T) {
		var (
			a     = assert.New(t)
			sep   = "."
			order = []string{"a", "c", "d", "e", "f", "g"}
			s     = "a.b.c.d.e.f.g"
			i     int
			it    cutiter.Iter
		)
		for k, ok := it.Start(s, sep); ok; k, ok = it.Advance() {
			a.Equal(order[i], k, "expected %d element of order: %v", i, order)
			i++
			// skip "b"
			if k == "a" {
				k, ok = it.Advance()
				a.True(ok, "advancement should be OK")
				a.Equal("b", k, "advancing when at 'a' should result in 'b'")
			}
		}
		a.Equal(len(order), i, "did not iterate over every element.")
	})
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
	{"empty", 0, 0, ""},
	{"short", 52, 1, "."},
	{"long", 1024, 256, "."},
	{"twoLongElements", 2, 1 << 16, "."},
	{"superLong", 2048, 1 << 16, "........"},
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
		})
	}
}

// TestBenchmarkCases tests the benchmark cases.
func TestBenchmarkCases(t *testing.T) {
	for _, bc := range benchCases {
		t.Run(bc.name, func(t *testing.T) {
			var (
				a     = assert.New(t)
				input = bc.genInput()
				it    cutiter.Iter
				i     int
			)
			for k, ok := it.Start(input, bc.sep); ok; k, ok = it.Advance() {
				i++
				a.Len(k, bc.elementLen)
			}
			a.Equal(bc.elements, i, "loop did not iterate over every element")
		})
	}
}
