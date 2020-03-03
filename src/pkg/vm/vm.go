package vm

import (
	"fmt"
	scrn "pkg/screen"
	"pkg/memory"
	decode "pkg/opcodes"
	"pkg/cpu"
	c "pkg/shared" //constants	
	"math/rand"
	"time"
)

type VirtualMachine struct {
	//timer
	Cpu cpu.CPU
	Ram memory.RAM
	Stack memory.Stack
	Screen *scrn.Screen
	TimerRate int // in hertz (60)
	DelayTimer int
	SoundTimer int
	KeyState [16]bool
}

/*
 TODO: - constants for VReg index
	   - dont draw if screen is not attached
	   - add halt after code complete, and make it exit main
*/

func (vm *VirtualMachine) task() {
    for range time.Tick(166 * time.Millisecond){ 
			if (vm.DelayTimer > 0) {
				vm.DelayTimer--
			}

			if (vm.SoundTimer > 0) {
				vm.SoundTimer--
			}
    }
}

func (vm *VirtualMachine) Init(romData []uint16) {
	vm.Ram.Init()
	vm.Ram.LoadCode(romData)
	vm.Cpu.PC = vm.Ram.GetCodeSegment()
	vm.Cpu.I = vm.Ram.GetFontSegment()
	go vm.task()
	rand.Seed(time.Now().UnixNano())
}

func (vm *VirtualMachine) KeyPress(key uint8) {
	if(vm.Cpu.IsHalted()) {
		vm.Cpu.Resume(key)
	}
	
	vm.KeyState[key] = true
}

func (vm *VirtualMachine) KeyRelease(key uint8) {	
	vm.KeyState[key] = false
}


// TODO , halt if screen is not attached?
func (vm *VirtualMachine) AttachScreen(s *scrn.Screen) {
	vm.Screen = s
}

func (vm *VirtualMachine) FetchDecodeExecute() {

	if(vm.Cpu.IsHalted()) {
		return;
	}

	// TODO: investigate how to make this cleaner
	// if we alter the flow with ret, goto, jmp etc
	// we want to skip pc +=2 
	//
	flowChanged := false

	// read instruction
	if(vm.Ram.ValidCodeAddress(vm.Cpu.PC)) {
		instruction := vm.Ram.Read16(vm.Cpu.PC)
	
		// decode instruction
		decode.DisassembleAndPrint(int(vm.Cpu.PC), instruction)
		decoded := decode.Disassemble(instruction)

		// execute handler
		switch decoded.Opcode {
		// OP_FLOW
		case c.OP_FLOW_CALL, 
			c.OP_FLOW_RETURN, 
			c.OP_FLOW_GOTO:
			flowChanged = vm.handleExecFlow(decoded)
		case c.OP_DISPLAY_CLEAR, 
			c.OP_DISPLAY_DRAWSPRITE:
			flowChanged = vm.handleDisplay(decoded)
		case c.OP_COND_SKIP_VX_EQ_NN, 
			c.OP_COND_SKIP_VX_NEQ_NN,
			c.OP_COND_SKIP_VX_EQ_VY, 
			c.OP_COND_SKIP_VX_NEQ_VY:
			flowChanged = vm.handleConditional(decoded)
		case c.OP_CONST_VX_EQUALS_NN, 
		     c.OP_CONST_VX_PLUS_EQUALS_NN:
			flowChanged = vm.handleConst(decoded)
		case c.OP_ASSIGN_VX_VY:
			flowChanged = vm.handleAssign(decoded)		
		case c.OP_BITOP_VX_EQUALS_VX_OR_VY, 
			c.OP_BITOP_VX_EQUALS_VX_AND_VY,
			c.OP_BITOP_VX_EQUALS_VX_XOR_VY, 
			c.OP_BITOP_VX_EQUALS_VX_RSHIFT,
			c.OP_BITOP_VX_EQUALS_VX_LSHIFT:
			flowChanged = vm.handleBitOp(decoded)
		case c.OP_MATH_VX_EQUALS_VX_PLUS_VY,
			c.OP_MATH_VX_EQUALS_VX_MINUS_VY,
			c.OP_MATH_VX_EQUALS_VY_MINUS_VX:
			flowChanged =vm.handleMath(decoded)
		case c.OP_MEM_I_EQUALS_NNN,
			c.OP_MEM_JMP_PC_EQUALS_V0_PLUS_NNN,
			c.OP_MEM_I_EQUALS_I_PLUS_VX,
			c.OP_MEM_I_EQUALS_SPRITE_MAP_CHAR_VX,
			c.OP_MEM_STORE_REGS_AT_I_PTR,
			c.OP_MEM_LOAD_REGS_FROM_I_PTR:
			flowChanged = vm.handleMemory(decoded) 
		case c.OP_KEY_SKIPNXT_NEQ_VX,
			c.OP_KEY_SKIPNXT_EQ_VX,
			c.OP_KEY_VX_EQUALS_KEYPRESS:
			flowChanged = vm.handleKeyOp(decoded)
		case c.OP_VX_EQUALS_DELAY_TIMER,
			c.OP_DELAY_TIMER_EQUALS_VX,
			c.OP_SOUND_TIMER_EQUALS_VX:
			flowChanged = vm.handleTimer(decoded)
		case c.OP_RAND:
			flowChanged = vm.handleRand(decoded)
		case c.OP_BCD_VX:
			flowChanged = vm.handleBCD(decoded)
		}

		// update program counter
		// might need new context
		if (!flowChanged) {
			vm.Cpu.PC += 2 
		}

		//fmt.Printf("%#v\n", vm.Cpu)
		vm.printState() 
		time.Sleep(time.Duration(1)*time.Millisecond)
	}

}

