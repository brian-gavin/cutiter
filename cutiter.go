// package cutiter contains types for iterating over strings using strings.Cut,
// to replace strings.Split usage when strings.SplitSeq doesn't get the job done.
package cutiter

import "strings"

type Iter struct {
	sep  string
	rest string
}

// Start initializes the iteration, and returns the first key.
func (it *Iter) Start(s, sep string) (string, bool) {
	it.sep = sep
	return it.next(s)
}

// Advance advances the iteration to the next key.
func (it *Iter) Advance() (string, bool) {
	return it.next(it.rest)
}

// next should only be used by (*pathIter).start and (*pathIter).advance. it will cut path,
// and return the "before", storing the "after" for the subsequent call.
// if path does not contain "." and is non-empty, it will return (path, true), storing "".
// if path is empty, it will return ("", false). this means the end of iteration.
func (it *Iter) next(path string) (key string, ok bool) {
	key, it.rest, _ = strings.Cut(path, it.sep)
	ok = key != ""
	return key, ok
}
