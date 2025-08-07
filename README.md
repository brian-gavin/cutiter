# cutiter

A Go package for iterating over strings using `strings.Cut`, as an alternative to `strings.Split`(when `strings.SplitSeq` doesn't get the job done).

The iterator(s[^1]) here are intended to be 0 allocation and extremely efficient.

```go
package main

import (
	"fmt"

	"github.com/brian-gavin/cutiter"
)

func main() {
	var it cutiter.Iter
	for k, ok := it.Start("hello.world", "."); ok; k, ok = it.Advance() {
		fmt.Printf("%s\n", k)
	}
}
```

## use case

1. You don't need to retain the result of `strings.Split` (this is also a usage of `strings.SplitSeq`)
2. You need to lookahead or advance mid-way through your iteration
3. But don't want to use a Pull iterator because of the allocations it does.

```go
// checkValidFormat checks that s is dotted notation with each element being <letter>+ OR a.<number>
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
```

## benchmark results

Benchmarks are done to compare `cutiter.Iter` to the performance of `strings.Split`, `strings.SplitSeq`, and `iter.Pull(strings.SplitSeq)`.
For most usual cases, `cutiter.Iter` should see similar performance to `strings.SplitSeq`, and improved performance over `strings.Split`.

Some edge cases that `strings.Cut` does not naturally handle have a performance hit. These cases are implemented for full `strings.Split` compatability. These benchmarks are also displayed to show the tradeoffs that `cutiter.Iter` makes. These cases may have improved performance down the line.

```
goos: darwin
goarch: amd64
pkg: github.com/brian-gavin/cutiter
cpu: Intel(R) Core(TM) i9-9980HK CPU @ 2.40GHz
BenchmarkIter
BenchmarkIter/empty
    cutiter_test.go:98: elements: 0 | elementLen: 0 | sep: "." | len(input): 0
BenchmarkIter/empty/cutiter.Iter
BenchmarkIter/empty/cutiter.Iter-16         	59478816	        18.62 ns/op	       0 B/op	       0 allocs/op
BenchmarkIter/empty/strings.Split
BenchmarkIter/empty/strings.Split-16        	36525802	        32.54 ns/op	      16 B/op	       1 allocs/op
BenchmarkIter/empty/strings.SplitSeq
BenchmarkIter/empty/strings.SplitSeq-16     	11807013	       101.8 ns/op	      88 B/op	       4 allocs/op
BenchmarkIter/empty/iter.Pull(strings.SplitSeq)
BenchmarkIter/empty/iter.Pull(strings.SplitSeq)-16         	 2210746	       545.3 ns/op	     416 B/op	      14 allocs/op

BenchmarkIter/short
    cutiter_test.go:98: elements: 52 | elementLen: 1 | sep: "." | len(input): 103
BenchmarkIter/short/cutiter.Iter
BenchmarkIter/short/cutiter.Iter-16                        	 2179500	       549.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkIter/short/strings.Split
BenchmarkIter/short/strings.Split-16                       	 2110801	       566.4 ns/op	     896 B/op	       1 allocs/op
BenchmarkIter/short/strings.SplitSeq
BenchmarkIter/short/strings.SplitSeq-16                    	 2047783	       591.2 ns/op	      88 B/op	       4 allocs/op
BenchmarkIter/short/iter.Pull(strings.SplitSeq)
BenchmarkIter/short/iter.Pull(strings.SplitSeq)-16         	  258468	      4500 ns/op	     416 B/op	      14 allocs/op

BenchmarkIter/long
    cutiter_test.go:98: elements: 1024 | elementLen: 256 | sep: "." | len(input): 263167
BenchmarkIter/long/cutiter.Iter
BenchmarkIter/long/cutiter.Iter-16                         	   66387	     18224 ns/op	       0 B/op	       0 allocs/op
BenchmarkIter/long/strings.Split
BenchmarkIter/long/strings.Split-16                        	   49908	     23568 ns/op	   18432 B/op	       1 allocs/op
BenchmarkIter/long/strings.SplitSeq
BenchmarkIter/long/strings.SplitSeq-16                     	   64971	     18184 ns/op	      88 B/op	       4 allocs/op
BenchmarkIter/long/iter.Pull(strings.SplitSeq)
BenchmarkIter/long/iter.Pull(strings.SplitSeq)-16          	   15163	     78977 ns/op	     416 B/op	      14 allocs/op

BenchmarkIter/twoLongElements
    cutiter_test.go:98: elements: 2 | elementLen: 65536 | sep: "." | len(input): 131073
BenchmarkIter/twoLongElements/cutiter.Iter
BenchmarkIter/twoLongElements/cutiter.Iter-16              	  444014	      2760 ns/op	       0 B/op	       0 allocs/op
BenchmarkIter/twoLongElements/strings.Split
BenchmarkIter/twoLongElements/strings.Split-16             	  354555	      3255 ns/op	      32 B/op	       1 allocs/op
BenchmarkIter/twoLongElements/strings.SplitSeq
BenchmarkIter/twoLongElements/strings.SplitSeq-16          	  414417	      2887 ns/op	      88 B/op	       4 allocs/op
BenchmarkIter/twoLongElements/iter.Pull(strings.SplitSeq)
BenchmarkIter/twoLongElements/iter.Pull(strings.SplitSeq)-16         	  344030	      3512 ns/op	     416 B/op	      14 allocs/op

BenchmarkIter/superLong
    cutiter_test.go:98: elements: 2048 | elementLen: 65536 | sep: "........" | len(input): 134234104
BenchmarkIter/superLong/cutiter.Iter
BenchmarkIter/superLong/cutiter.Iter-16                              	     136	   8759636 ns/op	       0 B/op	       0 allocs/op
BenchmarkIter/superLong/strings.Split
BenchmarkIter/superLong/strings.Split-16                             	      66	  17166985 ns/op	   32768 B/op	       1 allocs/op
BenchmarkIter/superLong/strings.SplitSeq
BenchmarkIter/superLong/strings.SplitSeq-16                          	     142	   8362771 ns/op	      88 B/op	       4 allocs/op
BenchmarkIter/superLong/iter.Pull(strings.SplitSeq)
BenchmarkIter/superLong/iter.Pull(strings.SplitSeq)-16               	     138	   8657333 ns/op	     419 B/op	      14 allocs/op

BenchmarkIter/shortEmptySeparator
    cutiter_test.go:98: elements: 52 | elementLen: 1 | sep: "" | len(input): 52
BenchmarkIter/shortEmptySeparator/cutiter.Iter
BenchmarkIter/shortEmptySeparator/cutiter.Iter-16                    	  225180	      5319 ns/op	     416 B/op	      16 allocs/op
BenchmarkIter/shortEmptySeparator/strings.Split
BenchmarkIter/shortEmptySeparator/strings.Split-16                   	 3143425	       376.4 ns/op	     896 B/op	       1 allocs/op
BenchmarkIter/shortEmptySeparator/strings.SplitSeq
BenchmarkIter/shortEmptySeparator/strings.SplitSeq-16                	 3430003	       353.0 ns/op	      72 B/op	       5 allocs/op
BenchmarkIter/shortEmptySeparator/iter.Pull(strings.SplitSeq)
BenchmarkIter/shortEmptySeparator/iter.Pull(strings.SplitSeq)-16     	  304188	      4067 ns/op	     400 B/op	      15 allocs/op

BenchmarkIter/emptyStringEmptySeparator
    cutiter_test.go:98: elements: 0 | elementLen: 0 | sep: "" | len(input): 0
BenchmarkIter/emptyStringEmptySeparator/cutiter.Iter
BenchmarkIter/emptyStringEmptySeparator/cutiter.Iter-16              	 2314797	       516.1 ns/op	     416 B/op	      16 allocs/op
BenchmarkIter/emptyStringEmptySeparator/strings.Split
BenchmarkIter/emptyStringEmptySeparator/strings.Split-16             	128797786	         9.303 ns/op	       0 B/op	       0 allocs/op
BenchmarkIter/emptyStringEmptySeparator/strings.SplitSeq
BenchmarkIter/emptyStringEmptySeparator/strings.SplitSeq-16          	11182239	       105.7 ns/op	      72 B/op	       5 allocs/op
BenchmarkIter/emptyStringEmptySeparator/iter.Pull(strings.SplitSeq)
BenchmarkIter/emptyStringEmptySeparator/iter.Pull(strings.SplitSeq)-16         	 2545087	       472.2 ns/op	     400 B/op	      15 allocs/op
```

[^1]: coming soon?
