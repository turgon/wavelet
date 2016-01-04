package main

import (
	"fmt"

	"github.com/turgon/wavelet"
)

func main() {

	var alphabet []string = []string{"a", "b", "c"}
	var symbols []string = []string{"a", "a", "b", "a", "b", "c", "a"}

	wt := wavelet.NewWaveletTree(alphabet, symbols)

	// Use Select to find the first occurrence of each
	// symbol in the alphabet; the largest can be used to
	// find the shortest prefix of symbols that includes
	// every symbol.
	var max uint
	for _, c := range alphabet {
		x := wt.Select(1, c)
		if x > max {
			max = x
		}
	}
	fmt.Printf("shortest prefix is: %v\n", symbols[:max])
}
