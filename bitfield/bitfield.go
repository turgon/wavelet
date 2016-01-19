// Package bitfield provides a low level data structure for bag-of-bits.
package bitfield

type segment uint16
const segmentSize uint = 16

// BitField uses a slice of bytes to store the bag-of-bits.
type BitField struct {
	Data []segment
	length uint
}

// NewBitField returns a BitField sized to the desired number of bits.
func NewBitField(bits uint) BitField {
	bytes := bits / segmentSize
	if bits % segmentSize > 0 {
		bytes++
	}

	return BitField{
		make([]segment, bytes),
		bits,
	}
}

// Set sets the bit at position to one.
func (bf BitField) Set(position uint) {
	bf.Data[position / segmentSize] |= (1 << ((segmentSize - 1) - position % segmentSize))
}

// Unset sets the bit at position to zero.
func (bf BitField) Unset(position uint) {
	bf.Data[position / segmentSize] &^= (1 << ((segmentSize - 1) - position % segmentSize))
}

// Test returns true if the bit at position is set.
func (bf BitField) Test(position uint) bool {
	return (bf.Data[position / segmentSize] & (1 << ((segmentSize - 1) - position % segmentSize))) != 0
}

// Len returns the number of bits in the bitfield.
func (bf BitField) Len() uint {
	return bf.length
}

// Resize returns a new BitField with new size and a copy of the
// original data. If the new copy is larger than the original, it
// will be padded with 0-bits. If smaller, bits are truncated.
func (bf BitField) Resize(bits uint) BitField {
	nbf := NewBitField(bits)

	copy(nbf.Data, bf.Data)

	if nbf.Len() % segmentSize == 0 {
		return nbf
	}

	ctr := uint(0)
	for i := nbf.Len(); i < bf.Len() && ctr < (segmentSize - (nbf.Len() % segmentSize)); i++ {
		nbf.Unset(i)
		ctr++
	}

	return nbf
}
