package rrr2

import (
	"testing"
	"github.com/turgon/wavelet/bitfield"
)

func TestNewRRRField(t *testing.T) {
        bf := bitfield.NewBitField(24)

        for i := uint(0); i < 4; i++ {
                bf.Set(i)
        }
        for i := uint(8); i < 24; i+=2 {
                bf.Set(i)
        }

	r := NewRRRField(&bf, 8, 10)
	if r.Class(0) != 4 {
		t.Errorf("NewRRRField has wrong class: %v", r.Class(0))
	}
	if r.Class(1) != 4 {
		t.Errorf("NewRRRField has wrong class: %v", r.Class(1))
	}

	if r.Offset(0) != 0 {
		t.Errorf("NewRRRField has wrong offset: %v", r.Offset(0))
	}
	if r.Offset(1) != 1 {
		t.Errorf("NewRRRField has wrong offset: %v", r.Offset(1))
	}
}

func TestNeedBits(t *testing.T) {
	var j uint8
	for i := uint64(2); i < 128; i *= 2 {
		j++
		if j != needBits(i) {
			t.Errorf("needBits returned %v for %v; should have been %v\n", needBits(i), i, j)
		}
	}
}

func TestBitsForLargest(t *testing.T) {
	var b uint8

	b = bitsForLargest(8)
	if b != 7 {
		t.Errorf("bitsForLargest for 8 was %v but should be 7", b)
	}

	b = bitsForLargest(64)
	if b != 61 {
		t.Errorf("bitsForLargest for 64 was %v but should be 61", b)
	}
}

func TestRRRFieldClass(t *testing.T) {
        bf := bitfield.NewBitField(30)

	bf.Set(5)

	bf.Set(10)
	bf.Set(11)

	bf.Set(15)
	bf.Set(16)
	bf.Set(17)

	bf.Set(20)
	bf.Set(21)
	bf.Set(22)
	bf.Set(23)

	bf.Set(25)
	bf.Set(26)
	bf.Set(27)
	bf.Set(28)
	bf.Set(29)

	r := NewRRRField(&bf, 5, 10)

	for i := uint(0); i <= 5; i++ {
		if r.Class(i) != uint64(i) {
			t.Errorf("Class returned wrong value: %v", r.Class(0))
		}
	}
}

func TestRRRFieldBlock(t *testing.T) {
	bf := bitfield.NewBitField(37)

	pl := []uint{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31}

	for _, p := range pl {
		bf.Set(p)
	}

	r := NewRRRField(&bf, 5, 10)

	bf2 := bitfield.NewBitField(37)

	for i := uint(0); i < r.Len() / uint(r.stepBits); i++ {
		b := r.Block(i)
		bf2 = bf2.CopyBits(b, i * r.blockSize, b.Len())
	}

	if bf.Len() != bf2.Len() {
		t.Fatalf("Rebuilt bitfield lengths don't match")
	}

	for i := uint(0); i < bf.Len(); i++ {
		if bf.Test(i) != bf2.Test(i) {
			t.Errorf("Rebuilt bitfield doesn't match at position %v", i)
		}
	}

	// Again, this time with a repeating pattern of 011 011 bits
	bf = bitfield.NewBitField(16)

	bf.Set(1)
	bf.Set(2)
	bf.Set(4)
	bf.Set(5)
	bf.Set(6)
	bf.Set(7)

	r = NewRRRField(&bf, 3, 1)

	bf2 = bitfield.NewBitField(16)

	for i := uint(0); i < r.Len() / uint(r.stepBits); i++ {
		b := r.Block(i)
		bf2 = bf2.CopyBits(b, i * r.blockSize, b.Len())
	}

	if bf.Len() != bf2.Len() {
		t.Fatalf("Rebuilt bitfield lengths don't match")
	}

	for i := uint(0); i < bf.Len(); i++ {
		if bf.Test(i) != bf2.Test(i) {
			t.Errorf("Rebuilt bitfield doesn't match at position %v", i)
		}
	}
}
