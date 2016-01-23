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

func (bf BitField) Sub(left uint, right uint) BitField {
	nbf := NewBitField(right - left)

	var position uint

	for i := left; i < right; i++ {
		if bf.Test(i) {
			nbf.Set(position)
		}

		position++
	}

	return nbf
}

// Popcount returns the number of bits set to 1 in the BitField.
func (bf BitField) Popcount() uint64 {
	var total uint64

	for _, b := range bf.Data {
		total += popcount16(uint16(b))
	}

	return total
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
