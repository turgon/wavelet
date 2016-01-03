// Package wavelet provides the Wavelet Tree, a succinct data
// structure. Wavelet Trees allow fast Rank and Select operations
// against strings with arbitrarily large alphabets by converting the
// input string into a form of binary tree.
package wavelet

import (
	"github.com/willf/bitset"
)

type WaveletTree struct {
	leftAlphabet  []string
	rightAlphabet []string
	left          *WaveletTree
	right         *WaveletTree
	data          *bitset.BitSet
}

// NewWaveletTree returns a pointer to a new WaveletTree given an
// alphabet and string to encode. Both the alphabet and input
// should be slices of strings. The alphabet should be a slice of
// distinct symbols used in the input text, and the input text
// should be an in-order sequence of the alphabet's symbols.
//
// Note: the tree will ignore any symbols in the input that are not
// found in the given alphabet.
func NewWaveletTree(ab []string, s []string) *WaveletTree {

	var left_s, right_s []string
	var lwt, rwt *WaveletTree

	left_ab := ab[:len(ab)/2]
	right_ab := ab[len(ab)/2:]

	var d bitset.BitSet
	var ctr uint = 0

	for _, x := range s {
		switch {
		case inSlice(x, left_ab):
			left_s = append(left_s, x)
			d.Set(ctr)
			d.Clear(ctr)
			ctr++
		case inSlice(x, right_ab):
			right_s = append(right_s, x)
			d.Set(ctr)
			ctr++
		}
	}

	if len(left_ab) > 1 {
		lwt = NewWaveletTree(left_ab, left_s)
	}

	if len(right_ab) > 1 {
		rwt = NewWaveletTree(right_ab, right_s)
	}

	var wt WaveletTree = WaveletTree{
		leftAlphabet:  left_ab,
		rightAlphabet: right_ab,
		left:          lwt,
		right:         rwt,
		data:          &d,
	}

	return &wt
}

// inSlice returns true/false depending on whether the symbol x is
// in the set of symbols l.
func inSlice(x string, l []string) bool {
	for _, y := range l {
		if x == y {
			return true
		}
	}
	return false
}

// Rank returns the number of occurrences of the symbol q within the
// first x symbols.
//
// Rank is a sort of anti-Select operation; given a symbol q and a
// position boundary, Rank tells you how many times the symbol appears.
//
// It works by counting up the set bits in the WT's data, then using
// that count to recurse into the appropriate child WT. The count at
// the leaf node is the final answer.
func (wt *WaveletTree) Rank(x uint, q string) uint {
	var tot uint

	for i := uint(0); i < x; i++ {
		if wt.data.Test(i) {
			tot++
		}
	}

	if inSlice(q, wt.leftAlphabet) {
		tot = x - tot

		if nil != wt.left {
			return wt.left.Rank(tot, q)
		}
	} else {
		if nil != wt.right {
			return wt.right.Rank(tot, q)
		}
	}

	return tot
}

// Select finds the xth occurrence of symbol q and returns its
// location.
//
// Select is a sort of anti-Rank operation; given a symbol and a rank,
// Select finds the position of that rank.
//
// It works by recursing down to the right leaf node, then using the
// location of the xth occurrence of symbol q there to limit search
// in its parent WT.
func (wt *WaveletTree) Select(x uint, q string) uint {

	var isLeft bool = false

	if inSlice(q, wt.leftAlphabet) {
		isLeft = true
		if nil != wt.left {
			x = wt.left.Select(x, q)
		}
	} else {
		isLeft = false
		if nil != wt.right {
			x = wt.right.Select(x, q)
		}
	}

	// find the local position of the desired q string,
	// then use it to replace x in the parent call.
	// when the root node is hit, x is now the overall position.

	var ctr, pos uint

	for pos = 0; pos < wt.data.Len() && ctr < x; pos++ {

		bitActive := wt.data.Test(pos)

		if isLeft && !bitActive {
			ctr++
		}
		if !isLeft && bitActive {
			ctr++
		}
	}

	return pos
}

// Iter returns an out-channel that will contain symbols from the
// input text in order. It's useful in conjunction with Go's range
// operation. You can think of this as a way to "decompress" the input
// back into slice-of-string form.
func (wt *WaveletTree) Iter() <-chan string {

	counters := make(map[*WaveletTree]uint)
	ch := make(chan string, wt.data.Len())

	go func() {
		for i := uint(0); i < wt.data.Len(); i++ {
			wt.iterate(counters, ch)
		}
		close(ch)
	}()

	return ch
}

// iterate is a helper function for use with Iter(). It walks through
// the symbols in the root WT and recurses into the child nodes,
// pulling symbols from the leaves.
func (wt *WaveletTree) iterate(m map[*WaveletTree]uint, ch chan string) {

	if wt.data.Test(m[wt]) {
		if nil != wt.right {
			wt.right.iterate(m, ch)
		} else {
			ch <- wt.rightAlphabet[0]
		}
	} else {
		if nil != wt.left {
			wt.left.iterate(m, ch)
		} else {
			ch <- wt.leftAlphabet[0]
		}
	}

	m[wt]++
}

// Size returns the number of bits of the bitsets used to encode the
// input text. It does not include pointer overhead or the slices of
// alphabet symbols.
func (wt *WaveletTree) Size() uint {
	var l, r uint

	if wt.left != nil {
		l = wt.left.Size()
	}

	if wt.right != nil {
		r = wt.right.Size()
	}

	return wt.data.Len() + l + r
}
