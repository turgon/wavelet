// Package wavelet provides the Wavelet Tree, a succinct data
// structure. Wavelet Trees allow fast Rank and Select operations
// against strings with arbitrarily large alphabets by converting the
// input string into a form of binary tree.
package wavelet

import (
	"github.com/willf/bitset"
	"sort"
	"strings"
)

type WaveletTree struct {
	leftAlphabet  []string
	rightAlphabet []string
	left          *WaveletTree
	right         *WaveletTree
	data          *bitset.BitSet
}

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

func inSlice(x string, l []string) bool {
	for _, y := range l {
		if x == y {
			return true
		}
	}
	return false
}

func alphabet(s string) []string {
	var r []string
	var chars []string = strings.Split(s, "")

	sort.Strings(chars)

	r = append(r, chars[0])

	for _, x := range chars {
		if x != r[len(r)-1] {
			r = append(r, x)
		}
	}
	return r
}

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

// Select finds the xth occurrence of string q and returns its
// location.
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

func (wt *WaveletTree) Iter() chan string {

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
