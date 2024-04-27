package core

import (
	"fmt"
	"io"

	"github.com/justinawrey/goretro/internal/audio"
	"github.com/justinawrey/goretro/internal/display"
	"github.com/justinawrey/goretro/internal/input"
)

// NES represents an entire nes system.  It contains a collection of modules
// which work together to power the entire system.
// NES is meant be used primarily as a high-level emulation API.
type nes struct {
	cpu  *cpu
	ppu  *ppu
	apu  *apu
	mem  *memory
	cart *cartridge

	disp  *display.Display
	input *input.Input
	audio *audio.Audio
}

func (n *nes) UseDisplay(d *display.Display) {
	n.disp = d
}

func (n *nes) UseInput(i *input.Input) {
	n.input = i
}

func (n *nes) UseAudio(a *audio.Audio) {
	n.audio = a
}

func (n *nes) UseCartridge(path string) error {
	cart, err := newCartridge(path)
	if err != nil {
		return err
	}

	n.cart = cart

	fmt.Println(cart)

	return nil
}

// New creates a new NES.
func NewNes() *nes {
	// Create all modules
	// cpu := newCpu()
	// ppu := newPpu()
	// apu := newApu()
	// mem := newMemory()
	//
	// // TODO: bring this back
	// // Set up memory mapped IO
	// // cpu.useMemory(mem)
	// // mem.assignMemoryMappedIO(ppu, apu)
	//
	// return &nes{
	// 	cpu: cpu,
	// 	ppu: ppu,
	// 	apu: apu,
	// 	mem: mem,
	// }
	return &nes{}
}

// OutputTo sets the nes to log its execution to io.Writer w.
func (n *nes) OutputTo(w io.Writer) {
	// nes.cpu.OutputTo(w)
}

// Start begins executing the loaded cartridge.
// For now, this is for nestest.
// TODO: make this actually start
func (n *nes) Start() {
	// nes.cpu.PC = 0xC000
	// nes.cpu.Step()
	// nes.cpu.Step()
	// nes.cpu.Step()
}

// Reset resets the nes to its initial power up state.
// func (n *nes) Reset() {
// 	resetAll(nes.cpu, nes.ppu, nes.apu, nes.disp, nes.mem)
// }
