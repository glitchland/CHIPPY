package cpu

import (
	cnst "pkg/shared" //constants	
)

type CPU struct {
	// registers array
	/*
		V0	uint8
		V1	uint8
		V2	uint8
		V3	uint8
		V4	uint8
		V5	uint8
		V6	uint8
		V7	uint8
		V8	uint8
		V9	uint8
		VA 	uint8
		VB 	uint8
		VC 	uint8
		VD 	uint8
		VE 	uint8
		VF 	uint8
	*/
	VRegs [16]uint8
	PC  uint16      			// program counter
	I   uint16      			// address register
	Instruction uint16
	HaltedForKeypress      bool //CPU is halted
	AwaitingKeyPressRegister uint8

}

func (c *CPU) IncPC() uint16 {
	c.PC += 2
	return c.PC
}

func (c *CPU) DecPC() uint16 {
	c.PC -= 2 
	return c.PC
}

func (c *CPU) SetVF(b bool) {
	if b {
		c.VRegs[0xf] = 1
	} else {
		c.VRegs[0xf] = 0
	}
}

func (c *CPU) Halt(reg uint8) {
	c.HaltedForKeypress = true
	c.AwaitingKeyPressRegister = reg
}

func (c *CPU) Resume(key uint8) {
	c.HaltedForKeypress = false
	c.VRegs[c.AwaitingKeyPressRegister] = key 
	c.AwaitingKeyPressRegister = cnst.KEYCODE_UNKNOWN
}

func (c *CPU) IsHalted() bool {
	return c.HaltedForKeypress
}