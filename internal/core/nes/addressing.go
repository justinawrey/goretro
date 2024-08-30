package nes

// CPU addressing modes
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

// GetAddressWithMode uses addressing mode addressingMode to get
// an address on which any instruction can execute.
// Must be used when c.PC is on an opcode address, otherwise
// the following addresses will be interpreted incorrectly.
func (c *cpu) getAddressWithMode(addressingMode int) (addr uint16) {
	switch addressingMode {
	case modeImplied:
		// Address will be unused for following two addressing modes; return 0
		fallthrough

	case modeAccumulator:
		return 0

	case modeRelative:
		// Instructions with modeRelative take 2 bytes:
		// 1. opcode
		// 2. 8 bit constant value
		// The address will only be jumped to if the branch succeeeds.
		// Note: relative addressing uses twos complement to branch both
		// forwards and backwards.
		offset := uint16(c.Read(c.pc + 1))
		if offset >= 0x80 {
			// interpret as negative number
			return c.pc + offset - 0x100
		}
		return c.pc + offset

	case modeImmediate:
		// Instructions with modeImmediate take 2 bytes:
		// 1. opcode
		// 2. 8 bit constant value
		return c.pc + 1

	case modeZeroPage:
		// Instructions with modeZeroPage take 2 bytes:
		// 1. opcode
		// 2. zero-page address
		return uint16(c.Read(c.pc + 1))

	case modeZeroPageX:
		// Same as modeZeroPage, but with zero page address being added to X register with wraparound
		return uint16(c.Read(c.pc+1)+c.x) & zeroPageEnd

	case modeZeroPageY:
		// Same as modeZeroPage, but with zero page address being added to Y register with wraparound
		return uint16(c.Read(c.pc+1)+c.y) & zeroPageEnd

	case modeAbsolute:
		// Instructions with modeAbsolute take 3 bytes:
		// 1. opcode
		// 2. least significant byte of address
		// 3. most significant byte of address
		return c.read16(c.pc + 1)

	case modeAbsoluteX:
		// Same as modeAbsolute, with address being added to contents of X register
		return c.read16(c.pc+1) + uint16(c.x)

	case modeAbsoluteY:
		// Same as modeAbsolute, with address being added to contents of Y register
		return c.read16(c.pc+1) + uint16(c.y)

	case modeIndirect:
		// Instructions with modeIndirect take 3 bytes:
		// 1. opcode
		// 2. least significant byte of address
		// 3. most significant byte of address
		// The formulated address, along with the next,
		// are then accessed again to get the final address.
		return c.read16(c.read16(c.pc + 1))

	case modeIndirectX:
		// Instructions with modeIndirectX take 2 bytes:
		// 1. opcode
		// 2. single byte
		// The byte is then added to the X register, which then
		// gives the least significant byte of the target address.
		return c.read16(uint16(c.Read(c.pc+1) + c.x))

	case modeIndirectY:
		// Instructions with modeIndirectY take 2 bytes:
		// 1. opcode
		// 2. least significant byte of zero page address
		// The zero page address is then accessed, and the data
		// is added to the Y register. The resulting data is the
		// target address.
		return c.read16(uint16(c.Read(c.pc+1))) + uint16(c.y)

	default:
		// shouldn't happen, but handle gracefully
		return 0
	}
}
