package cpu

import "fmt"

const (
	modeImplied = iota
	modeRelative
	modeAccumulator
	modeImmediate
	modeZeroPage
	modeZeroPageX
	modeZeroPageY
	modeAbsolute
	modeAbsoluteX
	modeAbsoluteY
	modeIndirect
	modeIndirectX
	modeIndirectY
)

type instruction struct {
	name           string
	addressingMode int
	cycleCost      int
	pageCrossCost  int
	byteCost       int
	execute        func(uint16)
}

// ErrInvalidOpcode is an invalid opcode error.
// Only official opcodes are supported.
type ErrInvalidOpcode byte

// Error() implements error
func (e ErrInvalidOpcode) Error() (repr string) {
	return fmt.Sprintf("invalid opcode: %v", byte(e))
}

// IsInvalidOpcodeErr returns whether or not err is of type
// ErrInvalidOpcode.
func IsInvalidOpcodeErr(err error) (invalid bool) {
	if _, ok := err.(ErrInvalidOpcode); ok {
		return true
	}
	return false
}

// initInstructions assembles instructions according to information
// from http://obelisk.me.uk/6502/reference.html.
// Instruction "execute" functions are assigned to c, i.e. set to make
// use of the memory and registers assigned to c.
func (c *CPU) initInstructions() {
	c.instructions = map[byte]instruction{
		// example
		0x69: {
			name:           "ADC",
			addressingMode: modeImmediate,
			byteCost:       2,
			cycleCost:      2,
			pageCrossCost:  0,
			execute:        c.ADC,
		},
		0x65: {
			"ADC",
			modeZeroPage,
			2,
			3,
			0,
			c.ADC,
		},
		0x75: {
			"ADC",
			modeZeroPageX,
			2,
			4,
			0,
			c.ADC,
		},
		0x6D: {
			"ADC",
			modeAbsolute,
			3,
			4,
			0,
			c.ADC,
		},
		0x7D: {
			"ADC",
			modeAbsoluteX,
			3,
			4,
			1,
			c.ADC,
		},
		0x79: {
			"ADC",
			modeAbsoluteY,
			3,
			4,
			1,
			c.ADC,
		},
		0x61: {
			"ADC",
			modeIndirectX,
			2,
			6,
			0,
			c.ADC,
		},
		0x71: {
			"ADC",
			modeIndirectY,
			2,
			5,
			1,
			c.ADC,
		},
		0x29: {
			"AND",
			modeImmediate,
			2,
			2,
			0,
			c.AND,
		},
		0x25: {
			"AND",
			modeZeroPage,
			2,
			3,
			0,
			c.AND,
		},
		0x35: {
			"AND",
			modeZeroPageX,
			2,
			4,
			0,
			c.AND,
		},
		0x2D: {
			"AND",
			modeAbsolute,
			3,
			4,
			0,
			c.AND,
		},
		0x3D: {
			"AND",
			modeAbsoluteX,
			3,
			4,
			1,
			c.AND,
		},
		0x39: {
			"AND",
			modeAbsoluteY,
			3,
			4,
			1,
			c.AND,
		},
		0x21: {
			"AND",
			modeIndirectX,
			2,
			6,
			0,
			c.AND,
		},
		0x31: {
			"AND",
			modeIndirectY,
			2,
			5,
			1,
			c.AND,
		},
		0x0A: {
			"ASL",
			modeAccumulator,
			1,
			2,
			0,
			c.ASL,
		},
		0x06: {
			"ASL",
			modeZeroPage,
			2,
			5,
			0,
			c.ASL,
		},
		0x16: {
			"ASL",
			modeZeroPageX,
			2,
			6,
			0,
			c.ASL,
		},
		0x0E: {
			"ASL",
			modeAbsolute,
			3,
			6,
			0,
			c.ASL,
		},
		0x1E: {
			"ASL",
			modeAbsoluteX,
			3,
			7,
			0,
			c.ASL,
		},
		// TODO: special case
		0x90: {
			"BCC",
			modeRelative,
			2,
			2,
			0,
			c.BCC,
		},
		// TODO: special case
		0xB0: {
			"BCS",
			modeRelative,
			2,
			2,
			0,
			c.BCS,
		},
		// TODO: special case
		0xF0: {
			"BEQ",
			modeRelative,
			2,
			2,
			0,
			c.BEQ,
		},
		0x24: {
			"BIT",
			modeZeroPage,
			2,
			3,
			0,
			c.BIT,
		},
		0x2C: {
			"BIT",
			modeAbsolute,
			3,
			4,
			0,
			c.BIT,
		},
		// TODO: special case
		0x30: {
			"BMI",
			modeRelative,
			2,
			2,
			0,
			c.BMI,
		},
		// TODO: special case
		0xD0: {
			"BNE",
			modeRelative,
			2,
			2,
			0,
			c.BNE,
		},
		// TODO: special case
		0x10: {
			"BPL",
			modeRelative,
			2,
			2,
			0,
			c.BPL,
		},
		0x00: {
			"BRK",
			modeImplied,
			1,
			7,
			0,
			c.BRK,
		},
		// TODO: special case
		0x50: {
			"BVC",
			modeRelative,
			2,
			2,
			0,
			c.BVC,
		},
		// TODO: special case
		0x70: {
			"BVS",
			modeRelative,
			2,
			2,
			0,
			c.BVS,
		},
		0x18: {
			"CLC",
			modeImplied,
			1,
			2,
			0,
			c.CLC,
		},
		0xD8: {
			"CLD",
			modeImplied,
			1,
			2,
			0,
			c.CLD,
		},
		0x58: {
			"CLI",
			modeImplied,
			1,
			2,
			0,
			c.CLI,
		},
		0xB8: {
			"CLV",
			modeImplied,
			1,
			2,
			0,
			c.CLV,
		},
		0xC9: {
			"CMP",
			modeImmediate,
			2,
			2,
			0,
			c.CMP,
		},
		0xC5: {
			"CMP",
			modeZeroPage,
			2,
			3,
			0,
			c.CMP,
		},
		0xD5: {
			"CMP",
			modeZeroPageX,
			2,
			4,
			0,
			c.CMP,
		},
		0xCD: {
			"CMP",
			modeAbsolute,
			3,
			4,
			0,
			c.CMP,
		},
		0xDD: {
			"CMP",
			modeAbsoluteX,
			3,
			4,
			1,
			c.CMP,
		},
		0xD9: {
			"CMP",
			modeAbsoluteY,
			3,
			4,
			1,
			c.CMP,
		},
		0xC1: {
			"CMP",
			modeIndirectX,
			2,
			6,
			0,
			c.CMP,
		},
		0xD1: {
			"CMP",
			modeIndirectY,
			2,
			5,
			1,
			c.CMP,
		},
		0xE0: {
			"CPX",
			modeImmediate,
			2,
			2,
			0,
			c.CPX,
		},
		0xE4: {
			"CPX",
			modeZeroPage,
			2,
			3,
			0,
			c.CPX,
		},
		0xEC: {
			"CPX",
			modeAbsolute,
			3,
			4,
			0,
			c.CPX,
		},
		0xC0: {
			"CPY",
			modeImmediate,
			2,
			2,
			0,
			c.CPY,
		},
		0xC4: {
			"CPY",
			modeZeroPage,
			2,
			3,
			0,
			c.CPY,
		},
		0xCC: {
			"CPY",
			modeAbsolute,
			3,
			4,
			0,
			c.CPY,
		},
		0xC6: {
			"DEC",
			modeZeroPage,
			2,
			5,
			0,
			c.DEC,
		},
		0xD6: {
			"DEC",
			modeZeroPageX,
			2,
			6,
			0,
			c.DEC,
		},
		0xCE: {
			"DEC",
			modeAbsolute,
			3,
			6,
			0,
			c.DEC,
		},
		0xDE: {
			"DEC",
			modeAbsoluteX,
			3,
			7,
			0,
			c.DEC,
		},
		0xCA: {
			"DEX",
			modeImplied,
			1,
			2,
			0,
			c.DEX,
		},
		0x88: {
			"DEY",
			modeImplied,
			1,
			2,
			0,
			c.DEY,
		},
		0x49: {
			"EOR",
			modeImmediate,
			2,
			2,
			0,
			c.EOR,
		},
		0x45: {
			"EOR",
			modeZeroPage,
			2,
			3,
			0,
			c.EOR,
		},
		0x55: {
			"EOR",
			modeZeroPageX,
			2,
			4,
			0,
			c.EOR,
		},
		0x4D: {
			"EOR",
			modeAbsolute,
			3,
			4,
			0,
			c.EOR,
		},
		0x5D: {
			"EOR",
			modeAbsoluteX,
			3,
			4,
			1,
			c.EOR,
		},
		0x59: {
			"EOR",
			modeAbsoluteY,
			3,
			4,
			1,
			c.EOR,
		},
		0x41: {
			"EOR",
			modeIndirectX,
			2,
			6,
			0,
			c.EOR,
		},
		0x51: {
			"EOR",
			modeIndirectY,
			2,
			5,
			1,
			c.EOR,
		},
		0xE6: {
			"INC",
			modeZeroPage,
			2,
			5,
			0,
			c.INC,
		},
		0xF6: {
			"INC",
			modeZeroPageX,
			2,
			6,
			0,
			c.INC,
		},
		0xEE: {
			"INC",
			modeAbsolute,
			3,
			6,
			0,
			c.INC,
		},
		0xFE: {
			"INC",
			modeAbsoluteX,
			3,
			7,
			0,
			c.INC,
		},
		0xE8: {
			"INX",
			modeImplied,
			1,
			2,
			0,
			c.INX,
		},
		0xC8: {
			"INY",
			modeImplied,
			1,
			2,
			0,
			c.INY,
		},
		0x4C: {
			"JMP",
			modeAbsolute,
			3,
			3,
			0,
			c.JMP,
		},
		0x6C: {
			"JMP",
			modeIndirect,
			3,
			5,
			0,
			c.JMP,
		},
		0x20: {
			"JSR",
			modeAbsolute,
			3,
			6,
			0,
			c.JSR,
		},
		0xA9: {
			"LDA",
			modeImmediate,
			2,
			2,
			0,
			c.LDA,
		},
		0xA5: {
			"LDA",
			modeZeroPage,
			2,
			3,
			0,
			c.LDA,
		},
		0xB5: {
			"LDA",
			modeZeroPageX,
			2,
			4,
			0,
			c.LDA,
		},
		0xAD: {
			"LDA",
			modeAbsolute,
			3,
			4,
			0,
			c.LDA,
		},
		0xBD: {
			"LDA",
			modeAbsoluteX,
			3,
			4,
			1,
			c.LDA,
		},
		0xB9: {
			"LDA",
			modeAbsoluteY,
			3,
			4,
			1,
			c.LDA,
		},
		0xA1: {
			"LDA",
			modeIndirectX,
			2,
			6,
			0,
			c.LDA,
		},
		0xB1: {
			"LDA",
			modeIndirectY,
			2,
			5,
			1,
			c.LDA,
		},
		0xA2: {
			"LDX",
			modeImmediate,
			2,
			2,
			0,
			c.LDX,
		},
		0xA6: {
			"LDX",
			modeZeroPage,
			2,
			3,
			0,
			c.LDX,
		},
		0xB6: {
			"LDX",
			modeZeroPageY,
			2,
			4,
			0,
			c.LDX,
		},
		0xAE: {
			"LDX",
			modeAbsolute,
			3,
			4,
			0,
			c.LDX,
		},
		0xBE: {
			"LDX",
			modeAbsoluteY,
			3,
			4,
			1,
			c.LDX,
		},
		0xA0: {
			"LDY",
			modeImmediate,
			2,
			2,
			0,
			c.LDY,
		},
		0xA4: {
			"LDY",
			modeZeroPage,
			2,
			3,
			0,
			c.LDY,
		},
		0xB4: {
			"LDY",
			modeZeroPageX,
			2,
			4,
			0,
			c.LDY,
		},
		0xAC: {
			"LDY",
			modeAbsolute,
			3,
			4,
			0,
			c.LDY,
		},
		0xBC: {
			"LDY",
			modeAbsoluteX,
			3,
			4,
			1,
			c.LDY,
		},
		0x4A: {
			"LSR",
			modeAccumulator,
			1,
			2,
			0,
			c.LSR,
		},
		0x46: {
			"LSR",
			modeZeroPage,
			2,
			5,
			0,
			c.LSR,
		},
		0x56: {
			"LSR",
			modeZeroPageX,
			2,
			6,
			0,
			c.LSR,
		},
		0x4E: {
			"LSR",
			modeAbsolute,
			3,
			6,
			0,
			c.LSR,
		},
		0x5E: {
			"LSR",
			modeAbsoluteX,
			3,
			7,
			0,
			c.LSR,
		},
		0xEA: {
			"NOP",
			modeImplied,
			1,
			2,
			0,
			c.NOP,
		},
		0x09: {
			"ORA",
			modeImmediate,
			2,
			2,
			0,
			c.ORA,
		},
		0x05: {
			"ORA",
			modeZeroPage,
			2,
			3,
			0,
			c.ORA,
		},
		0x15: {
			"ORA",
			modeZeroPageX,
			2,
			4,
			0,
			c.ORA,
		},
		0x0D: {
			"ORA",
			modeAbsolute,
			3,
			4,
			0,
			c.ORA,
		},
		0x1D: {
			"ORA",
			modeAbsoluteX,
			3,
			4,
			1,
			c.ORA,
		},
		0x19: {
			"ORA",
			modeAbsoluteY,
			3,
			4,
			1,
			c.ORA,
		},
		0x01: {
			"ORA",
			modeIndirectX,
			2,
			6,
			0,
			c.ORA,
		},
		0x11: {
			"ORA",
			modeIndirectY,
			2,
			5,
			1,
			c.ORA,
		},
		0x48: {
			"PHA",
			modeImplied,
			1,
			3,
			0,
			c.PHA,
		},
		0x08: {
			"PHP",
			modeImplied,
			1,
			3,
			0,
			c.PHP,
		},
		0x68: {
			"PLA",
			modeImplied,
			1,
			4,
			0,
			c.PLA,
		},
		0x28: {
			"PLP",
			modeImplied,
			1,
			4,
			0,
			c.PLP,
		},
		0x2A: {
			"ROL",
			modeAccumulator,
			1,
			2,
			0,
			c.ROL,
		},

		0x26: {
			"ROL",
			modeZeroPage,
			2,
			5,
			0,
			c.ROL,
		},
		0x36: {
			"ROL",
			modeZeroPageX,
			2,
			6,
			0,
			c.ROL,
		},
		0x2E: {
			"ROL",
			modeAbsolute,
			3,
			6,
			0,
			c.ROL,
		},
		0x3E: {
			"ROL",
			modeAbsoluteX,
			3,
			7,
			0,
			c.ROL,
		},
		0x6A: {
			"ROR",
			modeAccumulator,
			1,
			2,
			0,
			c.ROR,
		},

		0x66: {
			"ROR",
			modeZeroPage,
			2,
			5,
			0,
			c.ROR,
		},
		0x76: {
			"ROR",
			modeZeroPageX,
			2,
			6,
			0,
			c.ROR,
		},
		0x6E: {
			"ROR",
			modeAbsolute,
			3,
			6,
			0,
			c.ROR,
		},
		0x7E: {
			"ROR",
			modeAbsoluteX,
			3,
			7,
			0,
			c.ROR,
		},
		0x40: {
			"RTI",
			modeImplied,
			1,
			6,
			0,
			c.RTI,
		},
		0x60: {
			"RTS",
			modeImplied,
			1,
			6,
			0,
			c.RTS,
		},
		0xE9: {
			"SBC",
			modeImmediate,
			2,
			2,
			0,
			c.SBC,
		},
		0xE5: {
			"SBC",
			modeZeroPage,
			2,
			3,
			0,
			c.SBC,
		},
		0xF5: {
			"SBC",
			modeZeroPageX,
			2,
			4,
			0,
			c.SBC,
		},
		0xED: {
			"SBC",
			modeAbsolute,
			3,
			4,
			0,
			c.SBC,
		},
		0xFD: {
			"SBC",
			modeAbsoluteX,
			3,
			4,
			1,
			c.SBC,
		},
		0xF9: {
			"SBC",
			modeAbsoluteY,
			3,
			4,
			1,
			c.SBC,
		},
		0xE1: {
			"SBC",
			modeIndirectX,
			2,
			6,
			0,
			c.SBC,
		},
		0xF1: {
			"SBC",
			modeIndirectY,
			2,
			5,
			1,
			c.SBC,
		},
		0x38: {
			"SEC",
			modeImplied,
			1,
			2,
			0,
			c.SEC,
		},
		0xF8: {
			"SED",
			modeImplied,
			1,
			2,
			0,
			c.SED,
		},
		0x78: {
			"SEI",
			modeImplied,
			1,
			2,
			0,
			c.SEI,
		},
		0x85: {
			"STA",
			modeZeroPage,
			2,
			3,
			0,
			c.STA,
		},
		0x95: {
			"STA",
			modeZeroPageX,
			2,
			4,
			0,
			c.STA,
		},
		0x8D: {
			"STA",
			modeAbsolute,
			3,
			4,
			0,
			c.STA,
		},
		0x9D: {
			"STA",
			modeAbsoluteX,
			3,
			5,
			0,
			c.STA,
		},
		0x99: {
			"STA",
			modeAbsoluteY,
			3,
			5,
			0,
			c.STA,
		},
		0x81: {
			"STA",
			modeIndirectX,
			2,
			6,
			0,
			c.STA,
		},
		0x91: {
			"STA",
			modeIndirectY,
			2,
			6,
			0,
			c.STA,
		},
		0x86: {
			"STX",
			modeZeroPage,
			2,
			3,
			0,
			c.STX,
		},
		0x96: {
			"STX",
			modeZeroPageY,
			2,
			4,
			0,
			c.STX,
		},
		0x8E: {
			"STX",
			modeAbsolute,
			3,
			4,
			0,
			c.STX,
		},
		0x84: {
			"STY",
			modeZeroPage,
			2,
			3,
			0,
			c.STY,
		},
		0x94: {
			"STY",
			modeZeroPageX,
			2,
			4,
			0,
			c.STY,
		},
		0x8C: {
			"STY",
			modeAbsolute,
			3,
			4,
			0,
			c.STY,
		},
		0xAA: {
			"TAX",
			modeImplied,
			1,
			2,
			0,
			c.TAX,
		},
		0xA8: {
			"TAY",
			modeImplied,
			1,
			2,
			0,
			c.TAY,
		},
		0xBA: {
			"TSX",
			modeImplied,
			1,
			2,
			0,
			c.TSX,
		},
		0x8A: {
			"TXA",
			modeImplied,
			1,
			2,
			0,
			c.TXA,
		},
		0x9A: {
			"TXS",
			modeImplied,
			1,
			2,
			0,
			c.TXS,
		},
		0x98: {
			"TYA",
			modeImplied,
			1,
			2,
			0,
			c.TYA,
		},
	}
}