func (vm *VirtualMachine) printState() {
	state := "-------------------------\n"
	for i := 0; i < len(vm.Cpu.VRegs); i++ {
		state += fmt.Sprintf(" V%X = %X ", i, vm.Cpu.VRegs[i])
	}
	state +=  fmt.Sprintf("\n PC = %X I = %X\n", vm.Cpu.PC, vm.Cpu.I)
	fmt.Printf("%s\n-------------------------\n", state)
}

func (vm *VirtualMachine) handleExecFlow(instruction c.DecodedOpcode) bool {
	switch instruction.Opcode {
		case c.OP_FLOW_CALL:	
			vm.Stack.Push(vm.Cpu.PC + 2)
			vm.Cpu.PC = instruction.NNN 
		case c.OP_FLOW_RETURN:
			vm.Cpu.PC = vm.Stack.Pop()
		case c.OP_FLOW_GOTO:
			vm.Cpu.PC = instruction.NNN 
		case c.OP_FLOW_JMP_NNN:			
			vm.Cpu.PC = instruction.NNN
	}
	// 00EE return 
	// 1NNN goto  
	// 2NNN call 
	// BNNN jump [PC=V0+NNN]
	// 0NNN call XXX: ?
	return true 
}

func (vm *VirtualMachine) handleDisplay(instruction c.DecodedOpcode) bool {

	switch instruction.Opcode {
	case c.OP_DISPLAY_CLEAR:
		// 00E0 clear 
		vm.Screen.Clear()
	case c.OP_DISPLAY_DRAWSPRITE:
		// DXYN draw(Vx,Vy,N) 
		// Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels. 
		// Each row of 8 pixels is read as bit-coded starting from memory location I; I value doesn’t 
		// change after the execution of this instruction. As described above, VF is set to 1 if any 
		// screen pixels are flipped from set to unset when the sprite is drawn, and to 0 if that doesn’t
		// happen.

		// Get sprite buffer and position
		vx := vm.Cpu.VRegs[instruction.Vx] 
		vy := vm.Cpu.VRegs[instruction.Vy] // up+down
		
		fmt.Printf("DXYN draw(Vx[%d],Vy[%d],N(len)[%d])\n", vx, vy, instruction.N)

		//spriteWords := vm.Ram.Read16N(vm.Cpu.I, int(instruction.N/2))
		//spriteBytes := make([]uint8, 0)
		sprite := vm.Ram.Read8N(vm.Cpu.I, int(instruction.N))
		fmt.Printf("DXYN: %#v\n", sprite)
		vm.Cpu.SetVF(vm.Screen.Draw(int(vx), int(vy), sprite) )
	}
	return false 
}

