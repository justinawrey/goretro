package core

// apu is the audio processing unit of the nes.
type apu struct{}

// NewApu creates a new apu.
func newApu() (a *apu) {
	return &apu{}
}

// ReadRegister implements mmio.MemoryMappedIO.
func (a *apu) readRegister(address uint16) (data byte) { return 0x00 }

// WriteRegister implements mmio.MemoryMappedIO.
func (a *apu) writeRegister(address uint16, data byte) {}

// Init implements nes.Component.
func (a *apu) init() {}

// Clear implements nes.Component.
func (a *apu) clear() {}
