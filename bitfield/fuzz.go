package bitfield

import "fmt"

func Fuzz(data []byte) int {
	bits := uint(len(data) * 8)

	bf := NewBitField(bits)

	for i, b := range data {
		for j := uint(0); j < 8; j++ {
			if (b << j) & 1 > 0 {
				bf.Set(uint(i) * 8 + j)
			}
		}
	}


	for i, b := range data {
		for j := uint(0); j < 8; j++ {
			if (b << j) & 1 > 0 {
				if !bf.Test(uint(i) * 8 + j) {
					panic(fmt.Sprintf("bit %v should have been set!", uint(i)*8+j))
				}
			} else {
				if bf.Test(uint(i) * 8 + j) {
					panic(fmt.Sprintf("bit %v should not have been set!", uint(i)*8+j))
				}
			}
		}
	}

	return 1
}