func (vm *VirtualMachine) handleConditional(instruction c.DecodedOpcode) bool {
	// 3XNN if(Vx==NN) -> Skips the next instruction if VX equals NN. 
	// 4XNN if(Vx!=NN) -> Skips the next instruction if VX doesn't equal NN. 
	// 5XY0 if(Vx==Vy) -> Skips the next instruction if VX equals VY. 
	// 9XY0 if(Vx!=Vy) -> Skips the next instruction if VX doesn't equal VY. 

	vx := vm.Cpu.VRegs[instruction.Vx]
	vy := vm.Cpu.VRegs[instruction.Vy]
	nn := instruction.NN

	switch instruction.Opcode {
	case c.OP_COND_SKIP_VX_EQ_NN:
		fmt.Printf("OP_COND_SKIP_VX_EQ_NN: |%x| |%x| |%x|\n", instruction.Vx, vx, nn)
		if (vx == nn) {
			vm.Cpu.IncPC()
			vm.Cpu.IncPC()			
			return true		
		}
	case c.OP_COND_SKIP_VX_NEQ_NN:
		if (vx != nn) {
			vm.Cpu.IncPC()	
			vm.Cpu.IncPC()			
			return true	
		}
	case c.OP_COND_SKIP_VX_EQ_VY:
		if (vx == vy) {
			vm.Cpu.IncPC()
			vm.Cpu.IncPC()			
			return true	
		}
	case c.OP_COND_SKIP_VX_NEQ_VY:
		if (vx != vy) {
			vm.Cpu.IncPC()	
			vm.Cpu.IncPC()			
			return true
		}		
	}
	return false 
}

func (vm *VirtualMachine) handleConst(instruction c.DecodedOpcode) bool {
	switch instruction.Opcode {
	case c.OP_CONST_VX_EQUALS_NN:	
		// 6XNN Vx = NN  -> Sets VX to NN. 
		vm.Cpu.VRegs[instruction.Vx] = instruction.NN
	case c.OP_CONST_VX_PLUS_EQUALS_NN:
		// 7XNN Vx += NN -> Adds NN to VX. (Carry flag is not changed) 
		vm.Cpu.VRegs[instruction.Vx] = (vm.Cpu.VRegs[instruction.Vx] + instruction.NN)
	}
	return false 
}

func (vm *VirtualMachine) handleAssign(instruction c.DecodedOpcode) bool {
	fmt.Printf("handleAssign: [%#v]\n", instruction)
	vm.Cpu.VRegs[instruction.Vx] = vm.Cpu.VRegs[instruction.Vy]
	// 8XY0 Vx=Vy  -> Sets VX to the value of VY. 
	return false 
}

func (vm *VirtualMachine) handleBitOp(instruction c.DecodedOpcode) bool {
	// 8XY1 Vx=Vx|Vy -> Sets VX to VX or VY. (Bitwise OR operation) 
	// 8XY2 Vx=Vx&Vy -> Sets VX to VX and VY. (Bitwise AND operation) 
	// 8XY3 Vx=Vx^Vy -> Sets VX to VX xor VY. 
	// 8XY6 Vx>>=1   -> Stores the least significant bit of VX in VF and then shifts VX to the right by 1.
	// 8XYE Vx<<=1   -> Stores the most significant bit of VX in VF and then shifts VX to the left by 1.
	vx := vm.Cpu.VRegs[instruction.Vx] 
	vy := vm.Cpu.VRegs[instruction.Vy]

	switch instruction.Opcode {
	case c.OP_BITOP_VX_EQUALS_VX_OR_VY:
		vm.Cpu.VRegs[instruction.Vx] = vx | vy  
	case c.OP_BITOP_VX_EQUALS_VX_AND_VY:
		vm.Cpu.VRegs[instruction.Vx] = vx & vy  		
	case c.OP_BITOP_VX_EQUALS_VX_XOR_VY: 
		vm.Cpu.VRegs[instruction.Vx] = vx ^ vy  
	case c.OP_BITOP_VX_EQUALS_VX_RSHIFT:
		vm.Cpu.SetVF((vm.Cpu.VRegs[instruction.Vx] & 0x01) == 1)
		vm.Cpu.VRegs[instruction.Vx] = vx >> 1		
	case c.OP_BITOP_VX_EQUALS_VX_LSHIFT:
		vm.Cpu.SetVF((vm.Cpu.VRegs[instruction.Vx] >> 7) == 1)
		vm.Cpu.VRegs[instruction.Vx] = vx << 1		
	}

	return false
}