// ADC Add with Carry
func (c *CPU) ADC(address uint16) {
}

// AND Logical AND
func (c *CPU) AND(address uint16) {
}

// ASL Arithmetic Shift Left
func (c *CPU) ASL(address uint16) {
}

// BCC Branch if Carry Clear
func (c *CPU) BCC(address uint16) {
}

// BCS Branch if Carry Set
func (c *CPU) BCS(address uint16) {
}

// BEQ Branch if Equal
func (c *CPU) BEQ(address uint16) {
}

// BIT Bit Test
func (c *CPU) BIT(address uint16) {
}

// BMI Branch if Minus
func (c *CPU) BMI(address uint16) {
}

// BNE Branch if Not Equal
func (c *CPU) BNE(address uint16) {
}

// BPL Branch if Positive
func (c *CPU) BPL(address uint16) {
}

// BRK Force Interrupt
func (c *CPU) BRK(address uint16) {
}

// BVC Branch if Overflow Clear
func (c *CPU) BVC(address uint16) {
}

// BVS Branch if Overflow Set
func (c *CPU) BVS(address uint16) {
}

// CLC Clear Carry Flag
func (c *CPU) CLC(address uint16) {
	c.Status.C = false
}

// CLD Clear Decimal Mode
func (c *CPU) CLD(address uint16) {
	c.Status.D = false
}

