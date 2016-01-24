package bitfield

import (
	"testing"
)

func TestNewBitField(t *testing.T) {
	var bf BitField

	bf = NewBitField(16)
	if len(bf.Data) != 1 {
		t.Errorf("NewBitField returned wrongly sized field: %v", len(bf.Data))
	}

	bf = NewBitField(17)
	if len(bf.Data) != 2 {
		t.Errorf("NewBitField returned wrongly sized field: %v", len(bf.Data))
	}
}

func TestBitFieldTest(t *testing.T) {
	bf := NewBitField(17)
	bf.Set(0)
	if !bf.Test(0) {
		t.Errorf("BitField Test returned wrong value!")
	}
	if bf.Test(1) {
		t.Errorf("BitField Test returned wrong value!")
	}
	if bf.Test(8) {
		t.Errorf("BitField Test returned wrong value!")
	}
	if bf.Test(9) {
		t.Errorf("BitField Test returned wrong value!")
	}
	if bf.Test(16) {
		t.Errorf("BitField Test returned wrong value!")
	}
	if bf.Test(19) {
		t.Errorf("BitField Test returned wrong value!")
	}
}

func TestBitFieldUnset(t *testing.T) {
	bf := NewBitField(8)
	bf.Set(3)
	bf.Unset(3)
	if bf.Data[0] != 0 {
		t.Errorf("BitField Unset failed to clear bit!")
	}
}

func TestBitFieldSet(t *testing.T) {

	bf := NewBitField(16)
	for i := uint(0); i < 16; i++ {
		bf.Set(i)
	}
	for i := uint(0); i < 16; i++ {
		if !bf.Test(i) {
			t.Errorf("BitField Set or Test failed!")
		}
	}

	bf = NewBitField(17)

	bf.Set(0)
	if bf.Data[0] != 32768 {
		t.Errorf("BitField Set wrong value!")
	}

	bf.Set(1)
	if bf.Data[0] != 49152 {
		t.Errorf("BitField Set wrong value!")
	}

	bf.Set(8)
	if bf.Data[0] != 49280 {
		t.Errorf("BitField Set wrong value!")
	}

	bf.Set(9)
	if bf.Data[0] != 49344 {
		t.Errorf("BitField Set wrong value!")
	}

	bf.Set(16)
	if bf.Data[1] != 32768 {
		t.Errorf("BitField Set wrong value!")
	}

}

func TestLen(t *testing.T) {
	bf := NewBitField(17)

	if bf.Len() != 17 {
		t.Errorf("BitField Len returned wrong length!")
	}
}

func TestResize(t *testing.T) {
	bf := NewBitField(2)
	bf.Set(0)
	bf.Set(1)

	bf = bf.Resize(2)
	if !bf.Test(0) || !bf.Test(1) {
		t.Errorf("Resize didn't preserve bits!")
	}

	bf = bf.Resize(segmentSize)
	if !bf.Test(0) || !bf.Test(1) {
		t.Errorf("Resize didn't preserve bits!")
	}

	bf = bf.Resize(1)
	if !bf.Test(0) {
		t.Errorf("Resize didn't preserve bits!")
	}

	bf = bf.Resize(2)
	if !bf.Test(0) {
		t.Errorf("Resize didn't preserve bits!")
	}
	if bf.Test(1) {
		t.Errorf("Resize didn't pad bits!")
	}

	bf = NewBitField(17)
	for i := uint(0); i < 17; i++ {
		bf.Set(i)
	}
	bf = bf.Resize(1)
	if bf.Data[0] != 32768 {
		t.Errorf("Resize didn't clear bits!")
	}

	bf = NewBitField(9)
	bf.Set(1)
	bf = bf.Resize(8)
	if len(bf.Data) != 1 {
		t.Errorf("Resize to 8 has wrong data size!")
	}

	bf = NewBitField(1)
	bf.Set(0)
	bf = bf.Resize(0)
	if len(bf.Data) != 0 {
		t.Errorf("Resize to zero still has data!")
		t.Errorf("%v", bf)
	}
}

func TestSub(t *testing.T) {
	bf := NewBitField(3)
	bf.Set(0)
	bf.Set(1)
	nbf := bf.Sub(1, 3)
	if !nbf.Test(0) || nbf.Test(1) {
		t.Errorf("Sub did not properly set bits!")
	}
	if nbf.Len() != 2 {
		t.Errorf("Sub has wrong length!")
	}

	nbf = bf.Sub(1, 99)
	if nbf.Len() != 2 {
		t.Errorf("Sub has wrong length!")
	}
}

func TestPopcount16(t *testing.T) {
	for i := uint64(0); i < 65536; i++ {
		j := popcount16(uint16(i))
		var cnt uint64
		for k := uint8(0); k < 16; k++ {
			if (i >> k) & 1 > 0 {
				cnt++
			}
		}
		if j != cnt {
			t.Errorf("popcount16 is wrong: %v = %v", j, cnt)
		}
	}
}

func TestBitFieldPopcount(t *testing.T) {
	bf := NewBitField(33)
	for i := uint(0); i < 33; i++ {
		if bf.Popcount() != uint64(i) {
			t.Errorf("Popcount is wrong: %v = %v", i, bf.Popcount())
		}
		bf.Set(i)
	}
}

func TestNewBitFieldFromUint64(t *testing.T) {
	for i := uint64(0); i < 131072; i++ {
		bf := NewBitFieldFromUint64(17, i)
		if i != bf.Uint64(17) {
			t.Errorf("NewBitFieldFromUint64 of %v could not be restored", i)
		}
	}

	bf := NewBitField(100)
	bf.Set(0)
	bf.Set(65)
	if bf.Uint64(65) != 1 {
		t.Errorf("Uint64 failed to limit n to 64 bits")
	}
	if bf.Uint64(1000) != 1 {
		t.Errorf("Uint64 failed to limit n to 64 bits")
	}
}

func TestCopyBits(t *testing.T) {
	bf := NewBitField(5)
	bf.Set(0)
	bf.Set(2)

	length := bf.Len()
	for i := 0; i < 2; i++ {
		length *= 2
		bf = bf.CopyBits(bf, bf.Len(), bf.Len())
	}
	if length != bf.Len() {
		t.Errorf("CopyBits left us with the wrong bitlength")
	}
	for i := uint(0); i < bf.Len(); i++ {
		if i % 5 == 0 || i % 5 == 2 {
			if !bf.Test(i) {
				t.Errorf("CopyBits failed to properly copy bits from source")
			}
		} else {
			if bf.Test(i) {
				t.Errorf("CopyBits failed to properly copy bits from source")
			}
		}
	}
}
