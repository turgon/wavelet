package rrr

import (
	"github.com/turgon/wavelet/bitfield"
	"math/rand"
	"testing"
)

func TestRRRPair(t *testing.T) {
	for i := uint8(0); i < 16; i++ {
		for j := uint8(0); j < 16; j++ {
			rrrp := NewRRRPair(i, j)
			if rrrp.Class() != i {
				t.Fatalf("RRRPair has wrong class")
			}
			if rrrp.Offset() != j {
				t.Fatalf("RRRPair has wrong offset")
			}
		}
	}
}

func TestNewRRR(t *testing.T) {
	bf := bitfield.NewBitField(153)

	bf.Set(2)
	bf.Set(10)
	bf.Set(35)
	bf.Set(129)
	bf.Set(152)

	r := NewRRR(&bf)

	for i := uint(0); i < r.superblocks; i++ {
		for j := uint(0); j < superblockSize; j++ {
			// t.Errorf("%v %v %v\n", i * superblockSize, j, i * superblockSize + j)
		}
	}

}

func TestNewRRREmptyField(t *testing.T) {
	bf := bitfield.NewBitField(0)
	r := NewRRR(&bf)
	r.Rank(0)
}

func TestRank(t *testing.T) {
	bf := bitfield.NewBitField(8 * 5)
	for i := uint(0); i < 8*5; i++ {
		bf.Set(i)
	}
	r := NewRRR(&bf)
	for i := uint(1); i <= 8*5; i++ {
		if r.Rank(i) != uint64(i) {
			t.Errorf("%v\n", r)
			t.Errorf("%v\n", r.bf.Len())
			t.Errorf("%v\n", r.superRanks)
			t.Fatalf("%v %v\n", i, r.Rank(i))
		}
	}
	// Now test an edge case
	r.Rank(8*5 + 1)
}

func BenchmarkRank(b *testing.B) {

	rsSize := 1000

	rs := make([]RRR, rsSize)

	for i := 0; i < rsSize; i++ {
		var bf bitfield.BitField

		bits := 2 + rand.Intn(8190)

		bf = bitfield.NewBitField(uint(bits))
		for j := rand.Intn(bits / 2); j >= 0; j-- {
			bf.Set(uint(rand.Intn(bits)))
		}

		rs[i] = NewRRR(&bf)
	}

	for i := 0; i < b.N; i++ {
		var r RRR
		r = rs[i%rsSize]

		r.Rank(r.bf.Len() - 1)
	}
}

func TestPopcountByte(t *testing.T) {
	var x byte

	for i := uint16(0); i < 256; i++ {
		x = byte(i)
		if popcountByte(x) != naivePopCount(x) {
			t.Fatalf("popcountByte did not match naivePopCount for x=%v\n", x)
		}
	}
}

func TestPopcount16(t *testing.T) {
	for i := uint16(0); i < 65535; i++ {
		if popcount16(i) != (naivePopCount(byte(i >> 8)) + naivePopCount(byte((i << 8) >> 8))) {
			t.Fatalf("popcount16 did not match naivePopCount for i=%v\n", i)
		}
	}
}

func naivePopCount(x byte) (c uint64) {
	for i := 0; i < 8; i++ {
		if x&1 > 0 {
			c++
		}
		x >>= 1
	}
	return c
}