// CLI Clear Interrupt Disable
func (c *CPU) CLI(address uint16) {
	c.Status.I = false
}

// CLV Clear Overflow Flag
func (c *CPU) CLV(address uint16) {
	c.Status.V = false
}

// CMP Compare
func (c *CPU) CMP(address uint16) {
}

// CPX Compare X Register
func (c *CPU) CPX(address uint16) {
}

// CPY Compare Y Register
func (c *CPU) CPY(address uint16) {
}

// DEC Decrement Memory
func (c *CPU) DEC(address uint16) {
}

// DEX Decrement X Register
func (c *CPU) DEX(address uint16) {
}

// DEY Decrement Y Register
func (c *CPU) DEY(address uint16) {
}

// EOR Exclusive OR
func (c *CPU) EOR(address uint16) {
}

// INC Increment Register
func (c *CPU) INC(address uint16) {
}

// INX Increment X Register
func (c *CPU) INX(address uint16) {
}

// INY Increment Y Register
func (c *CPU) INY(address uint16) {
}

// JMP Jump
func (c *CPU) JMP(address uint16) {
}

// JSR Jump to Subroutine
func (c *CPU) JSR(address uint16) {
}

// LDA Load Accumulator
func (c *CPU) LDA(address uint16) {
}

// LDX Load X Register
func (c *CPU) LDX(address uint16) {
}

