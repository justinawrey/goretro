package cpu

const (
	modeImplied = iota
	modeRelative
	modeAccumulator
	modeImmediate
	modeZeroPage
	modeZeroPageX
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
	execute        func()
}

// initInstructions assembles instructions according to information
// from http://obelisk.me.uk/6502/reference.html.
// Instruction "execute" functions are assigned to c, i.e. set to make
// use of the memory and registers assigned to c.
func initInstructions(c *CPU) {
	c.instructions = map[byte]instruction{
		// example
		0x69: {
			name:           "ADC",
			addressingMode: modeImmediate,
			byteCost:       2,
			cycleCost:      2,
			pageCrossCost:  0,
			execute:        func() {},
		},
		0x65: {
			"ADC",
			modeZeroPage,
			2,
			3,
			0,
			func() {},
		},
		0x75: {
			"ADC",
			modeZeroPageX,
			2,
			4,
			0,
			func() {},
		},
		0x6D: {
			"ADC",
			modeAbsolute,
			3,
			4,
			0,
			func() {},
		},
		0x7D: {
			"ADC",
			modeAbsoluteX,
			3,
			4,
			1,
			func() {},
		},
		0x79: {
			"ADC",
			modeAbsoluteY,
			3,
			4,
			1,
			func() {},
		},
		0x61: {
			"ADC",
			modeIndirectX,
			2,
			6,
			0,
			func() {},
		},
		0x71: {
			"ADC",
			modeIndirectY,
			2,
			5,
			1,
			func() {},
		},
		0x29: {
			"AND",
			modeImmediate,
			2,
			2,
			0,
			func() {},
		},
		0x25: {
			"AND",
			modeZeroPage,
			2,
			3,
			0,
			func() {},
		},
		0x35: {
			"AND",
			modeZeroPageX,
			2,
			4,
			0,
			func() {},
		},
		0x2D: {
			"AND",
			modeAbsolute,
			3,
			4,
			0,
			func() {},
		},
		0x3D: {
			"AND",
			modeAbsoluteX,
			3,
			4,
			1,
			func() {},
		},
		0x39: {
			"AND",
			modeAbsoluteY,
			3,
			4,
			1,
			func() {},
		},
		0x21: {
			"AND",
			modeIndirectX,
			2,
			6,
			0,
			func() {},
		},
		0x31: {
			"AND",
			modeIndirectY,
			2,
			5,
			1,
			func() {},
		},
		0x0A: {
			"ASL",
			modeAccumulator,
			1,
			2,
			0,
			func() {},
		},
		0x06: {
			"ASL",
			modeZeroPage,
			2,
			5,
			0,
			func() {},
		},
		0x16: {
			"ASL",
			modeZeroPageX,
			2,
			6,
			0,
			func() {},
		},
		0x0E: {
			"ASL",
			modeAbsolute,
			3,
			6,
			0,
			func() {},
		},
		0x1E: {
			"ASL",
			modeAbsoluteX,
			3,
			7,
			0,
			func() {},
		},
		// TODO: special case
		0x90: {
			"BCC",
			modeRelative,
			2,
			2,
			0,
			func() {},
		},
		// TODO: special case
		0xB0: {
			"BCS",
			modeRelative,
			2,
			2,
			0,
			func() {},
		},
		// TODO: special case
		0xF0: {
			"BEQ",
			modeRelative,
			2,
			2,
			0,
			func() {},
		},
		0x24: {
			"BIT",
			modeZeroPage,
			2,
			3,
			0,
			func() {},
		},
		0x2C: {
			"BIT",
			modeAbsolute,
			3,
			4,
			0,
			func() {},
		},
		// TODO: special case
		0x30: {
			"BMI",
			modeRelative,
			2,
			2,
			0,
			func() {},
		},
		// TODO: special case
		0xD0: {
			"BNE",
			modeRelative,
			2,
			2,
			0,
			func() {},
		},
		// TODO: special case
		0x10: {
			"BPL",
			modeRelative,
			2,
			2,
			0,
			func() {},
		},
		0x00: {
			"BRK",
			modeImplied,
			1,
			7,
			0,
			func() {},
		},
		// TODO: special case
		0x50: {
			"BVC",
			modeRelative,
			2,
			2,
			0,
			func() {},
		},
		// TODO: special case
		0x70: {
			"BVS",
			modeRelative,
			2,
			2,
			0,
			func() {},
		},
		0x18: {
			"CLC",
			modeImplied,
			1,
			2,
			0,
			func() {},
		},
		0xD8: {
			"CLD",
			modeImplied,
			1,
			2,
			0,
			func() {},
		},
		0x58: {
			"CLI",
			modeImplied,
			1,
			2,
			0,
			func() {},
		},
		0xB8: {
			"CLV",
			modeImplied,
			1,
			2,
			0,
			func() {},
		},
		0xC9: {
			"CMP",
			modeImmediate,
			2,
			2,
			0,
			func() {},
		},
		0xC5: {
			"CMP",
			modeZeroPage,
			2,
			3,
			0,
			func() {},
		},
		0xD5: {
			"CMP",
			modeZeroPageX,
			2,
			4,
			0,
			func() {},
		},
		0xCD: {
			"CMP",
			modeAbsolute,
			3,
			4,
			0,
			func() {},
		},
		0xDD: {
			"CMP",
			modeAbsoluteX,
			3,
			4,
			1,
			func() {},
		},
		0xD9: {
			"CMP",
			modeAbsoluteY,
			3,
			4,
			1,
			func() {},
		},
		0xC1: {
			"CMP",
			modeIndirectX,
			2,
			6,
			0,
			func() {},
		},
		0xD1: {
			"CMP",
			modeIndirectY,
			2,
			5,
			1,
			func() {},
		},
		0xE0: {
			"CPX",
			modeImmediate,
			2,
			2,
			0,
			func() {},
		},
		0xE4: {
			"CPX",
			modeZeroPage,
			2,
			3,
			0,
			func() {},
		},
		0xEC: {
			"CPX",
			modeAbsolute,
			3,
			4,
			0,
			func() {},
		},
		0xC0: {
			"CPY",
			modeImmediate,
			2,
			2,
			0,
			func() {},
		},
		0xC4: {
			"CPY",
			modeZeroPage,
			2,
			3,
			0,
			func() {},
		},
		0xCC: {
			"CPY",
			modeAbsolute,
			3,
			4,
			0,
			func() {},
		},
		0xC6: {
			"DEC",
			modeZeroPage,
			2,
			5,
			0,
			func() {},
		},
		0xD6: {
			"DEC",
			modeZeroPageX,
			2,
			6,
			0,
			func() {},
		},
		0xCE: {
			"DEC",
			modeAbsolute,
			3,
			6,
			0,
			func() {},
		},
		0xDE: {
			"DEC",
			modeAbsoluteX,
			3,
			7,
			0,
			func() {},
		},
		0xCA: {
			"DEX",
			modeImplied,
			1,
			2,
			0,
			func() {},
		},
		0x88: {
			"DEY",
			modeImplied,
			1,
			2,
			0,
			func() {},
		},
		0x49: {
			"EOR",
			modeImmediate,
			2,
			2,
			0,
			func() {},
		},
		0x45: {
			"EOR",
			modeZeroPage,
			2,
			3,
			0,
			func() {},
		},
		0x55: {
			"EOR",
			modeZeroPageX,
			2,
			4,
			0,
			func() {},
		},
		0x4D: {
			"EOR",
			modeAbsolute,
			3,
			4,
			0,
			func() {},
		},
		0x5D: {
			"EOR",
			modeAbsoluteX,
			3,
			4,
			1,
			func() {},
		},
		0x59: {
			"EOR",
			modeAbsoluteY,
			3,
			4,
			1,
			func() {},
		},
		0x41: {
			"EOR",
			modeIndirectX,
			2,
			6,
			0,
			func() {},
		},
		0x51: {
			"EOR",
			modeIndirectY,
			2,
			5,
			1,
			func() {},
		},
		0xE6: {
			"INC",
			modeZeroPage,
			2,
			5,
			0,
			func() {},
		},
		0xF6: {
			"INC",
			modeZeroPageX,
			2,
			6,
			0,
			func() {},
		},
		0xEE: {
			"INC",
			modeAbsolute,
			3,
			6,
			0,
			func() {},
		},
		0xFE: {
			"INC",
			modeAbsoluteX,
			3,
			7,
			0,
			func() {},
		},
		0xE8: {
			"INX",
			modeImplied,
			1,
			2,
			0,
			func() {},
		},
		0xC8: {
			"INY",
			modeImplied,
			1,
			2,
			0,
			func() {},
		},
		0x4C: {
			"JMP",
			modeAbsolute,
			3,
			3,
			0,
			func() {},
		},
		0x6C: {
			"JMP",
			modeIndirect,
			3,
			5,
			0,
			func() {},
		},
	}
}

// SEC Set Carry Flag
func (cpu *CPU) SEC() {
	cpu.Status.C = true
}

// CLC Clear Carry Flag
func (cpu *CPU) CLC() {
	cpu.Status.C = false
}

// SEI Set Interrupt Disable
func (cpu *CPU) SEI() {
	cpu.Status.I = true
}

// CLI Clear Interrupt Disable
func (cpu *CPU) CLI() {
	cpu.Status.I = false
}

// SED Set Decimal Mode
func (cpu *CPU) SED() {
	cpu.Status.D = true
}

// CLD Clear Decimal Mode
func (cpu *CPU) CLD() {
	cpu.Status.D = false
}

// CLV Clear Overflow Flag
func (cpu *CPU) CLV() {
	cpu.Status.V = false
}
