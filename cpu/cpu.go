package cpu

import (
	"fmt"
	"log"
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

// Clear clears all bits of the status register,
// i.e. will set all flags to '0'.
func (sr *Status) Clear() {
	sr.C = false
	sr.Z = false
	sr.I = false
	sr.D = false
	sr.B = false
	sr.V = false
	sr.N = false
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

const (
	memSize = 0xFFFF // 6502 has a 64kB memory map

	// See table below for more details
	zeroPageEnd  = 0x00FF
	ramEnd       = 0x1FFF
	ppuEnd       = 0x3FFF
	apuIoEnd     = 0x4017
	testModeEnd  = 0x401F
	cartridgeEnd = 0xFFFF
)

// MemoryMap is the 64kB memory map contained within the CPU.
// The memory is organized as follows (https://wiki.nesdev.com/w/index.php/CPU_memory_map):
//
// AddressRange	Size	Device
// ---------------------------------------------
// $0000-$07FF	$0800	2KB internal RAM
// $0800-$0FFF	$0800	Mirrors of $0000-$07FF
// $1000-$17FF	$0800
// $1800-$1FFF	$0800
// $2000-$2007	$0008	NES PPU registers
// $2008-$3FFF	$1FF8	Mirrors of $2000-2007 (repeats every 8 bytes)
// $4000-$4017	$0018	NES APU and I/O registers
// $4018-$401F	$0008	APU and I/O functionality that is normally disabled. See CPU Test Mode.
// $4020-$FFFF	$BFE0	Cartridge space: PRG ROM, PRG RAM, and mapper registers (See Note)
type MemoryMap [memSize]byte

// String implements Stringer
// TODO: complete memory dump
func (m *MemoryMap) String() (repr string) {
	return ""
}

// Write writes to the memory map.
// TODO: make robust
func (m *MemoryMap) Write(to uint16, data byte) {
	m[to] = data
}

// Read reads from the memory map.
// TODO: make robust
func (m *MemoryMap) Read(from uint16) (b byte) {
	return m[from]
}

// Read16 reads two bytes, in little endian order, starting
// at memory location from.  The bytes are concatenated
// into a two byte word and returned.
func (m *MemoryMap) Read16(from uint16) (word uint16) {
	lo := uint16(m.Read(from))
	hi := uint16(m.Read(from + 1))
	return hi<<8 | lo
}

// CPU represents to 6502 and its associated registers and memory map.
// This should be declared and used as a singleton during emulator execution.
type CPU struct {
	*MemoryMap
	*Registers

	instructions map[byte]instruction
}

// NewCPU initializes a new 6502 CPU with all status bits, register, and memory
// initialized to zero.
func NewCPU() (c *CPU) {
	cpu := &CPU{
		MemoryMap: &MemoryMap{},
		Registers: &Registers{
			Status: &Status{},
		},
	}
	cpu.initInstructions()
	return cpu
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
		// The address containing the constant value
		// will only be accessed if the branch succeeeds.
		// TODO: special case?
		return c.PC + uint16(c.Read(c.PC+1))

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
		return uint16(c.Read(c.PC+1)+c.X) & zeroPageEnd

	case modeZeroPageY:
		// Same as modeZeroPage, but with zero page address being added to Y register with wraparound
		return uint16(c.Read(c.PC+1)+c.Y) & zeroPageEnd

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
// 1. Retrieving the opcode at current PC
// 2. Decoding the opcode.
// 3. Performing the instruction.
// 4. Incrementing the program counter by the correct amount.
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

	// 3. Perform instruction
	execute(instructionAddress)

	// 4. Increment program counter
	c.PC += uint16(byteCost)

	// TODO: use these vars later
	_ = name
	_ = cycleCost
	_ = pageCrossCost
}
