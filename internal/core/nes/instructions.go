package nes

import "fmt"

// instruction is a 6502 instruction.  It has a specific
// name, addressing mode, cycle cost, page cross cost, and byte cost.
type instruction struct {
	name               string
	addressingMode     int
	byteCost           int
	cycleCost          int
	pageCrossCycleCost int
	execute            func(uint16) // contains instruction logic
}

// ErrInvalidOpcode is an invalid opcode error.
// Only official opcodes are supported.
type ErrInvalidOpcode byte

// Error() implements error.
func (e ErrInvalidOpcode) Error() (repr string) {
	return fmt.Sprintf("invalid opcode: %v", byte(e))
}

// IsInvalidOpcodeErr returns whether or not err is of type
// ErrInvalidOpcode.
func IsInvalidOpcodeErr(err error) (invalid bool) {
	_, ok := err.(ErrInvalidOpcode)
	return ok
}

// decode decodes opcode opcode and returns relevant information.
func (c *cpu) decode(opcode byte) (instr *instruction, err error) {
	if instruction, ok := c.instructions[opcode]; ok {
		return instruction, nil
	}
	return &instruction{}, ErrInvalidOpcode(opcode)
}

// initInstructionLookupTable assembles instructions according to information
// from http://obelisk.me.uk/6502/reference.html.
// Instruction "execute" functions are assigned to c, i.e. set to make
// use of the memory and registers assigned to c.
func (c *cpu) initInstructionLookupTable() {
	c.instructions = map[byte]*instruction{
		// example
		0x69: {
			name:               "ADC",
			addressingMode:     modeImmediate,
			byteCost:           2,
			cycleCost:          2,
			pageCrossCycleCost: 0,
			execute:            c.ADC,
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
			c.ASLA,
		},
		0x06: {
			"ASL",
			modeZeroPage,
			2,
			5,
			0,
			c.ASLM,
		},
		0x16: {
			"ASL",
			modeZeroPageX,
			2,
			6,
			0,
			c.ASLM,
		},
		0x0E: {
			"ASL",
			modeAbsolute,
			3,
			6,
			0,
			c.ASLM,
		},
		0x1E: {
			"ASL",
			modeAbsoluteX,
			3,
			7,
			0,
			c.ASLM,
		},
		0x90: {
			"BCC",
			modeRelative,
			2,
			2,
			1,
			c.BCC,
		},
		0xB0: {
			"BCS",
			modeRelative,
			2,
			2,
			1,
			c.BCS,
		},
		0xF0: {
			"BEQ",
			modeRelative,
			2,
			2,
			1,
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
		0x30: {
			"BMI",
			modeRelative,
			2,
			2,
			1,
			c.BMI,
		},
		0xD0: {
			"BNE",
			modeRelative,
			2,
			2,
			1,
			c.BNE,
		},
		0x10: {
			"BPL",
			modeRelative,
			2,
			2,
			1,
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
		0x50: {
			"BVC",
			modeRelative,
			2,
			2,
			1,
			c.BVC,
		},
		0x70: {
			"BVS",
			modeRelative,
			2,
			2,
			1,
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
			c.LSRA,
		},
		0x46: {
			"LSR",
			modeZeroPage,
			2,
			5,
			0,
			c.LSRM,
		},
		0x56: {
			"LSR",
			modeZeroPageX,
			2,
			6,
			0,
			c.LSRM,
		},
		0x4E: {
			"LSR",
			modeAbsolute,
			3,
			6,
			0,
			c.LSRM,
		},
		0x5E: {
			"LSR",
			modeAbsoluteX,
			3,
			7,
			0,
			c.LSRM,
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
			c.ROLA,
		},
		0x26: {
			"ROL",
			modeZeroPage,
			2,
			5,
			0,
			c.ROLM,
		},
		0x36: {
			"ROL",
			modeZeroPageX,
			2,
			6,
			0,
			c.ROLM,
		},
		0x2E: {
			"ROL",
			modeAbsolute,
			3,
			6,
			0,
			c.ROLM,
		},
		0x3E: {
			"ROL",
			modeAbsoluteX,
			3,
			7,
			0,
			c.ROLM,
		},
		0x6A: {
			"ROR",
			modeAccumulator,
			1,
			2,
			0,
			c.RORA,
		},
		0x66: {
			"ROR",
			modeZeroPage,
			2,
			5,
			0,
			c.RORM,
		},
		0x76: {
			"ROR",
			modeZeroPageX,
			2,
			6,
			0,
			c.RORM,
		},
		0x6E: {
			"ROR",
			modeAbsolute,
			3,
			6,
			0,
			c.RORM,
		},
		0x7E: {
			"ROR",
			modeAbsoluteX,
			3,
			7,
			0,
			c.RORM,
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

// adcSbcHelper provides common logic for both ADC and SBC.
// This works because sbc can be implemented by invoking adc with data bits inverted.
// See http://forums.nesdev.com/viewtopic.php?p=19080#19080.
func (c *cpu) adcSbcHelper(data byte) {
	carry := convert(c.status.c)
	temp := c.a + data + carry
	c.status.c = (int(data) + int(carry) + int(c.a)) > zeroPageEnd
	c.status.v = ((c.a ^ temp) & (data ^ temp) & mask7) != 0
	c.a = temp
	c.status.setZN(c.a)
}

/* Instructions Start */

// ADC Add with Carry
// Fairly complicated, see http://www.obelisk.me.uk/6502/reference.html#ADC.
func (c *cpu) ADC(address uint16) {
	c.adcSbcHelper(c.Read(address))
	c.setPageCrossed(address)
}

// AND Logical AND
func (c *cpu) AND(address uint16) {
	c.a &= c.Read(address)
	c.status.setZN(c.a)
	c.setPageCrossed(address)
}

// ASLA Arithmetic Shift Left, acting on Accumulator.
// ASL is separated into two functions here for implementation
// reasons; one function which is called when ASL is called in
// modeAccumulator, and the other is called when ASL is
// called in any other addressing mode.
func (c *cpu) ASLA(address uint16) {
	c.status.c = c.a&mask7 != 0
	c.a <<= 1
	c.status.setZN(c.a)
}

// ASLM Arithmetic Shift Left, acting on Memory.
// See explanation for ASLA
func (c *cpu) ASLM(address uint16) {
	val := c.Read(address)
	c.status.c = val&mask7 != 0
	val <<= 1
	c.write(address, val)
	c.status.setZN(val)
}

// BCC Branch if Carry Clear
func (c *cpu) BCC(address uint16) {
	if !c.status.c {
		c.branchTo(address)
	}
}

// BCS Branch if Carry Set
func (c *cpu) BCS(address uint16) {
	if c.status.c {
		c.branchTo(address)
	}
}

// BEQ Branch if Equal
func (c *cpu) BEQ(address uint16) {
	if c.status.z {
		c.branchTo(address)
	}
}

// BIT Bit Test
func (c *cpu) BIT(address uint16) {
	val := c.Read(address)
	c.status.z = c.a&val == 0
	c.status.v = val&mask6 != 0
	c.status.n = val&mask7 != 0
}

// BMI Branch if Minus
func (c *cpu) BMI(address uint16) {
	if c.status.n {
		c.branchTo(address)
	}
}

// BNE Branch if Not Equal
func (c *cpu) BNE(address uint16) {
	if !c.status.z {
		c.branchTo(address)
	}
}

// BPL Branch if Positive
func (c *cpu) BPL(address uint16) {
	if !c.status.n {
		c.branchTo(address)
	}
}

// BRK Force Interrupt
func (c *cpu) BRK(address uint16) {
	c.GenerateInterrupt(irq)
	c.status.b = true
}

// BVC Branch if Overflow Clear
func (c *cpu) BVC(address uint16) {
	if !c.status.v {
		c.branchTo(address)
	}
}

// BVS Branch if Overflow Set
func (c *cpu) BVS(address uint16) {
	if c.status.v {
		c.branchTo(address)
	}
}

// CLC Clear Carry Flag
func (c *cpu) CLC(address uint16) {
	c.status.c = false
}

// CLD Clear Decimal Mode
func (c *cpu) CLD(address uint16) {
	c.status.d = false
}

// CLI Clear Interrupt Disable
func (c *cpu) CLI(address uint16) {
	c.status.i = false
}

// CLV Clear Overflow Flag
func (c *cpu) CLV(address uint16) {
	c.status.v = false
}

// CMP Compare
func (c *cpu) CMP(address uint16) {
	val := c.Read(address)
	c.status.setZN(c.a - val)
	c.status.c = c.a >= val
	c.setPageCrossed(address)
}

// CPX Compare X Register
func (c *cpu) CPX(address uint16) {
	val := c.Read(address)
	c.status.setZN(c.x - val)
	c.status.c = c.x >= val
}

// CPY Compare Y Register
func (c *cpu) CPY(address uint16) {
	val := c.Read(address)
	c.status.setZN(c.y - val)
	c.status.c = c.y >= val
}

// DEC Decrement Memory
func (c *cpu) DEC(address uint16) {
	val := c.Read(address) - 1
	c.write(address, val)
	c.status.setZN(val)
}

// DEX Decrement X Register
func (c *cpu) DEX(address uint16) {
	c.x--
	c.status.setZN(c.x)
}

// DEY Decrement Y Register
func (c *cpu) DEY(address uint16) {
	c.y--
	c.status.setZN(c.y)
}

// EOR Exclusive OR
func (c *cpu) EOR(address uint16) {
	c.a ^= c.Read(address)
	c.status.setZN(c.a)
	c.setPageCrossed(address)
}

// INC Increment Register
func (c *cpu) INC(address uint16) {
	newVal := c.Read(address) + 0x01
	c.write(address, newVal)
	c.status.setZN(newVal)
}

// INX Increment X Register
func (c *cpu) INX(address uint16) {
	c.x++
	c.status.setZN(c.x)
}

// INY Increment Y Register
func (c *cpu) INY(address uint16) {
	c.y++
	c.status.setZN(c.y)
}

// JMP Jump
// TODO: original 6502 bug?
func (c *cpu) JMP(address uint16) {
	c.pc = address
}

// JSR Jump to Subroutine
func (c *cpu) JSR(address uint16) {
	c.push16(c.pc)
	c.pc = address
}

// LDA Load Accumulator
func (c *cpu) LDA(address uint16) {
	c.a = c.Read(address)
	c.status.setZN(c.a)
	c.setPageCrossed(address)
}

// LDX Load X Register
func (c *cpu) LDX(address uint16) {
	c.x = c.Read(address)
	c.status.setZN(c.x)
	c.setPageCrossed(address)
}

// LDY Load Y Register
func (c *cpu) LDY(address uint16) {
	c.y = c.Read(address)
	c.status.setZN(c.y)
	c.setPageCrossed(address)
}

// LSRA Logical Shift Right, acting on Accumulator.
// LSR is separated into two functions here for implementation
// reasons; one function which is called when LSR is called in
// modeAccumulator, and the other is called when LSR is
// called in any other addressing mode.
func (c *cpu) LSRA(address uint16) {
	c.status.c = c.a&mask0 == 1
	c.a >>= 1
	c.status.setZN(c.a)
}

// LSRM Logical Shift Right, acting on Memory.
// See explanation for LSRA.
func (c *cpu) LSRM(address uint16) {
	val := c.Read(address)
	c.status.c = val&mask0 == 1
	val >>= 1
	c.write(address, val)
	c.status.setZN(val)
}

// NOP No Operation
func (c *cpu) NOP(address uint16) {
}

// ORA Logical Inclusive OR
func (c *cpu) ORA(address uint16) {
	c.a |= c.Read(address)
	c.status.setZN(c.a)
	c.setPageCrossed(address)
}

// PHA Push Accumulator
func (c *cpu) PHA(address uint16) {
	c.pushStack(c.a)
}

// PHP Push Processor Status
func (c *cpu) PHP(address uint16) {
	c.pushStack(c.status.asByte())
}

// PLA Pull Accumulator
func (c *cpu) PLA(address uint16) {
	c.a = c.pullStack()
	c.status.setZN(c.a)
}

// PLP Pull Processor Status
func (c *cpu) PLP(address uint16) {
	c.status.fromByte(c.pullStack())
}

// ROLA Rotate Left, acting on Accumulator.
// ROL is separated into two functions here for implementation
// reasons; one function which is called when ROL is called in
// modeAccumulator, and the other is called when ROL is
// called in any other addressing mode.
func (c *cpu) ROLA(address uint16) {
	carry := c.a&mask7 != 0
	c.a <<= 1
	if c.status.c {
		c.a |= 0x01
	} else {
		c.a &= 0xFE
	}
	c.status.c = carry
	c.status.setZN(c.a)
}

// ROLM Rotate Left, acting on Memory.
// See explanation for ROLA.
func (c *cpu) ROLM(address uint16) {
	val := c.Read(address)
	carry := val&mask7 != 0
	val <<= 1
	if c.status.c {
		val |= 0x01
	} else {
		val &= 0xFE
	}
	c.status.c = carry
	c.write(address, val)
	c.status.setZN(val)
}

// RORA Rotate Right, acting on Accumulator.
// ROR is separated into two functions here for implementation
// reasons; one function which is called when ROR is called in
// modeAccumulator, and the other is called when ROR is
// called in any other addressing mode.
func (c *cpu) RORA(address uint16) {
	carry := c.a&mask0 == 1
	c.a >>= 1
	if c.status.c {
		c.a |= 0x80
	} else {
		c.a &= 0x7F
	}
	c.status.c = carry
	c.status.setZN(c.a)
}

// RORM Rotate Right, acting on Memory
// See explanation for RORA
func (c *cpu) RORM(address uint16) {
	val := c.Read(address)
	carry := val&mask0 == 1
	val >>= 1
	if c.status.c {
		val |= 0x80
	} else {
		val &= 0x7F
	}
	c.status.c = carry
	c.write(address, val)
	c.status.setZN(val)
}

// RTI Return from Interrupt
func (c *cpu) RTI(address uint16) {
	c.status.fromByte(c.pullStack())
	c.pc = c.pull16()
	c.status.b = false
	c.status.i = false
}

// RTS Return from Subroutine
func (c *cpu) RTS(address uint16) {
	c.pc = c.pull16()
}

// SBC Subtract with Carry
// Fairly complicated, see http://www.obelisk.me.uk/6502/reference.html#SDC.
func (c *cpu) SBC(address uint16) {
	c.adcSbcHelper(c.Read(address) ^ zeroPageEnd)
	c.setPageCrossed(address)
}

// SEC Set Carry Flag
func (c *cpu) SEC(address uint16) {
	c.status.c = true
}

// SED Set Decimal Mode
func (c *cpu) SED(address uint16) {
	c.status.d = true
}

// SEI Set Interrupt Disable
func (c *cpu) SEI(address uint16) {
	c.status.i = true
}

// STA Store Accumulator
func (c *cpu) STA(address uint16) {
	c.write(address, c.a)
}

// STX Store X Register
func (c *cpu) STX(address uint16) {
	c.write(address, c.x)
}

// STY Store Y Register
func (c *cpu) STY(address uint16) {
	c.write(address, c.y)
}

// TAX Transfer Accumulator to X
func (c *cpu) TAX(address uint16) {
	c.x = c.a
	c.status.setZN(c.x)
}

// TAY Transfer Accumulator to Y
func (c *cpu) TAY(address uint16) {
	c.y = c.a
	c.status.setZN(c.y)
}

// TSX Transfer Stack Pointer to X
func (c *cpu) TSX(address uint16) {
	c.x = c.sp
	c.status.setZN(c.x)
}

// TXA Transfer X to Accumulator
func (c *cpu) TXA(address uint16) {
	c.a = c.x
	c.status.setZN(c.a)
}

// TXS Transfer X to Stack Pointer
func (c *cpu) TXS(address uint16) {
	c.sp = c.x
}

// TYA Transfer Y to Accumulator
func (c *cpu) TYA(address uint16) {
	c.a = c.y
	c.status.setZN(c.a)
}
