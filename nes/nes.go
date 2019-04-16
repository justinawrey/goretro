package nes

import (
	"io"
	"log"

	"github.com/justinawrey/nes/apu"
	"github.com/justinawrey/nes/cartridge"
	"github.com/justinawrey/nes/cpu"
	"github.com/justinawrey/nes/display"
	"github.com/justinawrey/nes/memory"
	"github.com/justinawrey/nes/ppu"
)

// Module describes a module of the nes which can be determanistically initialized and cleared.
type Module interface {
	// Init initializes the module to its correct start up state.
	Init()

	// Clear clears the module to its correct power down state.
	Clear()
}

// initAll initializes all provided modules to the correct start up state.
func initAll(modules ...Module) {
	for _, m := range modules {
		m.Init()
	}
}

// clearAll clears all provided modules.
func clearAll(modules ...Module) {
	for _, m := range modules {
		m.Clear()
	}
}

// resetAll resets all provided modules.
func resetAll(modules ...Module) {
	for _, m := range modules {
		m.Clear()
		m.Init()
	}
}

// NES represents an entire nes system.  It contains a collection of modules
// which work together to power the entire system.
// NES is meant be used primarily as a high-level emulation API.
type NES struct {
	cpu  *cpu.CPU
	ppu  *ppu.PPU
	apu  *apu.APU
	mem  *memory.Memory
	disp *display.Display
	cart *cartridge.Cartridge
}

// New creates a new NES.
func New() (nes *NES) {
	// Create all modules
	cpu := cpu.New()
	ppu := ppu.New()
	apu := apu.New()
	mem := memory.New()
	disp := display.New()

	// Set up memory mapped IO
	cpu.UseMemory(mem)
	mem.AssignMemoryMappedIO(ppu, apu)

	// Use correct display
	ppu.UseDisplay(disp)

	// Get all modules to correct start up state
	initAll(cpu, ppu, apu, disp)

	return &NES{cpu: cpu, ppu: ppu, apu: apu, mem: mem, disp: disp}
}

// Load loads a cartridge with name name into the nes.
func (nes *NES) Load(name string) {
	cart := cartridge.New()
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
	nes.cpu.Step()
}

// Reset resets the nes to its initial power up state.
func (nes *NES) Reset() {
	resetAll(nes.cpu, nes.ppu, nes.apu, nes.disp, nes.mem)
}
