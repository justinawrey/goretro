package cpu

import "fmt"

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
// TODO: make robust
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
// TODO: handle case where opcode is invalid
func (c *CPU) Decode(opcode byte) (name string, addressingMode, cycleCost, pageCrossCost, byteCost int, execute func()) {
	instruction := c.instructions[opcode]
	return instruction.name,
		instruction.addressingMode,
		instruction.cycleCost,
		instruction.pageCrossCost,
		instruction.byteCost,
		instruction.execute
}

// GetAddressWithMode ...
// TODO: complete
// TODO: maybe assert c.PC is on an opcode address
func (c *CPU) GetAddressWithMode(addressingMode int) (addr uint16) {
	switch addressingMode {
	case modeImplied:
		// TODO: flag as unneeded somehow
		return 0
	case modeRelative:
		// TODO: complete
		return 0
	case modeAccumulator:
		// TODO: flag as unneeded somehow
		return 0
	case modeImmediate:
		// TODO: complete
		return 0
	case modeZeroPage:
		// TODO: complete
		return 0
	case modeZeroPageX:
		// TODO: complete
		return 0
	case modeZeroPageY:
		// TODO: complete
		return 0
	case modeIndirect:
		// Same as modeAbsolute
		fallthrough
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
	case modeIndirectX:
		// Instructions with modeIndirectX take 2 bytes:
		// 1. opcode
		// 2. zero page byte address
		// byte retrieved in 2. is then added to the contents of X register with zero page wraparound
		// TODO: might be able to be simplified
		return (uint16(c.Read(c.PC+1)) + uint16(c.X)) & zeroPageEnd
	case modeIndirectY:
		// Same as modeIndirectX, with byte retrieved being added to contents of Y register with zero page wraparound
		// TODO: might be able to be simplified
		return (uint16(c.Read(c.PC+1)) + uint16(c.Y)) & zeroPageEnd
	default:
		// TODO: shouldn't happen, but should handle gracefully
		return 0
	}
}