func (vm *VirtualMachine) handleMath(instruction c.DecodedOpcode) bool {
	// 8XY4 Vx += Vy -> Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't. 
	// 8XY5 Vx -= Vy -> VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't. 
	// 8XY7 Vx=Vy-Vx -> Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't. 
	vx := vm.Cpu.VRegs[instruction.Vx] 
	vy := vm.Cpu.VRegs[instruction.Vy]

	switch instruction.Opcode {
	case c.OP_MATH_VX_EQUALS_VX_PLUS_VY:
		vm.Cpu.SetVF((uint16(vm.Cpu.VRegs[instruction.Vx]) + uint16(vm.Cpu.VRegs[instruction.Vy])) > 0xff)
		vm.Cpu.VRegs[instruction.Vx] = vx + vy
	case c.OP_MATH_VX_EQUALS_VX_MINUS_VY:
		vm.Cpu.SetVF(!(vm.Cpu.VRegs[instruction.Vx] < vm.Cpu.VRegs[instruction.Vy]))
		vm.Cpu.VRegs[instruction.Vx] = vx - vy		
	case c.OP_MATH_VX_EQUALS_VY_MINUS_VX:	
		vm.Cpu.SetVF((vm.Cpu.VRegs[instruction.Vx] < vm.Cpu.VRegs[instruction.Vy]))
		vm.Cpu.VRegs[instruction.Vx] = vy - vx		
	}	
	return false 
}

func (vm *VirtualMachine) handleMemory(instruction c.DecodedOpcode) bool {
	// ANNN I = NNN -> Sets I to the address NNN. 
	// FX1E I +=Vx  -> Adds VX to I. VF is set to 1 when there is a range overflow (I+VX>0xFFF), and to 0 when there isn't
	// FX29 I=sprite_addr[Vx] -> Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font. 
	// FX55 reg_dump(Vx,&I) -> Stores V0 to VX (including VX) in memory starting at address I. The offset from I is increased by 1 for each value written, but I itself is left unmodified
	// FX65 reg_load(Vx,&I) -> Fills V0 to VX (including VX) with values from memory starting at address I. The offset from I is increased by 1 for each value written, but I itself is left unmodified.
	switch instruction.Opcode {
	case c.OP_MEM_I_EQUALS_NNN:
		vm.Cpu.I = instruction.NNN
	case c.OP_MEM_JMP_PC_EQUALS_V0_PLUS_NNN:
		vm.Cpu.PC =  uint16(vm.Cpu.VRegs[0]) + instruction.NNN //TODO: Hide this direct assignment in the CPU
		return true 
	case c.OP_MEM_I_EQUALS_I_PLUS_VX:	
		vm.Cpu.I = vm.Cpu.I + uint16(vm.Cpu.VRegs[instruction.Vx])
	case c.OP_MEM_I_EQUALS_SPRITE_MAP_CHAR_VX:
		vx := vm.Cpu.VRegs[instruction.Vx]
		fontBase := vm.Ram.GetFontSegment()
		vm.Cpu.I = fontBase + uint16(c.FONT_BYTE_LEN * vx)
	case c.OP_MEM_STORE_REGS_AT_I_PTR:	
		//Stores V0 to VX (including VX) in memory starting at address I
		for i := 0; i <= int(instruction.Vx); i++ {
			//fmt.Printf("HANDLE MEMORY STORE: i:%d %d\n", i, vm.Cpu.VRegs[i])
			vm.Ram.Write8(vm.Cpu.I + uint16(i), vm.Cpu.VRegs[i])
		}
	case c.OP_MEM_LOAD_REGS_FROM_I_PTR:
		//Fills V0 to VX (including VX) with values from memory starting at address I.
		for i := 0; i <= int(instruction.Vx); i++ {
			v := vm.Ram.Read8(vm.Cpu.I + uint16(i))
			//fmt.Printf("HANDLE MEMORY LOAD : i:%d %d\n", i, v)
			vm.Cpu.VRegs[i] = v
		}
	}
	return false 
}