// LDY Load Y Register
func (c *CPU) LDY(address uint16) {
}

// LSR Logical Shift Right
func (c *CPU) LSR(address uint16) {
}

// NOP No Operation
func (c *CPU) NOP(address uint16) {
}

// ORA Logical Inclusive OR
func (c *CPU) ORA(address uint16) {
}

// PHA Push Accumulator
func (c *CPU) PHA(address uint16) {
}

// PHP Push Processor Status
func (c *CPU) PHP(address uint16) {
}

// PLA Pull Accumulator
func (c *CPU) PLA(address uint16) {
}

// PLP Pull Processor Status
func (c *CPU) PLP(address uint16) {
}

// ROL Rotate Left
func (c *CPU) ROL(address uint16) {
}

// ROR Rotate Right
func (c *CPU) ROR(address uint16) {
}

// RTI Return from Interrupt
func (c *CPU) RTI(address uint16) {
}

// RTS Return from Subroutine
func (c *CPU) RTS(address uint16) {
}

// SBC Subtract with Carry
func (c *CPU) SBC(address uint16) {
}

// SEC Set Carry Flag
func (c *CPU) SEC(address uint16) {
	c.Status.C = true
}

// SED Set Decimal Mode
func (c *CPU) SED(address uint16) {
	c.Status.D = true
}

// SEI Set Interrupt Disable
func (c *CPU) SEI(address uint16) {
	c.Status.I = true
}

// STA Store Accumulator
func (c *CPU) STA(address uint16) {
}

// STX Store X Register
func (c *CPU) STX(address uint16) {
}

// STY Store Y Register
func (c *CPU) STY(address uint16) {
}

// TAX Transfer Accumulator to X
func (c *CPU) TAX(address uint16) {
}

// TAY Transfer Accumulator to Y
func (c *CPU) TAY(address uint16) {
}

// TSX Transfer Stack Pointer to X
func (c *CPU) TSX(address uint16) {
}

// TXA Transfer X to Accumulator
func (c *CPU) TXA(address uint16) {
}

// TXS Transfer X to Stack Pointer
func (c *CPU) TXS(address uint16) {
}

// TYA Transfer Y to Accumulator
func (c *CPU) TYA(address uint16) {
}
