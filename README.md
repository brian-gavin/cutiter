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

## benchmark results

Benchmarks are done to compare `cutiter.Iter` to the performance of `strings.Split` and `strings.SplitSeq`.

```
goarch: amd64
pkg: github.com/brian.gavin/cutiter
cpu: Intel(R) Core(TM) i9-9980HK CPU @ 2.40GHz
BenchmarkIter
BenchmarkIter/empty
    cutiter_test.go:81: elements: 0 | elementLen: 0 | sep: "" | len(input): 0
BenchmarkIter/empty/cutiter.Iter
BenchmarkIter/empty/cutiter.Iter-16         	150675615	         8.588 ns/op	       0 B/op	       0 allocs/op
BenchmarkIter/empty/strings.Split
BenchmarkIter/empty/strings.Split-16        	81887461	        13.89 ns/op	       0 B/op	       0 allocs/op
BenchmarkIter/empty/strings.SplitSeq
BenchmarkIter/empty/strings.SplitSeq-16     	 9480236	       109.2 ns/op	      72 B/op	       5 allocs/op
BenchmarkIter/short
    cutiter_test.go:81: elements: 52 | elementLen: 1 | sep: "." | len(input): 103
BenchmarkIter/short/cutiter.Iter
BenchmarkIter/short/cutiter.Iter-16         	 2138184	       550.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkIter/short/strings.Split
BenchmarkIter/short/strings.Split-16        	 2065755	       572.8 ns/op	     896 B/op	       1 allocs/op
BenchmarkIter/short/strings.SplitSeq
BenchmarkIter/short/strings.SplitSeq-16     	 2033740	       597.8 ns/op	      88 B/op	       4 allocs/op
BenchmarkIter/long
    cutiter_test.go:81: elements: 1024 | elementLen: 256 | sep: "." | len(input): 263167
BenchmarkIter/long/cutiter.Iter
BenchmarkIter/long/cutiter.Iter-16          	   65427	     17846 ns/op	       0 B/op	       0 allocs/op
BenchmarkIter/long/strings.Split
BenchmarkIter/long/strings.Split-16         	   49592	     23335 ns/op	   18432 B/op	       1 allocs/op
BenchmarkIter/long/strings.SplitSeq
BenchmarkIter/long/strings.SplitSeq-16      	   61321	     19328 ns/op	      88 B/op	       4 allocs/op
BenchmarkIter/twoLongElements
    cutiter_test.go:81: elements: 2 | elementLen: 65536 | sep: "." | len(input): 131073
BenchmarkIter/twoLongElements/cutiter.Iter
BenchmarkIter/twoLongElements/cutiter.Iter-16         	  415492	      2818 ns/op	       0 B/op	       0 allocs/op
BenchmarkIter/twoLongElements/strings.Split
BenchmarkIter/twoLongElements/strings.Split-16        	  329936	      3321 ns/op	      32 B/op	       1 allocs/op
BenchmarkIter/twoLongElements/strings.SplitSeq
BenchmarkIter/twoLongElements/strings.SplitSeq-16     	  385611	      2947 ns/op	      88 B/op	       4 allocs/op
BenchmarkIter/superLong
    cutiter_test.go:81: elements: 2048 | elementLen: 65536 | sep: "........" | len(input): 134234104
BenchmarkIter/superLong/cutiter.Iter
BenchmarkIter/superLong/cutiter.Iter-16               	     129	   9103489 ns/op	       0 B/op	       0 allocs/op
BenchmarkIter/superLong/strings.Split
BenchmarkIter/superLong/strings.Split-16              	      62	  17712717 ns/op	   32768 B/op	       1 allocs/op
BenchmarkIter/superLong/strings.SplitSeq
BenchmarkIter/superLong/strings.SplitSeq-16           	     135	   8830863 ns/op	      88 B/op	       4 allocs/op
PASS
ok  	github.com/brian.gavin/cutiter	18.155s
```
