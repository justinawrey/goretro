package nes

import (
	"log"

	"github.com/justinawrey/nes/apu"
	"github.com/justinawrey/nes/cartridge"
	"github.com/justinawrey/nes/cpu"
	"github.com/justinawrey/nes/display"
	"github.com/justinawrey/nes/memory"
	"github.com/justinawrey/nes/ppu"
)

type Module interface {
	Init()
	Clear()
}

func initAll(modules ...Module) {
	for _, m := range modules {
		m.Init()
	}
}

func clearAll(modules ...Module) {
	for _, m := range modules {
		m.Clear()
	}
}

func resetAll(modules ...Module) {
	for _, m := range modules {
		m.Clear()
		m.Init()
	}
}

type NES struct {
	cpu  *cpu.CPU
	ppu  *ppu.PPU
	apu  *apu.APU
	mem  *memory.Memory
	disp *display.Display
	cart *cartridge.Cartridge
}

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

func (nes *NES) Load(name string) {
	cart := cartridge.New()
	cart.Load(name)
	nes.mem.AssignMemoryMappedIO(cart)
	nes.cart = cart

	log.Println("loaded: ", name)
	log.Println("mapper: ", nes.cart.MapperNum)
}

func (nes *NES) Start() {
	// TODO:
}

func (nes *NES) Reset() {
	resetAll(nes.cpu, nes.ppu, nes.apu, nes.disp, nes.mem)
}
