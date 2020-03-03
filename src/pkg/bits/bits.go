package bits

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// LastByte returns the last byte of a word
func LastByte(input uint16) uint16 {
	return input & 0xFF
}

// FirstByte returns the first byte of a word
func FirstByte(input uint16) uint16 {
	return ((input & 0xFF00) >> 8)
}

// Nibble fetches a nibble at offset
func Nibble(input uint16, offset uint8) (uint16, error) {
	if offset > 3 || offset < 0 {
		return 0, &errorString{"Cannot extract at index"}
	}
	mask := uint16(1 << (15 - offset*4))
	mask |= uint16(1 << (14 - offset*4))
	mask |= uint16(1 << (13 - offset*4))
	mask |= uint16(1 << (12 - offset*4))
	input &= mask
	return (input >> (12 - (offset * 4))), nil
}

// FlipBit flips a bit at offset
func FlipBit(input uint16, offset uint8) uint16 {
	input ^= (1 << (15 - offset))
	return input
}

// SetBit sets a bit at offset
func SetBit(input uint16, offset uint8) uint16 {
	input |= (1 << (15 - offset))
	return input
}

// GetBit gets a bit at offset
func GetBit16(input uint16, offset uint8) bool {
	if IsBitSet16(input, offset) {
		return true
	}
	return false
}

func GetBit8(input uint8, offset uint8) bool {
	if IsBitSet8(input, offset) {
		return true
	}
	return false
}

// BoolToIntStr converts a bool to a string
func BoolToIntStr(v bool) string {
	if v {
		return "1"
	}
	return "0"
}

// UnsetBit unsets a bit at offset
func UnsetBit(input uint16, offset uint8) uint16 {
	input &^= (1 << (15 - offset))
	return input
}

// IsBitSet checks if a bit at position is set
func IsBitSet16(n uint16, pos uint8) bool {
	val := n & (1 << (15 - pos))
	return (val > 0)
}

// indexed left to right
func IsBitSet8(n uint8, pos uint8) bool {
	val := n & (1 << (7 - pos))
	return (val > 0)
}

func BytesToWord(hiByte uint8, loByte uint8) uint16 {
	wrd := (uint16(loByte) & 0x00FF) 
	wrd += (uint16(hiByte) << 8)
	return wrd
}

func WordToBytes(wrd uint16) (uint8, uint8) {
	loByte := uint8(wrd & 0x00FF) 
	hiByte := uint8((wrd >> 8) & 0x00FF)
	return hiByte, loByte
}