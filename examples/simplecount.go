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
