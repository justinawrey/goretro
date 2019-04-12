// Package apu provides functionality related to the nes audio processing unit.
package apu

// APU is the audio processing unit of the nes.
type APU struct {
}

// New creates a new APU.
func New() (a *APU) {
	return &APU{}
}

// ReadRegister implements mmio.MemoryMappedIO.
func (a *APU) ReadRegister(address uint16) (data byte) { return 0x00 }

// WriteRegister implements mmio.MemoryMappedIO.
func (a *APU) WriteRegister(address uint16, data byte) {}

// Init implements nes.Module.
func (a *APU) Init() {}

// Clear implements nes.Module.
func (a *APU) Clear() {}
