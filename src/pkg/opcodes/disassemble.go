package opcodes

import (
	"fmt"
	"pkg/bits"
	c "pkg/shared" //constants	
)

// This is for the dissassembler
func DisassembleAndPrint(addr int, oc uint16) {

	fmt.Printf("0x%04X\t%02X %02X\t", addr, ((oc & 0xff00) >> 8), oc&0xff)

	decoded := Disassemble(oc)

	switch decoded.Opcode {
	case c.OP_FLOW_CALL:
		fmt.Printf("CALL 0x%03X\n", decoded.NNN) 

	case c.OP_FLOW_JMP_NNN:
		fmt.Printf("JMP 0x%03X\n", decoded.NNN) 

	case c.OP_FLOW_RETURN:
		fmt.Printf("RET\n") 

	case c.OP_DISPLAY_CLEAR:
		fmt.Printf("CLS\n")

	case c.OP_FLOW_GOTO:
		fmt.Printf("GOTO 0x%03X\n", decoded.NNN)

	case c.OP_COND_SKIP_VX_EQ_NN:
		fmt.Printf("SKIPNXT IF V%1X == 0x%02X\n", decoded.Vx, decoded.NN) 
		// [if Vx == NN] Skips the next instruction if VX equals NN. 
		// (Usually the next instruction is a jump to skip a code block)")

	case c.OP_COND_SKIP_VX_NEQ_NN:
		fmt.Printf("SKIPNXT IF V%1X != 0x%02X\n", decoded.Vx, decoded.NN) 
		//4XNN [if Vx != NN] Skips the next instruction if VX doesn't equal
		// NN. (Usually the next instruction is a jump to skip a code block)")

	case c.OP_COND_SKIP_VX_EQ_VY:
		fmt.Printf("SKIPNXT IF V%1X == V%01X\n", decoded.Vx, decoded.Vy) 
		//5XY0 [if Vx == Vy] Skips the next instruction if VX equals VY. 
		//(Usually the next instruction is a jump to skip a code block)")

	case c.OP_CONST_VX_EQUALS_NN:
		fmt.Printf("V%1X = %02X\n", decoded.Vx, decoded.NN) 
		//"6XNN [Vx = NN] Sets VX to NN.")

	case c.OP_CONST_VX_PLUS_EQUALS_NN:
		fmt.Printf("V%1X += %02X\n", decoded.Vx, decoded.NN) 
		//"7XNN [Vx += NN] Adds NN to VX. (Carry flag is not changed)")

	case c.OP_ASSIGN_VX_VY:
		fmt.Printf("V%1X = V%1X\n", decoded.Vx, decoded.Vy) 
		//"8XY0 [Vx=Vy] Sets VX to the value of VY.")

	case c.OP_BITOP_VX_EQUALS_VX_OR_VY:
		fmt.Printf("V%1X = V%1X | V%1X\n", decoded.Vx, decoded.Vx, decoded.Vy) 
		//"8XY1 [Vx=Vx|Vy] Sets VX to VX or VY. (Bitwise OR operation)")

	case c.OP_BITOP_VX_EQUALS_VX_AND_VY:
		fmt.Printf("V%1X = V%1X & V%1X\n", decoded.Vx, decoded.Vx, decoded.Vy) 
		//"8XY2 [Vx=Vx&Vy] Sets VX to VX and VY. (Bitwise AND operation)")

	case c.OP_BITOP_VX_EQUALS_VX_XOR_VY:
		fmt.Printf("V%1X = V%1X ^ V%1X\n", decoded.Vx, decoded.Vx, decoded.Vy) 
		//"8XY3 [Vx=Vx^Vy] Sets VX to VX xor VY.")

	case c.OP_MATH_VX_EQUALS_VX_PLUS_VY:
		fmt.Printf("V%1X += V%1X\t\t; VF = 1 when there's a carry, else 0\n", decoded.Vx, decoded.Vy) 
		//"8XY4 [Vx += Vy] Adds VY to VX. VF is set to 1 when there's a carry, 
		// and to 0 when there isn't.")

	case c.OP_MATH_VX_EQUALS_VX_MINUS_VY:
		fmt.Printf("V%1X -= V%1X\t\t; VF = 0 when there's a borrow, else 1\n", decoded.Vx, decoded.Vy) 
		// "8XY5 [Vx -= Vy] VY is subtracted from VX. VF is set to 0 when there's a borrow, 
		// and 1 when there isn't.")

	case c.OP_BITOP_VX_EQUALS_VX_RSHIFT:
		fmt.Printf("V%1X = V%1X = V%1X >> 1\t\t; VF = least significant bit of V%1X before the shift\n", decoded.Vx, decoded.Vy, decoded.Vy, decoded.Vy) 
		// "8XY6 [Vx=Vy=Vy>>1] Shifts VY right by one and copies the result to VX. VF is set to the 
		// value of the least significant bit of VY before the shift.")

	case c.OP_MATH_VX_EQUALS_VY_MINUS_VX:
		fmt.Printf("V%1X = V%1X - V%1X\t\t; VF = 0 when there's a borrow, else 1\n", decoded.Vx, decoded.Vy, decoded.Vx) 
		// "8XY7 [Vx=Vy-Vx] Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 
		// when there isn't.")

	case c.OP_BITOP_VX_EQUALS_VX_LSHIFT:
		fmt.Printf("V%1X = V%1X = V%1X << 1\t\t; VF = most significant bit of V%1X before the shift\n", decoded.Vx, decoded.Vy, decoded.Vy, decoded.Vy) 
		// "8XYE [Vx=Vy=Vy<<1] Shifts VY left by one and copies the result to VX. VF is 
		// set to the value of the most significant bit of VY before the shift.")

	case c.OP_COND_SKIP_VX_NEQ_VY:
		fmt.Printf("SKIPNXT IF V%1X != V%1X\n", decoded.Vx, decoded.Vy) 
		// "9XY0 [if(Vx!=Vy)] Skips the next instruction if VX doesn't equal VY. 
		// (Usually the next instruction is a jump to skip a code block)")

	case c.OP_MEM_I_EQUALS_NNN:
		fmt.Printf("I = 0x%03X\n", decoded.NNN) 
		// "ANNN [I = NNN] Sets I to the address NNN.")

	case c.OP_MEM_JMP_PC_EQUALS_V0_PLUS_NNN:
		fmt.Printf("PC = V0 + 0x%03X\n", decoded.NNN) 
		// "BNNN [PC=V0+NNN] Jumps to the address NNN plus V0.")

	case c.OP_RAND:
		fmt.Printf("V%1X = RND() & 0x%02X\n", decoded.Vx, decoded.NN) 
		// "CXNN [Vx=rand()&NN] Sets VX to the result of a bitwise and operation on a 
		// random number (Typically: 0 to 255) and NN.")

	case c.OP_DISPLAY_DRAWSPRITE:
		fmt.Printf("DRAWSPRITE(V%1X,V%1X,%1X)\n", decoded.Vx, decoded.Vy, decoded.N) 
		// "DXYN [draw(Vx,Vy,N)] Draws a sprite at coordinate (VX, VY) that has a
		//  width of 8 pixels and a height of N pixels.")

	case c.OP_KEY_SKIPNXT_NEQ_VX:
		fmt.Printf("SKIPNXT IF KEYPRESS() != #V%1X\n", decoded.Vx) 
		// "EXA1 [if(key()!=Vx)] Skips the next instruction if the key stored in VX 
		// isn't pressed. (Usually the next instruction is a jump to skip a code block)")

	case c.OP_KEY_SKIPNXT_EQ_VX:
		fmt.Printf("SKIPNXT IF KEYPRESS() == #V%1X\n", decoded.Vx) 
		// "EX9E [if(key()==Vx)] Skips the next instruction if the key stored in VX is 
		// pressed. (Usually the next instruction is a jump to skip a code block)")

	case c.OP_VX_EQUALS_DELAY_TIMER:
		fmt.Printf("V%1X = GET_DELAY_TIMER()\n", decoded.Vx) 
		//"FX07 [Vx = get_delay()] Sets VX to the value of the delay timer.")
	
	case c.OP_KEY_VX_EQUALS_KEYPRESS:
		fmt.Printf("V%1X = GET_KEYPRESS()\n", decoded.Vx) 
		// "FX0A [Vx = get_key()] A key press is awaited, and then stored in VX. 
		// (Blocking Operation. All instruction halted until next key event)")

	case c.OP_DELAY_TIMER_EQUALS_VX:
		fmt.Printf("TIMER_DELAY = V%1X\n", decoded.Vx) 
		// "FX15 [delay_timer(Vx)] Sets the delay timer to VX.")

	case c.OP_SOUND_TIMER_EQUALS_VX:
		fmt.Printf("TIMER_SOUND = V%1X\n", decoded.Vx) 
		// "FX18 [sound_timer(Vx)] Sets the sound timer to VX.")

	case c.OP_MEM_I_EQUALS_I_PLUS_VX:
		fmt.Printf("I += V%1X\n", decoded.Vx) 
		// "FX1E [I +=Vx] Adds VX to I.")

	case c.OP_MEM_I_EQUALS_SPRITE_MAP_CHAR_VX:
		fmt.Printf("I = SPRITE_MAP(#V%1X)\n", decoded.Vx) 
		// "FX29 [I=sprite_addr[Vx]] Sets I to the location of the sprite for the 
		// character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font.")

	case c.OP_BCD_VX:
		fmt.Printf("BCD(V%1X)\t\t\t; *(I+0)=BCD(3) , *(I+1)=BCD(2),  *(I+2)=BCD(1)\n", decoded.Vx) 
		// "FX33 [set_BCD(Vx);] *(I+0)=BCD(3);*(I+1)=BCD(2);*(I+2)=BCD(1);")

	case c.OP_MEM_STORE_REGS_AT_I_PTR:
		fmt.Printf("STORE((V0 ... V%1X) -> #I)\n", decoded.Vx) 
		// "FX55 [reg_dump(Vx,&I)] Stores V0 to VX (including VX) in memory starting at 
		// address I. I is increased by 1 for each value written.")

	case c.OP_MEM_LOAD_REGS_FROM_I_PTR:
		fmt.Printf("LOAD(#I -> (V0 ... V%1X))\n", decoded.Vx) 
		// "FX65 [reg_load(Vx,&I)] Fills V0 to VX (including VX) with values from memory 
		// starting at address I. I is increased by 1 for each value written.")
	}
}

