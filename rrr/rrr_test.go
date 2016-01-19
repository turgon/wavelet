package rrr

import (
	"github.com/turgon/wavelet/bitfield"
	"testing"
	"math/rand"
)

func TestNewRRR(t *testing.T) {
	bf := bitfield.NewBitField(137)

	bf.Set(2)
	bf.Set(10)
	bf.Set(35)

	r := NewRRR(&bf)

	for i := uint(0); i < r.superblocks; i++ {
		for j := uint(0); j < superblockSize; j++ {
			// t.Errorf("%v %v %v\n", i * superblockSize, j, i * superblockSize + j)
		}
	}

	t.Errorf("%v\n", r)
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
	for i := uint(0); i < 8*5; i++ {
		if r.Rank(i) != uint64(i) {
			t.Fatalf("%v\n", r.Rank(i))
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
		for j := rand.Intn(bits/2); j >= 0; j-- {
			bf.Set(uint(rand.Intn(bits)))
		}

		rs[i] = NewRRR(&bf)
	}

	for i := 0; i < b.N; i++ {
		var r RRR
		r = rs[i % rsSize]

		r.Rank(r.bf.Len()-1)
	}
}

func TestPopcountBye(t *testing.T) {
	var x byte

	for i := uint16(0); i < 256; i++ {
		x = byte(i)
		if popcountByte(x) != naivePopCount(x) {
			t.Fatalf("popcountByte did not match naivePopCount for x=%v\n", x)
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
