package rrr

import (
	"github.com/turgon/wavelet/bitfield"
)

// for now let's go with block size of 8 bits
// and superblock size of 4 blocks.
// so a superblock is 32 bits.

const blockSize = 8      // bits
const superblockSize = 16 // blocks

type RRRPair struct {
	// only 3 bits needed here
	Class uint8

	// 8 choose 4 = 70 is the max, requiring 7 bits
	Offset uint8
}

type RRR struct {
	bf          *bitfield.BitField
	blocks      uint
	superblocks uint

	superRanks []uint64

	global map[uint8][]byte
	pairs []RRRPair
}

func NewRRR(bf *bitfield.BitField) RRR {
	// how many blocks are in this bitfield?
	blocks := uint(bf.Len() / blockSize)
	if bf.Len()%blockSize > 0 {
		blocks++
	}
	superblocks := uint(blocks / superblockSize)
	if blocks%superblockSize > 0 {
		superblocks++
	}

	globals := make(map[uint8][]byte)

	r := RRR{
		bf,
		blocks,
		superblocks,
		make([]uint64, superblocks, superblocks),
		globals,
		make([]RRRPair, blocks),
	}

	if bf.Len() == 0 {
		return r
	}

	popMap := make(map[byte]uint64)
	popCnt := make(map[uint64]uint64)

	// Build popMap, then invert it into globals
	for _, b := range bf.Data {
		popMap[b] = popcountByte(b)
		if _, ok := 
			popCnt[popMap[b]]++
	}
	for pc, cnt := range popCnt {
		globals[uint8(pc)] = make([]byte, cnt)
	}
	for b, pc := range popMap {
		globals[uint8(pc)] = append(globals[uint8(pc)], b)
	}

	var tot uint64
	for i := uint(0); i < superblocks-1; i++ {
		for j := uint(0); j < superblockSize; j++ {
			loc := i * superblockSize + j
			if loc < uint(len(bf.Data)) {
				tot += popMap[bf.Data[loc]]
			}
		}
		r.superRanks[i+1] = tot
	}

	return r
}

func (r *RRR) Rank(x uint) (tot uint64) {
	// I need to return the sum of ranks for all prior superblocks
	// plus the sum of the ranks for all blocks within this
	// superblock prior to x,
	// plus the rank of this block up to x.

	length := r.bf.Len()

	if x >= length {
		x = length - 1
	}

	if length == 0 {
		return 0
	}

	block := x / blockSize

	superblock := block / superblockSize

	tot = r.superRanks[superblock]

	for i := superblock * superblockSize; i < block; i++ {
		tot += popcountByte(r.bf.Data[i])
	}

	blockOffset := x % blockSize

	if block != uint(len(r.bf.Data)) {
		tot += popcountByte((r.bf.Data[block] >> (8 - blockOffset)))
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
