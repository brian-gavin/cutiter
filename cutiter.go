// package cutiter contains iterating over strings like strings.Split with zero allocation.
// Iter should be used with a for loop as:
//
//	var it cutiter.Iter
//	for k, ok := it.Start("a.b.c.d", "."); ok; k, ok := it.Advance() {
//		fmt.Print(k)
//	}
//
// Outputs:
//
//	abcd
//
// It can additionally advance mid-iteration to skip elements.
//
//	var it cutiter.Iter
//	for k, ok := it.Start("a.0.b.c.d", "."); ok; k, ok := it.Advance() {
//		if k == "a" {
//			k, ok = it.Advance()
//			continue
//		}
//		fmt.Print(k)
//	}
//
// Outputs:
//
//	abcd
//
// Strings without the separator will have 1 iteration, as expected:
//
//	var it cutiter.Iter
//	for k, ok := it.Start("hello", "."); ok; k, ok := it.Advance() {
//		fmt.Print(k)
//	}
//
// Outputs:
//
//	hello
package cutiter

import "strings"

// Iter iterates over an input string using strings.Cut to avoid allocation.
// It should be used with a for loop as:
//
//	var it cutiter.Iter
//	for k, ok := it.Start("a.b.c.d", "."); ok; k, ok := it.Advance() {
//		fmt.Print(k)
//	}
//
// If the string given to s does not contain the separator, Start(s,sep) will return (s, true).
// The subsequent call to Advance() will return ("", false).
// If s is empty, then Start() will return ("", false).
type Iter struct {
	sep  string
	rest string
}

// Start initializes the iteration, and returns the first key.
// If s does not contain sep, it will return (s, true).
// The subsequent call to Advance() will return ("", false).
func (it *Iter) Start(s, sep string) (string, bool) {
	it.sep = sep
	return it.next(s)
}

// Advance advances the iteration to the next key.
// When iteration is completed, it returns ("", false)
func (it *Iter) Advance() (string, bool) {
	return it.next(it.rest)
}

// next is the internal method used by iter.Start and iter.Advance. it will cut s,
// and return the "before", storing the "after" for the subsequent call.
// if path does not contain "." and is non-empty, it will return (path, true), storing "".
// if path is empty, it will return ("", false). this means the end of iteration.
func (it *Iter) next(s string) (key string, ok bool) {
	key, it.rest, _ = strings.Cut(s, it.sep)
	ok = key != ""
	return key, ok
}
