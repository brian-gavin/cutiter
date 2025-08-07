# cutiter

A Go package for iterating over strings using `strings.Cut`, as an alternative to `strings.Split` when `strings.SplitSeq` doesn't get the job done.

The iterator(s) here are intended to be 0 allocation and extremely efficient.

## use-case

1. When you don't need to retain the result of `strings.Split` (this is also a usage of `strings.SplitSeq`)
2. You may need to lookahead or advance mid-way through your iteration, such as:

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
3. You also don't want to use a Pull iterator because of the allocations it does.

## benchmark results

Benchmarks are done to compare `cutiter.Iter` to the performance of `strings.Split`, `strings.SplitSeq`, and `iter.Pull(strings.SplitSeq)`.

```
goos: darwin
goarch: amd64
pkg: github.com/brian.gavin/cutiter
cpu: Intel(R) Core(TM) i9-9980HK CPU @ 2.40GHz
BenchmarkIter/empty/cutiter.Iter-16         	160384808	         7.464 ns/op	       0 B/op	       0 allocs/op
BenchmarkIter/empty/strings.Split-16        	122008760	         9.836 ns/op	       0 B/op	       0 allocs/op
BenchmarkIter/empty/strings.SplitSeq-16     	10861251	       108.0 ns/op	      72 B/op	       5 allocs/op
BenchmarkIter/empty/iter.Pull(strings.SplitSeq)-16         	 1603045	       794.1 ns/op	     448 B/op	      16 allocs/op

BenchmarkIter/short/cutiter.Iter-16                        	 2048434	       564.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkIter/short/strings.Split-16                       	 2079636	       573.9 ns/op	     896 B/op	       1 allocs/op
BenchmarkIter/short/strings.SplitSeq-16                    	 2040880	       588.5 ns/op	      88 B/op	       4 allocs/op
BenchmarkIter/short/iter.Pull(strings.SplitSeq)-16         	  239355	      4803 ns/op	     464 B/op	      15 allocs/op

BenchmarkIter/long/cutiter.Iter-16                         	   61976	     17920 ns/op	       0 B/op	       0 allocs/op
BenchmarkIter/long/strings.Split-16                        	   48513	     23973 ns/op	   18432 B/op	       1 allocs/op
BenchmarkIter/long/strings.SplitSeq-16                     	   67183	     17900 ns/op	      88 B/op	       4 allocs/op
BenchmarkIter/long/iter.Pull(strings.SplitSeq)-16          	   14235	     84755 ns/op	     463 B/op	      14 allocs/op

BenchmarkIter/twoLongElements/cutiter.Iter-16              	  405403	      2800 ns/op	       0 B/op	       0 allocs/op
BenchmarkIter/twoLongElements/strings.Split-16             	  337844	      3408 ns/op	      32 B/op	       1 allocs/op
BenchmarkIter/twoLongElements/strings.SplitSeq-16          	  376053	      2962 ns/op	      88 B/op	       4 allocs/op
BenchmarkIter/twoLongElements/iter.Pull(strings.SplitSeq)-16         	  284551	      3836 ns/op	     463 B/op	      14 allocs/op

BenchmarkIter/superLong/cutiter.Iter-16                              	     134	   8859670 ns/op	       0 B/op	       0 allocs/op
BenchmarkIter/superLong/strings.Split-16                             	      62	  17950772 ns/op	   32768 B/op	       1 allocs/op
BenchmarkIter/superLong/strings.SplitSeq-16                          	     134	   8907881 ns/op	      88 B/op	       4 allocs/op
BenchmarkIter/superLong/iter.Pull(strings.SplitSeq)-16               	     130	   9107284 ns/op	     448 B/op	      14 allocs/op
```