// Disassemble opcode
func Disassemble(oc uint16) c.DecodedOpcode {

	decoded := c.DecodedOpcode{}

	ocn, err := bits.Nibble(oc, 0)
	if err != nil {
		fmt.Printf("Error getting opcode nibble from %d\n", oc)
	}

	decoded.Vx = uint8(((oc & 0x0F00) >> 8))
	decoded.Vy = uint8(((oc & 0x00F0) >> 4))
	decoded.N = uint8((oc & 0x000F))
	decoded.NN = uint8((oc & 0x00FF))
	decoded.NNN = (oc & 0x0FFF)

	switch ocn {
	case 0x0:
		// 0 is a little awkward, so we handle it differently
		if (oc == 0x00E0) {     
			decoded.Opcode = c.OP_DISPLAY_CLEAR
		} else if (oc == 0x00EE) { 
			decoded.Opcode = c.OP_FLOW_RETURN
		} else {
			decoded.Opcode = c.OP_FLOW_JMP_NNN
		}
		return decoded 
	case 0x1:
		decoded.Opcode = c.OP_FLOW_GOTO
		return decoded 
	case 0x2:
		decoded.Opcode = c.OP_FLOW_CALL
		return decoded 
	case 0x3:
		decoded.Opcode = c.OP_COND_SKIP_VX_EQ_NN 
		return decoded 
	case 0x4:
		decoded.Opcode = c.OP_COND_SKIP_VX_NEQ_NN
		return decoded 
	case 0x5:
		decoded.Opcode = c.OP_COND_SKIP_VX_EQ_VY
		return decoded 
	case 0x6:
		decoded.Opcode = c.OP_CONST_VX_EQUALS_NN
		return decoded 
	case 0x7:
		decoded.Opcode = c.OP_CONST_VX_PLUS_EQUALS_NN
		return decoded 
	case 0x8:
		// handle subtypes
		st, _ := bits.Nibble(oc, 3)
		switch st {
		case 0x0:
			decoded.Opcode = c.OP_ASSIGN_VX_VY
			return decoded 
		case 0x1:
			decoded.Opcode = c.OP_BITOP_VX_EQUALS_VX_OR_VY
			return decoded 
		case 0x2:
			decoded.Opcode = c.OP_BITOP_VX_EQUALS_VX_AND_VY
			return decoded 
		case 0x3:
			decoded.Opcode = c.OP_BITOP_VX_EQUALS_VX_XOR_VY
			return decoded 
		case 0x4:
			decoded.Opcode = c.OP_MATH_VX_EQUALS_VX_PLUS_VY
			return decoded 
		case 0x5:
			decoded.Opcode = c.OP_MATH_VX_EQUALS_VX_MINUS_VY
			return decoded 
		case 0x6:
			decoded.Opcode = c.OP_BITOP_VX_EQUALS_VX_RSHIFT
			return decoded 
		case 0x7:
			decoded.Opcode = c.OP_MATH_VX_EQUALS_VY_MINUS_VX
			return decoded 
		case 0xE:
			decoded.Opcode = c.OP_BITOP_VX_EQUALS_VX_LSHIFT
			return decoded 
		}
	case 0x9:
		decoded.Opcode = c.OP_COND_SKIP_VX_NEQ_VY
		return decoded 
	case 0xA:
		decoded.Opcode = c.OP_MEM_I_EQUALS_NNN
		return decoded 
	case 0xB:
		decoded.Opcode = c.OP_MEM_JMP_PC_EQUALS_V0_PLUS_NNN
		return decoded 
	case 0xC:
		decoded.Opcode = c.OP_RAND
		return decoded 
	case 0xD:
		decoded.Opcode = c.OP_DISPLAY_DRAWSPRITE
		return decoded 
	case 0xE:
		// handle subtypes
		st, _ := bits.Nibble(oc, 3)
		switch st {
		case 0x1:
			decoded.Opcode = c.OP_KEY_SKIPNXT_NEQ_VX
			return decoded 
		case 0xE:
			decoded.Opcode = c.OP_KEY_SKIPNXT_EQ_VX
			return decoded 
		}
	case 0xF:
		// handle subtypes
		st := bits.LastByte(oc)
		switch st {
		case 0x07:
			decoded.Opcode = c.OP_VX_EQUALS_DELAY_TIMER
			return decoded 
		case 0x0A:
			decoded.Opcode = c.OP_KEY_VX_EQUALS_KEYPRESS
			return decoded 
		case 0x15:
			decoded.Opcode = c.OP_DELAY_TIMER_EQUALS_VX
			return decoded 
		case 0x18:
			decoded.Opcode = c.OP_SOUND_TIMER_EQUALS_VX
			return decoded
		case 0x1E:
			decoded.Opcode = c.OP_MEM_I_EQUALS_I_PLUS_VX
			return decoded 
		case 0x29:
			decoded.Opcode = c.OP_MEM_I_EQUALS_SPRITE_MAP_CHAR_VX
			return decoded 
		case 0x33:
			decoded.Opcode = c.OP_BCD_VX
			return decoded 
		case 0x55:
			decoded.Opcode = c.OP_MEM_STORE_REGS_AT_I_PTR
			return decoded 
		case 0x65:
			decoded.Opcode = c.OP_MEM_LOAD_REGS_FROM_I_PTR
			return decoded 
		}
	}

	decoded.Opcode = c.OP_INVALID
	return decoded
}
