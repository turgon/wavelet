// +build gofuzz

package rrr

import (
	"fmt"

	"github.com/turgon/wavelet/bitfield"
)

func Fuzz(data []byte) int {
	var cnt uint64

	bf := bitfield.NewBitField(uint(len(data) * 8))

	for j, b := range data {
		for i := uint(0); i < 8; i++ {
			if ((b << i) & 1) > 0 {
				bf.Set(uint(j)*8 + i)
				cnt++
			}
		}
	}

	r := NewRRR(&bf)

	rank := r.Rank(uint(len(data)*8))

	if rank != cnt {
		panic(fmt.Sprintf("Calculated rank is wrong! should be %v, was %v\n", cnt, rank))
	}

	return 1
}
