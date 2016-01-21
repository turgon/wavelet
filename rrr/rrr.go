package rrr

import (
	"fmt"
	"github.com/turgon/wavelet/bitfield"
)

// for now let's go with block size of 16 bits
// and superblock size of 8 blocks.
// so a superblock is 128 bits.

// since a block is 16 bits, its class can be described using just
// 4 bits, and there are less than 14 different combinations of bits
// at its max of 16 choose 8, which also needs for 4 bits. this means
// a uint8 can describe both the class and offset perfectly.

const blockSize = 16     // bits
const superblockSize = 8 // blocks

type RRRPair struct {
	value uint8
}

func NewRRRPair(class uint8, offset uint8) RRRPair {
	var v uint8

	v = class
	v <<= 4

	v |= offset

	return RRRPair{
		v,
	}
}

func (rrrp RRRPair) Class() uint8 {
	return rrrp.value >> 4
}

func (rrrp RRRPair) Offset() uint8 {
	return rrrp.value & 15
}

type RRR struct {
	bf          *bitfield.BitField
	length	    uint
	blocks      uint
	superblocks uint

	superRanks []uint64

	global map[uint8][]uint16
	pairs  []RRRPair
}

func NewRRR(bf *bitfield.BitField) RRR {
	// how many blocks are in this bitfield?
	bl := uint(bf.Len() / blockSize)
	if bf.Len()%blockSize > 0 {
		bl++
	}
	sbl := uint(bl / superblockSize)
	if bl%superblockSize > 0 {
		sbl++
	}

	r := RRR{
		bf,
		bf.Len(),
		bl,
		sbl,
		make([]uint64, sbl, sbl),
		make(map[uint8][]uint16),
		make([]RRRPair, bl),
	}

	if bf.Len() == 0 {
		return r
	}

	popMap := make(map[uint16]uint64)
	popCnt := make(map[uint64]uint64)
	seen := make(map[uint16]bool)
	offsets := make(map[uint16]int)

	for _, b := range bf.Data {
		bu := uint16(b)
		popMap[bu] = popcount16(bu)
		if !seen[bu] {
			popCnt[popMap[bu]]++
			seen[bu] = true
		}
	}
	for pc, cnt := range popCnt {
		r.global[uint8(pc)] = make([]uint16, cnt)
	}
	for b, pc := range popMap {
		pos := len(r.global[uint8(pc)]) - 1
		r.global[uint8(pc)][pos] = b
		offsets[b] = pos
	}

	for i, b := range bf.Data {
		bu := uint16(b)
		pc := popMap[bu]
		cl := uint8(pc)
		os := uint8(offsets[bu])
		r.pairs[i] = NewRRRPair(cl, os)
		r.superRanks[i/superblockSize] += pc
	}

	return r
}

func (r *RRR) Rank(x uint) (tot uint64) {
	// I need to return the sum of ranks for all prior superblocks
	// plus the sum of the ranks for all blocks within this
	// superblock prior to x,
	// plus the rank of this block up to x.

	if x >= r.length {
		x = r.length - 1
	}

	if r.length == 0 {
		return 0
	}

	block := x / blockSize

	superblock := block / superblockSize

	fmt.Printf("Starting Rank; x = %v, r.length = %v, block = %v, superblock = %v\n", x, r.length, block, superblock)

	// First set tot to the sum of all the preceding superblocks.
	if superblock > 0 {
		tot = r.superRanks[superblock]
		fmt.Printf("Added super rank to tot, now = %v\n", tot)
	}

	// Next add the class value of all the preceding blocks within
	// this superblock.
	for i := superblock * superblockSize; i < block; i++ {
		// tot += popcount16(uint16(r.bf.Data[i]))
		tot += uint64(r.pairs[i].Class())
		fmt.Printf("Added class of block %v to tot, now = %v\n", i, tot)
	}

	// Finally, add the popcount of this block, shifted to exclude
	// bits after position x.

	blockOffset := x % blockSize

	if block != uint(len(r.pairs)) {
		fmt.Printf("Entire RRR: %v\n", r)
		fmt.Printf("All RRRPairs: %v\n", r.pairs)
		fmt.Printf("RRRPair is %v\n", r.pairs[block])
		os := r.pairs[block].Offset()
		cl := r.pairs[block].Class()
		fmt.Printf("class is %v\n", cl)
		bls := r.global[cl]
		fmt.Printf("offset is %v\n", os)
		bl := bls[os]
		fmt.Printf("blockOffset is %v\n", blockOffset)
		bu := bl >> (superblockSize - blockOffset)
		tot += popcount16(bu)
		fmt.Printf("Added remainder of block %v = %16.16b to tot, now = %v\n", block, bu, tot)
		// tot += popcount16((uint16(r.bf.Data[block]) >> (16 - blockOffset)))
	}

	return tot
}

func popcountByte(z byte) uint64 {

	b0 := z & 85
	b1 := (z >> 1) & 85

	c := b0 + b1

	d0 := c & 51
	d2 := (c >> 2) & 51

	e := d0 + d2

	f0 := e & 15
	f4 := (e >> 4) & 15

	return uint64(f0 + f4)
}

func popcount16(z uint16) uint64 {

	b0 := z & 21845
	b1 := (z >> 1) & 21845

	c := b0 + b1

	d0 := c & 13107
	d2 := (c >> 2) & 13107

	e := d0 + d2

	f0 := e & 3855
	f4 := (e >> 4) & 3855

	g := f0 + f4

	h0 := g & 255
	h8 := (g >> 8) & 255

	return uint64(h0 + h8)
}
