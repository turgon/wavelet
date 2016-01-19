// Package wavelet provides the Wavelet Tree, a succinct data
// structure. Wavelet Trees allow fast Rank and Select operations
// against strings with arbitrarily large alphabets by converting the
// input string into a form of binary tree.
package wavelet

import (
	"github.com/turgon/wavelet/bitfield"
)

type WaveletTree struct {
	leftAlphabet  []string
	rightAlphabet []string
	left          *WaveletTree
	right         *WaveletTree
	data          bitfield.BitField
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

	var d bitfield.BitField = bitfield.NewBitField(uint(len(s)))
	var ctr uint = 0

	for _, x := range s {
		switch {
		case inSlice(x, left_ab):
			left_s = append(left_s, x)
			ctr++
		case inSlice(x, right_ab):
			right_s = append(right_s, x)
			d.Set(ctr)
			ctr++
		}
	}

	if ctr < uint(len(s)) {
		d = d.Resize(ctr)
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
		data:          d,
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

	if x > wt.data.Len() {
		x = wt.data.Len()
	}

	// this calculation is slow because it's O(n)
	// not to mention that this is the most naive way to popcount
	for i := uint(0); i < x && i < wt.data.Len(); i++ {
		if wt.data.Test(i) {
			tot++
		}
	}

	// this recursion is slow because inSlice is O(n)
	if inSlice(q, wt.leftAlphabet) {
		tot = x - tot
		if nil != wt.left {
			return wt.left.Rank(tot, q)
		}
	} else if inSlice(q, wt.rightAlphabet) {
		if nil != wt.right {
			return wt.right.Rank(tot, q)
		}
	}

	return tot
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
