package rrr2

import (
	"math"
	"math/big"

	"github.com/turgon/wavelet/bitfield"
)

type RRRField struct {
	bitfield.BitField

	blockSize uint
	superSize uint

	superRanks uint64

	classBits uint
	offsetBits uint

	stepBits uint

	lastBlockLen uint

	global map[uint64][]bitfield.BitField
}

// Class fetches the class of the block at the given position.
func (r RRRField) Class(block uint) uint64 {
	start := r.stepBits * block
	end := start + r.classBits

	bf := r.Sub(start, end)
	return bf.Uint64(r.classBits)
}

// Offset fetches the offset of the block at the given position.
func (r RRRField) Offset(block uint) uint64 {
	start := r.stepBits * block + r.classBits
	end := start + r.offsetBits

	bf := r.Sub(start, end)

	return bf.Uint64(r.offsetBits)
}

func (r RRRField) Block(block uint) bitfield.BitField {
	b := r.global[r.Class(block)][r.Offset(block)]

	if block == r.Len() / r.stepBits - 1 {
		b = b.Resize(r.lastBlockLen)
	}

	return b
}

// NewRRRField returns an RRRField generated from the bitfield
// given as bf. The bitfield will be blocked in blockSize bits,
// and each superSize blocks will be summarized.
func NewRRRField(bf *bitfield.BitField, blockSize uint, superSize uint) RRRField {
	var r RRRField

	r.blockSize = blockSize
	r.superSize = superSize

	r.classBits = needBits(uint64(blockSize + 1))
	r.offsetBits = bitsForLargest(blockSize)

	r.stepBits = r.classBits + r.offsetBits

	r.global = make(map[uint64][]bitfield.BitField, 0)

	offMap := make(map[uint64]map[uint64]uint64)
	offMax := make(map[uint64]uint64)

	for i := uint(0); i < bf.Len(); i += blockSize {
		sub := bf.Sub(i, i + blockSize)
		subpc := sub.Popcount()

		r.lastBlockLen = sub.Len()

		// create a bitfield consisting of the number of
		// bits necessary to describe any class, and set
		// it to this class.
		cf := bitfield.NewBitFieldFromUint64(r.classBits, subpc)
		r.BitField = r.CopyBits(cf, r.stepBits * (i / blockSize), cf.Len())

		if _, ok := offMap[subpc]; !ok {
			offMap[subpc] = make(map[uint64]uint64)
			r.global[subpc] = make([]bitfield.BitField, 0)
		}

		subval := sub.Uint64(sub.Len())

		var offset uint64

		if os, ok := offMap[subpc][subval]; !ok {
			offMap[subpc][subval] = offMax[subpc]
			offset = offMax[subpc]
			r.global[subpc] = append(r.global[subpc], sub)
			offMax[subpc]++
		} else {
			offset = os
		}


		of := bitfield.NewBitFieldFromUint64(r.offsetBits, offset)
		r.BitField = r.CopyBits(of, r.stepBits * (i / blockSize) + r.classBits, of.Len())
	}

	return r
}

// needBits takes a number x and returns the number of bits needed
// to store the range [0, x).
func needBits(x uint64) uint {
	return uint(math.Ceil(math.Log2(float64(x))))
}

// bitsForLargest takes a number n of RRR classes, figures out the
// size of the largest class, and returns the number of bits needed
// to enumerate that class. For example, given 8 classes, the largest
// class has 8/2 = 4 bits set with (8 choose 4) = 70 permutations.
// Thus, bitsForLargest(8) = ceiling(log2(70)) = 7 bits.
func bitsForLargest(n uint) uint {

	if n > 64 {
		panic("bitForLargest: n can't exceed 64")
	}

	i := new(big.Int)
	i.Binomial(int64(n), int64(n/2))
	return needBits(uint64(i.Int64()))
}
