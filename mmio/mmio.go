package mmio

// MemoryMappedIO specifies a module which can be read from / written to via memory mapped io.
// Rather than directly exposing the entire ppu / apu / joystick objects to memory,
// we instead only expose their read / write methods.
// See Memory for main usage.
type MemoryMappedIO interface {
	ReadRegister(uint16) byte
	WriteRegister(uint16, byte)
}
