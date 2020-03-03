package screen

import (
	"bytes"
	"fmt"
	"pkg/bits"
)

// http://craigthomas.ca/blog/2015/02/19/writing-a-chip-8-emulator-draw-command-part-3/
/*
	x – this specifies the register that stores the X coordinate where you want to draw the sprite.
		 Valid X coordinates range from 0 to 63. Values larger than 63 will cause the sprite to wrap
		 horizontally across the screen.
	y – this specifies the register that stores the Y coordinate where you want to draw the sprite.
		 Valid X coordinate range from 0 to 31. Values larger than 31 will cause the sprite to wrap
		 vertically across the screen.
    n – this specifies how many bytes the sprite is. Valid number of bytes range from 0 to 15.

	Graphics are drawn to the screen solely by drawing sprites, which are 8 pixels wide and may
	be from 1 to 15 pixels in height. Sprite pixels are XOR'd with corresponding screen pixels.
	In other words, sprite pixels that are set flip the color of the corresponding screen pixel,
	while unset sprite pixels do nothing. The carry flag (VF) is set to 1 if any screen pixels are
	flipped from set to unset when a sprite is drawn and set to 0 otherwise. This is used for
	collision detection.

"E"	 Binary  	Hex
**** 11110000   0xF0
*    10000000   0x80
**** 11110000   0xF0
*    10000000   0x80
**** 11110000   0xF0
These sprites are 5 bytes long, or 8x5 pixels.
*/

// WIDTH of LCD screen
const WIDTH = 64

// HEIGHT of LCD screen
const HEIGHT = 32

const PIXEL_ON = 0xff
const PIXEL_OFF = 0x66

type Screen struct {
	Height int
	Width  int
	Pixels [][]uint8
}

func (s *Screen) Init() {

	s.Height = HEIGHT
	s.Width = WIDTH

	s.Pixels = make([][]uint8, s.Height)
	for r := range s.Pixels {
		s.Pixels[r] = make([]uint8, s.Width*4)
	}

	s.Clear()
}

// The original implementation of the Chip-8 language used a
// 64x32-pixel monochrome display with this format:
// +----------------------------+
// |(0,0)                 (63,0)|
// |                            |
// |(0,31)               (63,31)|
// +----------------------------+
func (s *Screen) SetPixel(x int, y int, bitindex int, row int, ox int, oy int) {

	fmt.Printf("Inside draw (setpixel) X: %d Y: %d bitindex: %d, row: %d || original x: %d y: %d\n", x, y, bitindex, row, ox, oy)
	if !s.validPixelIndex(x, y) {
		return
	}

	s.Pixels[y][(x*4)+0] = PIXEL_ON
	s.Pixels[y][(x*4)+1] = PIXEL_ON
	s.Pixels[y][(x*4)+2] = PIXEL_ON
	s.Pixels[y][(x*4)+3] = PIXEL_ON
}

func (s *Screen) UnsetPixel(x int, y int) {

	if !s.validPixelIndex(x, y) {
		return
	}

	s.Pixels[y][(x*4)+0] = PIXEL_OFF
	s.Pixels[y][(x*4)+1] = PIXEL_OFF
	s.Pixels[y][(x*4)+2] = PIXEL_OFF
	s.Pixels[y][(x*4)+3] = PIXEL_OFF
}

func (s *Screen) Clear() {
	for r := range s.Pixels {
		for c := range s.Pixels[r] {
			s.Pixels[r][c] = PIXEL_OFF
		}
	}
}

func (s *Screen) IsPixelSet(x int, y int) bool {

	if !s.validPixelIndex(x, y) {
		return false
	}

	if s.Pixels[y][(x*4)+0] == PIXEL_ON &&
		s.Pixels[y][(x*4)+1] == PIXEL_ON &&
		s.Pixels[y][(x*4)+2] == PIXEL_ON &&
		s.Pixels[y][(x*4)+3] == PIXEL_ON {
		return true
	} else {
		return false
	}
}

// return a value for CPU VF
func (s *Screen) Draw(x int, y int, sprite []uint8) bool {

	setFlag := false
	//for row, b := range sprite {
	for spriteByteIndex := 0; spriteByteIndex < len(sprite); spriteByteIndex++ {
		b := sprite[spriteByteIndex]	
		for bitindex := 0; bitindex < 8; bitindex++ {
			if bits.GetBit8(b, uint8(bitindex)) {
				// if pixel is already set, unset it, set VF flag in CPU to 1
				if s.IsPixelSet(x + bitindex, spriteByteIndex + y) {
					s.UnsetPixel(x + bitindex, spriteByteIndex + y)
					setFlag = true
				} else {
					fmt.Printf("setpixel :%d %#v\n", spriteByteIndex, sprite)
					s.SetPixel(x + bitindex, spriteByteIndex + y, bitindex, spriteByteIndex, x, y)
				}
			}
		}
	}
	return setFlag
}

func (s *Screen) validPixelIndex(x int, y int) bool {
	if (x < 0 || x > WIDTH-1) || (y < 0 || y > HEIGHT-1) {
		fmt.Printf("X:[%d] Y:[%d] Out of bounds!\n", x, y)
		return false
	}
	return true
}

// DRAWSPRITE
func (s *Screen) RefreshPixelBytes() {
	s.Clear()
}

func (s *Screen) GetPixels() []uint8 {
	return bytes.Join(s.Pixels, nil)
}
