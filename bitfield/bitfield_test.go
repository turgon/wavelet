package bitfield

import (
	"testing"
)

func TestNewBitField(t *testing.T) {
	var bf BitField

	bf = NewBitField(16)
	if len(bf.Data) != 2 {
		t.Errorf("NewBitField returned wrongly sized field: %v", len(bf.Data))
	}

	bf = NewBitField(17)
	if len(bf.Data) != 3 {
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
	if bf.Data[0] != 128 {
		t.Errorf("BitField Set wrong value!")
	}

	bf.Set(1)
	if bf.Data[0] != 192 {
		t.Errorf("BitField Set wrong value!")
	}

	bf.Set(8)
	if bf.Data[1] != 128 {
		t.Errorf("BitField Set wrong value!")
	}

	bf.Set(9)
	if bf.Data[1] != 192 {
		t.Errorf("BitField Set wrong value!")
	}

	bf.Set(16)
	if bf.Data[2] != 128 {
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
	if bf.Data[0] != 128 {
		t.Errorf("Resize didn't clear bits!")
	}

	bf = NewBitField(9)
	bf.Set(1)
	bf = bf.Resize(8)
	if len(bf.Data) != 1 {
		t.Errorf("Resize to 8 has wrong data size!")
	}

	bf = NewBitField(1)
	bf.Set(1)
	bf = bf.Resize(0)
	if len(bf.Data) != 0 {
		t.Errorf("Resize to zero still has data!")
	}
}
