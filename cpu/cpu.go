package cpu

import (
	"fmt"
	"log"

	"github.com/justinawrey/nes/memory"
)

// Status holds data for each status flag
// in the 6502 status register.
type Status struct {
	C bool // Carry
	Z bool // Zero result
	I bool // Interrupt disable
	D bool // Decimal mode
	B bool // Break command
	V bool // Overflow
	N bool // Zero result
}

// String implements Stringer
func (sr *Status) String() (repr string) {
	convert := func(bit bool) string {
		if bit {
			return "1"
		}
		return "0"
	}

	return fmt.Sprintf("%s%s0%s%s%s%s%s",
		convert(sr.N),
		convert(sr.V),
		convert(sr.B),
		convert(sr.D),
		convert(sr.I),
		convert(sr.Z),
		convert(sr.C),
	)
}

// setZ sets the zero flag of sr according to the contents of reg.
func (sr *Status) setZ(reg byte) {
	if reg == 0x00 {
		sr.Z = true
	} else {
		sr.Z = false
	}
}

// setN sets the negative flag of sr according to the contents of reg.
func (sr *Status) setN(reg byte) {
	if reg&0x80 != 0x00 {
		sr.N = true
	} else {
		sr.N = false
	}
}

// setZN sets both the zero flag and negative flag of sr
// according to the contents of reg.
func (sr *Status) setZN(reg byte) {
	sr.setZ(reg)
	sr.setN(reg)
}

// Registers holds data for each register
// used by the 6502.
type Registers struct {
	// Special purpose registers
	Status *Status // Status register
	PC     uint16  // Program counter
	SP     byte    // Stack pointer

	// General purpose registers
	A byte // Accumulator register
	X byte // Index register X
	Y byte // Index register Y
}

// String implements Stringer
func (r *Registers) String() (repr string) {
	return fmt.Sprintf("%6s | %v\n%6s | %v\n%6s | %v\n%6s | %v\n%6s | %v\n%6s | %v\n",
		"Status",
		r.Status,
		"PC",
		r.PC,
		"SP",
		r.SP,
		"A",
		r.A,
		"X",
		r.X,
		"Y",
		r.Y,
	)
}

// CPU represents to 6502 and its associated registers and memory map.
// This should be declared and used as a singleton during emulator execution.
type CPU struct {
	*memory.Memory
	*Registers

	instructions map[byte]instruction
}

// New initializes a new 6502 CPU with all status bits, register, and memory
// initialized to zero. Memory is the shared memory that the CPU will access.
func New() (c *CPU) {
	cpu := &CPU{
		Registers: &Registers{
			Status: &Status{},
		},
	}
	return cpu
}

// TODO: comment
func (c *CPU) UseMemory(m *memory.Memory) {
	c.Memory = m
}

// TODO: comment
func (c *CPU) Init() {
	c.initInstructions()
	// TODO: start up state
}

// Clear sets every register in c (including PC, SP, and status) to 0x00.
// Retains memory linked through UseMemory.
func (c *CPU) Clear() {
	*c.Registers = Registers{
		Status: &Status{},
	}
}

func (c *CPU) Reset() {
	c.Clear()
	c.Init()
}

// Decode decodes opcode opcode and returns relevant information.
func (c *CPU) Decode(opcode byte) (name string, addressingMode, cycleCost, pageCrossCost, byteCost int, execute func(uint16), err error) {
	if instruction, ok := c.instructions[opcode]; ok {
		return instruction.name,
			instruction.addressingMode,
			instruction.cycleCost,
			instruction.pageCrossCost,
			instruction.byteCost,
			instruction.execute,
			nil
	}
	return "", 0, 0, 0, 0, nil, ErrInvalidOpcode(opcode)
}

// GetAddressWithMode uses addressing mode addressingMode to get
// an address on which any instruction can execute.
// Must be used when c.PC is on an opcode address, otherwise
// the following addresses will be interpreted incorrectly.
func (c *CPU) GetAddressWithMode(addressingMode int) (addr uint16) {
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
		offset := uint16(c.Read(c.PC + 1))
		if offset >= 0x80 {
			// interpret as negative number
			return c.PC + offset - 0x100
		}
		return c.PC + offset

	case modeImmediate:
		// Instructions with modeImmediate take 2 bytes:
		// 1. opcode
		// 2. 8 bit constant value
		return c.PC + 1

	case modeZeroPage:
		// Instructions with modeZeroPage take 2 bytes:
		// 1. opcode
		// 2. zero-page address
		return uint16(c.Read(c.PC + 1))

	case modeZeroPageX:
		// Same as modeZeroPage, but with zero page address being added to X register with wraparound
		return uint16(c.Read(c.PC+1)+c.X) & memory.ZeroPageEnd

	case modeZeroPageY:
		// Same as modeZeroPage, but with zero page address being added to Y register with wraparound
		return uint16(c.Read(c.PC+1)+c.Y) & memory.ZeroPageEnd

	case modeAbsolute:
		// Instructions with modeAbsolute take 3 bytes:
		// 1. opcode
		// 2. least significant byte of address
		// 3. most significant byte of address
		return c.Read16(c.PC + 1)

	case modeAbsoluteX:
		// Same as modeAbsolute, with address being added to contents of X register
		return c.Read16(c.PC+1) + uint16(c.X)

	case modeAbsoluteY:
		// Same as modeAbsolute, with address being added to contents of Y register
		return c.Read16(c.PC+1) + uint16(c.Y)

	case modeIndirect:
		// Instructions with modeIndirect take 3 bytes:
		// 1. opcode
		// 2. least significant byte of address
		// 3. most significant byte of address
		// The formulated address, along with the next,
		// are then accessed again to get the final address.
		return c.Read16(c.Read16(c.PC + 1))

	case modeIndirectX:
		// Instructions with modeIndirectX take 2 bytes:
		// 1. opcode
		// 2. single byte
		// The byte is then added to the X register, which then
		// gives the least significant byte of the target address.
		return c.Read16(uint16(c.Read(c.PC+1) + c.X))

	case modeIndirectY:
		// Instructions with modeIndirectY take 2 bytes:
		// 1. opcode
		// 2. least significant byte of zero page address
		// The zero page address is then accessed, and the data
		// is added to the Y register. The resulting data is the
		// target address.
		return c.Read16(uint16(c.Read(c.PC+1))) + uint16(c.Y)

	default:
		// shouldn't happen, but handle gracefully
		return 0
	}
}

// Step performs a single step of the CPU.
// Briefly, this consists of:
// 1. Retrieving the opcode at current PC.
// 2. Decoding the opcode.
// 3. Incrementing the program counter by the correct amount.
// 4. Performing the instruction.
// TODO: flesh out
func (c *CPU) Step() {
	// 1. Retrieve opcode at current PC
	opcode := c.Read(c.PC)

	// 2. Decode opcode
	name, addressingMode, cycleCost, pageCrossCost, byteCost, execute, err := c.Decode(opcode)
	if IsInvalidOpcodeErr(err) {
		// If the opcode is invalid, continue to the next instruction.
		log.Println(err)
		return
	}
	instructionAddress := c.GetAddressWithMode(addressingMode)

	// 3. Increment program counter
	c.PC += uint16(byteCost)

	// 4. Perform instruction
	execute(instructionAddress)

	// TODO: use these vars later
	_ = name
	_ = cycleCost
	_ = pageCrossCost
}
