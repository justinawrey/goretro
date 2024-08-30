package nes

// apu is the audio processing unit of the nes.
type apu struct{}

// NewApu creates a new apu.
func newApu() (a *apu) {
	return &apu{}
}
