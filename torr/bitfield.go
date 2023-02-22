package torr

type Bitfield []byte

func (bf Bitfield) hasPiece(index int) bool {
	bfIndex := index / 8
	offset := index % 8

	return bf[bfIndex]>>(7-offset)&1 != 0
}

func (bf Bitfield) setPiece(index int) {
	byteIndex := index / 8
	offset := index % 8

	bf[byteIndex] |= 1 << (7 - offset)
}
