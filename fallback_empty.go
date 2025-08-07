package cutiter

import (
	"iter"
	"strings"
)

// fallbackEmpty is used when the sep is "". strings.Cut does not split the elements on the rune
// like Split does. so we need to resort to a fallback implementation that will not be as efficient
// as Cut.
// Right now, this is based on iter.Pull(iter.SplitSeq(s, ""))
type fallbackEmpty struct {
	itNext func() (string, bool)
	itStop func()
}

func newFallbackEmpty(s string) *fallbackEmpty {
	var f fallbackEmpty
	f.itNext, f.itStop = iter.Pull(strings.SplitSeq(s, ""))
	return &f
}

func (f *fallbackEmpty) next() (string, bool) {
	key, ok := f.itNext()
	if !ok {
		f.itStop()
	}
	return key, ok
}
