package core

import (
	"io"
	"log"

	"github.com/justinawrey/goretro/internal/display"
	"github.com/justinawrey/goretro/internal/input"
)

// initAll initializes all provided modules to the correct start up state.
func initAll(modules ...Component) {
	for _, m := range modules {
		m.init()
	}
}

// clearAll clears all provided modules.
func clearAll(modules ...Component) {
	for _, m := range modules {
		m.clear()
	}
}

// resetAll resets all provided modules.
func resetAll(modules ...Component) {
	for _, m := range modules {
		m.clear()
		m.init()
	}
}

// NES represents an entire nes system.  It contains a collection of modules
// which work together to power the entire system.
// NES is meant be used primarily as a high-level emulation API.
type NES struct {
	cpu  *CPU
	ppu  *PPU
	apu  *APU
	mem  *Memory
	cart *Cartridge

	disp  *display.Display
	input *input.Input
}

// New creates a new NES.
func NewNes() (nes *NES) {
	// Create all modules
	cpu := newCpu()
	ppu := newPpu()
	apu := newApu()
	mem := newMemory()

	disp := display.NewDisplay()
	input := input.NewInput()

	// Set up memory mapped IO
	cpu.useMemory(mem)
	mem.assignMemoryMappedIO(ppu, apu)

	// Use correct display
	ppu.useDisplay(disp)

	// Get all modules to correct start up state
	initAll(cpu, ppu, apu, disp, input)

	return &NES{
		cpu:   cpu,
		ppu:   ppu,
		apu:   apu,
		mem:   mem,
		disp:  disp,
		input: input,
	}
}

// Load loads a cartridge with name name into the nes.
func (nes *NES) Load(name string) {
	cart := NewCartridge()
	cart.Load(name)
	nes.mem.AssignMemoryMappedIO(cart)
	nes.cart = cart

	log.Println("loaded: ", name)
	log.Println("mapper: ", nes.cart.MapperNum)
}

// OutputTo sets the nes to log its execution to io.Writer w.
func (nes *NES) OutputTo(w io.Writer) {
	nes.cpu.OutputTo(w)
}

// Start begins executing the loaded cartridge.
// For now, this is for nestest.
// TODO: make this actually start
func (nes *NES) Start() {
	nes.cpu.PC = 0xC000
	nes.cpu.Step()
	nes.cpu.Step()
	nes.cpu.Step()
}

// Reset resets the nes to its initial power up state.
func (nes *NES) Reset() {
	resetAll(nes.cpu, nes.ppu, nes.apu, nes.disp, nes.mem)
}
