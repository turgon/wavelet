package main

import (
	"fmt"
	"github.com/willf/bitset"
	"sort"
	"strings"
)

type WaveletTree struct {
	alphabet []string
	left     *WaveletTree
	right    *WaveletTree
	data     *bitset.BitSet
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
		alphabet: ab,
		left:     lwt,
		right:    rwt,
		data:     &d,
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

func (wt *WaveletTree) Iter() chan string {

	counters := make(map[*WaveletTree]uint)
	ch := make(chan string, wt.data.Len())

	go func() {
		for i:= uint(0); i < wt.data.Len(); i++ {
			wt.iterate(counters, ch)
		}
		close(ch)
	}()

	return ch
}

func (wt *WaveletTree) iterate(m map[*WaveletTree]uint, ch chan string) {

	if len(wt.alphabet) > 2 {
		if wt.data.Test(m[wt]) {
			wt.right.iterate(m, ch)
		} else {
			if wt.left == nil {
				ch <- wt.alphabet[0]
			} else {
				wt.left.iterate(m, ch)
			}
		}
	} else {
		if wt.data.Test(m[wt]) {
			ch <- wt.alphabet[1]
		} else {
			ch <- wt.alphabet[0]
		}
	}

	m[wt]++
}

func main() {
	var str string = "Sing, Goddess, of the wrath of Achilles"

	ab := alphabet(str)
	chars := strings.Split(str, "")

	wt := NewWaveletTree(ab, chars)

	var newstr string

	for x := range wt.Iter() {
		newstr += x
	}
	fmt.Println(newstr)
}