func (vm *VirtualMachine) handleRand(instruction c.DecodedOpcode) bool {
	// CXNN Vx=rand()&NN -> Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN. 
	num := uint8(rand.Intn(255)) & instruction.NN 
	vm.Cpu.VRegs[instruction.Vx] = num
	return false
}

func (vm *VirtualMachine) handleBCD(instruction c.DecodedOpcode) bool {
	// Store BCD representation of Vx in memory locations I, I+1, and I+2.

	// The interpreter takes the decimal value of Vx, and places the hundreds 
	// digit in memory at location in I, the tens digit at location I+1, and 
	// the ones digit at location I+2.

	b0 := vm.Cpu.VRegs[instruction.Vx] / 100
	b1 := (vm.Cpu.VRegs[instruction.Vx] % 100) / 10
	b2 := vm.Cpu.VRegs[instruction.Vx] % 10

	vm.Ram.Write8(vm.Cpu.I+0, b0) 
	vm.Ram.Write8(vm.Cpu.I+1, b1) 
	vm.Ram.Write8(vm.Cpu.I+2, b2) 

	return false 
}

func (vm *VirtualMachine) handleKeyOp(instruction c.DecodedOpcode) bool {
	// EX9E if(key()==Vx)  -> Skips the next instruction if the key stored in VX is pressed. 
	// EXA1 if(key()!=Vx)  -> Skips the next instruction if the key stored in VX isn't pressed. 
	// FX0A Vx = get_key() -> A key press is awaited, and then stored in VX. (Blocking Operation. All instruction halted until next key event) 

	vx := instruction.Vx

	switch instruction.Opcode {
	case c.OP_KEY_SKIPNXT_NEQ_VX:
		if (!vm.KeyState[vx]) {
			vm.Cpu.IncPC()	
			vm.Cpu.IncPC()			
			return true
		}	
	case c.OP_KEY_SKIPNXT_EQ_VX:
		if (vm.KeyState[vx]) {
			vm.Cpu.IncPC()	
			vm.Cpu.IncPC()			
			return true
		}	
	case c.OP_KEY_VX_EQUALS_KEYPRESS:
		vm.Cpu.Halt(vx)
	}
	return false
}

func (vm *VirtualMachine) handleTimer(instruction c.DecodedOpcode) bool {

	switch instruction.Opcode {
	case c.OP_VX_EQUALS_DELAY_TIMER: // timers
		vm.Cpu.VRegs[instruction.Vx] = uint8(vm.DelayTimer) 
	case c.OP_DELAY_TIMER_EQUALS_VX:
		vm.DelayTimer = int(vm.Cpu.VRegs[instruction.Vx])
	case c.OP_SOUND_TIMER_EQUALS_VX:
		vm.SoundTimer = int(vm.Cpu.VRegs[instruction.Vx])
	}
	// FX07 Vx = get_delay() -> Sets VX to the value of the delay timer. 
	// FX15 delay_timer(Vx)  -> Sets the delay timer to VX. 
	// FX18 sound_timer(Vx) -> Sets the sound timer to VX. 
	return false 
}

func (vm *VirtualMachine) handleSound(instruction c.DecodedOpcode) bool {
	return false 
}