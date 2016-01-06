# Wavelet Tree

Hello!

I've implemented a [Wavelet Tree][1] data structure. For a great introduction to Wavelet Trees, see Alex Bowe's excellent blog post, [Wavelet Trees - an Introduction][2]

## Status

The Wavelet Tree structure is fully functional and well-tested, but very no-frills. I'd like to spend time implementing the RRR structure as well, which can provide O(1) Rank and Select operations. The current implementation of Rank and Select in the Wavelet Tree is O(s log s * n log n) for an n-symbol input with alphabet size of s. Pretty bad! Moving to RRR should make these operations closer to O(log n).

## Examples

There are code examples in the [examples directory][3], but here's a short one:

```
package main

import (
	"fmt"

	"github.com/turgon/wavelet"
)

func main() {

	var alphabet []string = []string{"a", "b", "c"}
	var symbols []string = []string{"a", "a", "b", "a", "b", "c", "a"}

	wt := wavelet.NewWaveletTree(alphabet, symbols)

	// Use Rank to count the number of occurrences of each
	// symbol in the alphabet
	for _, c := range alphabet {
		fmt.Printf("%s: %v times\n", c, wt.Rank(uint(len(symbols)), c))
	}
}
```

## Roadmap

* Use go-fuzz to do fuzzy testing
* Implement RRR


[1]: https://en.wikipedia.org/wiki/Wavelet_Tree
[2]: http://alexbowe.com/wavelet-trees/
[3]: examples
