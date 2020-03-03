package memory

import (
	"fmt"
	s "pkg/shared"
	b "pkg/bits"
)

const RAMEND = 0xFFF
const RAMSTART = 0x000
const CODE = 0x200
const ETI_660 = 0x600
const RESERVED = 0x1FF

/*
	Memory Map:
	+---------------+= 0xFFF (4095) End of Chip-8 RAM
	|               |
	|               |
	|               |
	|               | <- fonts are somewhere after the code
	|               |
	| 0x200 to 0xFFF|
	|     Chip-8    |
	| Program / Data|
	|     Space     |
	|               |
	|               |
	|               |
	+- - - - - - - -+= 0x600 (1536) Start of ETI 660 Chip-8 programs
	|               |
	|               |
	|               |
	+---------------+= 0x200 (512) Start of most Chip-8 programs
	| 0x000 to 0x1FF|
	| Reserved for  |
	|  interpreter  |
	+---------------+= 0x000 (0) Start of Chip-8 RAM
*/

type RAM struct {
	ram Mem
	FontAddress uint16
	Top uint16
	CodeEnd uint16
}

func (r *RAM) Init() {
	r.Top = uint16(RAMEND - 1)
	r.ram.Init(RAMEND)
}

func (r *RAM) ToStr() string {
	s := ""
	for i := 1; i <= RAMEND; i++ {
		v, e := r.ram.GetWordAt(uint16(i - 1))
		if e != nil {
			panic("Unable to read memory")
		}
		s += fmt.Sprintf("%02X ", v)
		if i%16 == 0 {
			s += fmt.Sprintf("\n")
		}
	}
	return s
}

func (r *RAM) GetCodeSegment() uint16 {
	return CODE
}

func (r *RAM) GetFontSegment() uint16 {
	return r.FontAddress
}

func (r *RAM) LoadCode(romData []uint16) {
	address := CODE // set the code base address
	for _, instruction := range romData {
		r.Write16(uint16(address), instruction)
		address += 2
	}
	r.CodeEnd = uint16(address-2) // this marks the end of the code
	// update where we should pack fonts
	r.FontAddress = uint16(address) 

	for _, fontByte := range s.FONT_SPRITES {
		r.Write8(uint16(address), fontByte)
		address++
	}
}

func (r *RAM) Write2BytesAsWord(addr uint16, hiByte uint8, loByte uint8) {
	wrd := b.BytesToWord(hiByte, loByte)
	_ = r.ram.SetWordAt(uint16(addr), wrd)
}

func (r *RAM) Write16(addr uint16, v uint16) {
	_ = r.ram.SetWordAt(uint16(addr), v)
}

func (r *RAM) Read16(addr uint16) uint16 {
	v, _ := r.ram.GetWordAt(uint16(addr))
	return v
}

func (r *RAM) Write8(addr uint16, v uint8) {
	_ = r.ram.SetByteAt(uint16(addr), v)
}

func (r *RAM) Read8(addr uint16) uint8 {
	v, _ := r.ram.GetByteAt(uint16(addr))
	return v
}

func (r *RAM) Read16N(addr uint16, length int) []uint16 {

	l := uint16(length)
	chunk := make([]uint16, l)

	if (!r.boundsValid(addr+l)) {
		return chunk;
	}

	for i := uint16(0); i < l; i++ {
		wrd := r.Read16(addr+i)
		chunk = append(chunk, wrd)
	}

	return chunk
}

func (r *RAM) Read8N(addr uint16, length int) []uint8 {

	l := uint16(length)
	chunk := make([]uint8, l)

	if (!r.boundsValid(addr+l)) {
		return chunk;
	}

	for i := uint16(0); i < l; i++ {
		wrd := r.Read8(addr+i)
		chunk[i] = wrd
	}

	return chunk
}

func (r *RAM) ValidCodeAddress(addr uint16) bool {
	if (addr > r.CodeEnd) {
		return false;
	} else {
		return true;
	}
}

// make this an error
func (r *RAM) boundsValid(addr uint16) bool {
	if (addr >= RAMEND) {
		return false 
	} else {
		return true
	}
}