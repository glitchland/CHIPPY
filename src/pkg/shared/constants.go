package shared

const GLFW_RELEASE = 0
const GLFW_PRESS = 1

// Keycodes
const (
	KEYCODE_UNKNOWN = uint8(100)
	KEYCODE_0 = 0x0
	KEYCODE_1 = 0x1
	KEYCODE_2 = 0x2
	KEYCODE_3 = 0x3
	KEYCODE_4 = 0x4
	KEYCODE_5 = 0x5
	KEYCODE_6 = 0x6
	KEYCODE_7 = 0x7
	KEYCODE_8 = 0x8
	KEYCODE_9 = 0x9
	KEYCODE_A = 0xA
	KEYCODE_B = 0xB
	KEYCODE_C = 0xC
	KEYCODE_D = 0xE
	KEYCODE_E = 0xD
	KEYCODE_F = 0xF
)
// Sprites
// 0xF?? - 0xFFF built in 4x5 pixel font set, A-F, 1-9.
// packed at the end of the code
var FONT_SPRITES = [...]uint8 {
	0xF0, 0x90, 0x90, 0x90, 0xF0,  //0
	0x20, 0x60, 0x20, 0x20, 0x70,  //1
	0xF0, 0x10, 0xF0, 0x80, 0xF0,  //2
	0xF0, 0x10, 0xF0, 0x10, 0xF0,  //3
	0x90, 0x90, 0xF0, 0x10, 0x10,  //4
	0xF0, 0x80, 0xF0, 0x10, 0xF0,  //5
	0xF0, 0x80, 0xF0, 0x90, 0xF0,  //6
	0xF0, 0x10, 0x20, 0x40, 0x40,  //7
	0xF0, 0x90, 0xF0, 0x90, 0xF0,  //8
	0xF0, 0x90, 0xF0, 0x10, 0xF0,  //9
	0xF0, 0x90, 0xF0, 0x90, 0x90,  //A
	0xE0, 0x90, 0xE0, 0x90, 0xE0,  //B
	0xF0, 0x80, 0x80, 0x80, 0xF0,  //C
	0xE0, 0x90, 0x90, 0x90, 0xE0,  //D
	0xF0, 0x80, 0xF0, 0x80, 0xF0,  //E
	0xF0, 0x80, 0xF0, 0x80, 0x80,  //F 
}

const OP_INVALID     = 0

// Control flow
const OP_FLOW_CALL   = 1
const OP_FLOW_RETURN = 2
const OP_FLOW_GOTO   = 3
const OP_FLOW_JMP_NNN = 4

// Display
const OP_DISPLAY_CLEAR = 5
const OP_DISPLAY_DRAWSPRITE = 6

// Conditionals
const OP_COND_SKIP_VX_EQ_NN   = 7
const OP_COND_SKIP_VX_NEQ_NN  = 8
const OP_COND_SKIP_VX_EQ_VY   = 9
const OP_COND_SKIP_VX_NEQ_VY  = 10

const OP_CONST_VX_EQUALS_NN      = 11
const OP_CONST_VX_PLUS_EQUALS_NN = 12

const OP_ASSIGN_VX_VY = 13         

const OP_BITOP_VX_EQUALS_VX_OR_VY  = 14
const OP_BITOP_VX_EQUALS_VX_AND_VY = 15
const OP_BITOP_VX_EQUALS_VX_XOR_VY = 16
const OP_BITOP_VX_EQUALS_VX_RSHIFT = 17
const OP_BITOP_VX_EQUALS_VX_LSHIFT = 18

const OP_MATH_VX_EQUALS_VX_PLUS_VY  = 19
const OP_MATH_VX_EQUALS_VX_MINUS_VY = 20
const OP_MATH_VX_EQUALS_VY_MINUS_VX = 21

const OP_MEM_I_EQUALS_NNN = 22
const OP_MEM_I_EQUALS_I_PLUS_VX = 23
const OP_MEM_JMP_PC_EQUALS_V0_PLUS_NNN = 24
const OP_MEM_I_EQUALS_SPRITE_MAP_CHAR_VX = 25
const OP_MEM_STORE_REGS_AT_I_PTR = 26
const OP_MEM_LOAD_REGS_FROM_I_PTR = 27

const OP_RAND = 28

const OP_KEY_SKIPNXT_EQ_VX = 29
const OP_KEY_SKIPNXT_NEQ_VX = 30
const OP_KEY_VX_EQUALS_KEYPRESS= 31

const OP_VX_EQUALS_DELAY_TIMER = 32
const OP_DELAY_TIMER_EQUALS_VX = 33
const OP_SOUND_TIMER = 34
const OP_SOUND_TIMER_EQUALS_VX = 35

const OP_BCD_VX = 36

const FONT_BYTE_LEN = 5

const VF_REG = 0xF

type DecodedOpcode struct {
	Opcode int 
	Vx uint8 
	Vy uint8 
	N uint8 
	NN uint8 
	NNN uint16 
}