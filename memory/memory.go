package memory

import "github.com/justinawrey/nes/ppu"

const (
	// 6502 has a 64kB memory map
	memSize         = 0x10000
	internalRAMSize = 0x2000

	// See table below for more details
	ramMirrorFreq  = 0x0800
	ppuMirrorFreq  = 0x0008
	ppuMirrorStart = 0x2000
	ZeroPageEnd    = 0x00FF
	ramEnd         = 0x1FFF
	ppuEnd         = 0x3FFF
	apuIoEnd       = 0x4017
	testModeEnd    = 0x401F
	cartridgeEnd   = 0xFFFF
)

// memoryMap specifies a module which can be read from / written to via memory mapped io.
// Rather than directly exposing the entire ppu / apu / joystick objects to memory,
// we instead only expose their read / write methods.
// See Memory for main usage.
type memoryMap interface {
	WriteRegister(uint16, byte)
	ReadRegister(uint16) byte
}

// Memory is the 64kB memory map contained within the CPU.
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
type Memory struct {
	internal     [internalRAMSize]byte
	ppuMemoryMap memoryMap
}

// New constructs a new Memory.
func New() (m *Memory) {
	return &Memory{}
}

// AssignMemoryMappedIO sets up writing to / reading from memory to be memory mapped
// with the specified argument modules.
func (m *Memory) AssignMemoryMappedIO(ppu *ppu.PPU) {
	m.ppuMemoryMap = ppu
}

// Read reads a byte of data from the memory map at address.
func (m *Memory) Read(address uint16) (data byte) {
	switch {
	case address <= ramEnd:
		// Internal CPU RAM.  Mirrored memory at a frequency of 0x0800.
		// We can make a small shortcut by only writing
		// to a single 'chunk' of mirrored memory using a modulus.
		return m.internal[address%ramMirrorFreq]
	case address <= ppuEnd:
		// Memory mapped IO for PPU.  Mirrored memory at a frequency of 0x0008.
		// Same modulus trick as above.
		register := (address % ppuMirrorFreq) + ppuMirrorStart
		return m.ppuMemoryMap.ReadRegister(register)
	default:
		// TODO: handle the rest
		return 0x00
	}
}

// Write writes a byte of data to the memory map at address.
func (m *Memory) Write(address uint16, data byte) {
	switch {
	case address <= ramEnd:
		m.internal[address%ramMirrorFreq] = data
	case address <= ppuEnd:
		register := (address % ppuMirrorFreq) + ppuMirrorStart
		m.ppuMemoryMap.WriteRegister(register, data)
	default:
		// TODO: handle the rest
	}
}

// Clear sets all data in the memory map to 0x00.
func (m *Memory) Clear() {
	for i := range m.internal {
		m.internal[i] = 0x00
	}
}

// Read16 reads two bytes, in little endian order, starting
// at memory location from.  The bytes are concatenated
// into a two byte word and returned.
func (m *Memory) Read16(from uint16) (word uint16) {
	lo := uint16(m.Read(from))
	hi := uint16(m.Read(from + 1))
	return hi<<8 | lo
}
