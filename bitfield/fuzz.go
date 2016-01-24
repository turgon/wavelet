package bitfield

import "fmt"

func Fuzz(data []byte) int {
	bits := uint(len(data)) * segmentSize

	bf := NewBitField(bits)

	var pc uint64
	for i, b := range data {
		for j := uint(0); j < segmentSize; j++ {
			if (b << j) & 1 > 0 {
				bf.Set(uint(i) * segmentSize + j)
				pc++
			}
		}
	}


	for i, b := range data {
		for j := uint(0); j < segmentSize; j++ {
			if (b << j) & 1 > 0 {
				if !bf.Test(uint(i) * segmentSize + j) {
					panic(fmt.Sprintf("bit %v should have been set!", uint(i)*segmentSize+j))
				}
			} else {
				if bf.Test(uint(i) * segmentSize + j) {
					panic(fmt.Sprintf("bit %v should not have been set!", uint(i)*segmentSize+j))
				}
			}
		}
	}

	if pc != bf.Popcount() {
		panic(fmt.Sprintf("Popcount should have been %v but was %v!", pc, bf.Popcount()))
	}

	bff := bf.Sub(bf.Len()/2, bf.Len())
	bf = bf.Resize(bf.Len()/2)

	if pc < bf.Popcount() {
		panic(fmt.Sprintf("Low Popcount shouldn't be greater than original popcount after downsizing"))
	}

	if pc < bff.Popcount() {
		panic(fmt.Sprintf("High Popcount shouldn't be greater than original popcount after downsizing"))
	}

	if pc != bf.Popcount() + bff.Popcount() {
		panic(fmt.Sprintf("Sum of High and Low Popcounts should match original popcount"))
	}

	bf = bf.CopyBits(bff, bf.Len(), bff.Len())

	if pc != bf.Popcount() {
		panic(fmt.Sprintf("Popcount should have been %v after reconstruction but was %v!", pc, bf.Popcount()))
	}

	return 1
}
