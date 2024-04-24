package core

// memoryMappedIO specifies a module which can be read from / written to via memory mapped io.
// Each MemoryMappedIO module exposes ReadRegister and WriteRegister functions which perform
// differently for each module.  For example, the nes PPU may enable writing to certain control registers
// and input devices may enable reading from certain input feedback registers.
type memoryMappedIO interface {
	// readRegister reads a byte of data from a certain special purpose 16-bit register
	readRegister(uint16) byte

	// writeRegister writes a byte of data to a certain special purpose 16-bit register
	writeRegister(uint16, byte)
}
