package opcodes

import (
	"fmt"
	"pkg/bits"
)

// Detect take opcode, and print value
func Detect(oc uint16) {
	ocn, err := bits.Nibble(oc, 0)
	if err != nil {
		fmt.Printf("Error getting opcode nibble from %d\n", oc)
	} else {
		fmt.Printf("[%04X] ~> ", oc)
	}
	switch ocn {
	case 0:
		fmt.Println("0NNN [call] Calls RCA 1802 program at address NNN. Not necessary for most ROMs.")
	case 0x00EE:
		fmt.Println("00EE [return] Returns from a subroutine.")
	case 0x00E0:
		fmt.Println("00E0 [disp_clear] Clears the screen.")
	case 0x1:
		fmt.Println("1NNN [goto NNN] Jumps to address NNN.")
	case 0x2:
		fmt.Println("2NNN [*(0xNNN)()] Calls subroutine at NNN.")
	case 0x3:
		fmt.Println("3XNN [if Vx == NN] Skips the next instruction if VX equals NN. (Usually the next instruction is a jump to skip a code block)")
	case 0x4:
		fmt.Println("4XNN [if Vx != NN] Skips the next instruction if VX doesn't equal NN. (Usually the next instruction is a jump to skip a code block)")
	case 0x5:
		fmt.Println("5XY0 [if Vx == Vy] Skips the next instruction if VX equals VY. (Usually the next instruction is a jump to skip a code block)")
	case 0x6:
		fmt.Println("6XNN [Vx = NN] Sets VX to NN.")
	case 0x7:
		fmt.Println("7XNN [Vx += NN] Adds NN to VX. (Carry flag is not changed)")
	case 0x8:
		// handle subtypes
		st, _ := bits.Nibble(oc, 3)
		switch st {
		case 0x0:
			fmt.Println("8XY0 [Vx=Vy] Sets VX to the value of VY.")
		case 0x1:
			fmt.Println("8XY1 [Vx=Vx|Vy] Sets VX to VX or VY. (Bitwise OR operation)")
		case 0x2:
			fmt.Println("8XY2 [Vx=Vx&Vy] Sets VX to VX and VY. (Bitwise AND operation)")
		case 0x3:
			fmt.Println("8XY3 [Vx=Vx^Vy] Sets VX to VX xor VY.")
		case 0x4:
			fmt.Println("8XY4 [Vx += Vy] Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.")
		case 0x5:
			fmt.Println("8XY5 [Vx -= Vy] VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't.")
		case 0x6:
			fmt.Println("8XY6 [Vx=Vy=Vy>>1] Shifts VY right by one and copies the result to VX. VF is set to the value of the least significant bit of VY before the shift.")
		case 0x7:
			fmt.Println("8XY7 [Vx=Vy-Vx] Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't.")
		case 0xE:
			fmt.Println("8XYE [Vx=Vy=Vy<<1] Shifts VY left by one and copies the result to VX. VF is set to the value of the most significant bit of VY before the shift.")
		}
	case 0x9:
		fmt.Println("9XY0 [if(Vx!=Vy)] Skips the next instruction if VX doesn't equal VY. (Usually the next instruction is a jump to skip a code block)")
	case 0xA:
		fmt.Println("ANNN [I = NNN] Sets I to the address NNN.")
	case 0xB:
		fmt.Println("BNNN [PC=V0+NNN] Jumps to the address NNN plus V0.")
	case 0xC:
		fmt.Println("CXNN [Vx=rand()&NN] Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN.")
	case 0xD:
		fmt.Println("DXYN [draw(Vx,Vy,N)] Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels.")
	case 0xE:
		// handle subtypes
		st, _ := bits.Nibble(oc, 3)
		switch st {
		case 0x1:
			fmt.Println("EXA1 [if(key()!=Vx)] Skips the next instruction if the key stored in VX isn't pressed. (Usually the next instruction is a jump to skip a code block)")
		case 0xE:
			fmt.Println("EX9E [if(key()==Vx)] Skips the next instruction if the key stored in VX is pressed. (Usually the next instruction is a jump to skip a code block)")
		}
	case 0xF:
		// handle subtypes
		st := bits.LastByte(oc)
		switch st {
		case 0x07:
			fmt.Println("FX07 [Vx = get_delay()] Sets VX to the value of the delay timer.")
		case 0x0A:
			fmt.Println("FX0A [Vx = get_key()] A key press is awaited, and then stored in VX. (Blocking Operation. All instruction halted until next key event)")
		case 0x15:
			fmt.Println("FX15 [delay_timer(Vx)] Sets the delay timer to VX.")
		case 0x18:
			fmt.Println("FX18 [sound_timer(Vx)] Sets the sound timer to VX.")
		case 0x1E:
			fmt.Println("FX1E [I +=Vx] Adds VX to I.")
		case 0x29:
			fmt.Println("FX29 [I=sprite_addr[Vx]] Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font.")
		case 0x33:
			fmt.Println("FX33 [set_BCD(Vx);] *(I+0)=BCD(3);*(I+1)=BCD(2);*(I+2)=BCD(1);")
		case 0x55:
			fmt.Println("FX55 [reg_dump(Vx,&I)] Stores V0 to VX (including VX) in memory starting at address I. I is increased by 1 for each value written.")
		case 0x65:
			fmt.Println("FX65 [reg_load(Vx,&I)] Fills V0 to VX (including VX) with values from memory starting at address I. I is increased by 1 for each value written.")
		}
	}

}
