package bitfield

type BitField struct {
	Data []byte
	length uint
}

// NewBitField returns a BitField sized to the desired number of bits.
func NewBitField(bits uint) BitField {
	bytes := bits / 8
	if bits % 8 > 0 {
		bytes++
	}

	return BitField{
		make([]byte, bytes),
		bits,
	}
}

func (bf BitField) Set(position uint) {
	bf.Data[position / 8] |= (1 << (7 - position % 8))
}

func (bf BitField) Unset(position uint) {
	bf.Data[position / 8] &^= (1 << (7 - position % 8))
}

func (bf BitField) Test(position uint) bool {
	return (bf.Data[position / 8] & (1 << (7 - position % 8))) != 0
}

func (bf BitField) Len() uint {
	return bf.length
}

// Resize returns a new BitField with new size and a copy of the
// original data. If the new copy is larger than the original, it
// will be padded with 0-bits. If smaller, bits are truncated.
func (bf BitField) Resize(bits uint) BitField {
	nbf := NewBitField(bits)

	copy(nbf.Data, bf.Data)

	ctr := uint(0)
	for i := nbf.Len(); i < bf.Len() && ctr < (8 - (nbf.Len() % 8)); i++ {
		nbf.Unset(i)
		ctr++
	}

	return nbf
}
