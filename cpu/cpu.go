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
func (sr *Status) String() string {
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
// used by the 6502,
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
func (r *Registers) String() string {
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
func (m *MemoryMap) String() string {
	return ""
}

func (m *MemoryMap) Write(to uint16, data byte) {
	m[to] = data
}

func (m *MemoryMap) Read(from uint16) byte {
	return m[from]
}

// CPU represents to 6502 and its associated registers and memory map.
// This should be declared and used as a singleton during emulator execution.
type CPU struct {
	*MemoryMap
	*Registers
}
