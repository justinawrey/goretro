package nes

import (
	"github.com/justinawrey/nes/apu"
	"github.com/justinawrey/nes/cartridge"
	"github.com/justinawrey/nes/cpu"
	"github.com/justinawrey/nes/display"
	"github.com/justinawrey/nes/memory"
	"github.com/justinawrey/nes/ppu"
)

type module interface {
	Init()
	Clear()
}

func initAll(modules ...module) {
	for _, m := range modules {
		m.Init()
	}
}

func clearAll(modules ...module) {
	for _, m := range modules {
		m.Clear()
	}
}

func resetAll(modules ...module) {
	for _, m := range modules {
		m.Clear()
		m.Init()
	}
}

type NES struct {
	cpu *cpu.CPU
	ppu *ppu.PPU
	apu *apu.APU
	mem *memory.Memory
	dis *display.Display
}

func New() (nes *NES) {
	// Create all modules
	cpu := cpu.New()
	ppu := ppu.New()
	apu := apu.New()
	mem := memory.New()
	dis := display.New()

	// Set up memory mapped IO
	cpu.UseMemory(mem)
	mem.AssignMemoryMappedIO(ppu, apu)

	// Use correct display
	ppu.UseDisplay(dis)

	// Get all modules to correct start up state
	initAll(cpu, ppu, apu, dis)

	return &NES{cpu, ppu, apu, mem, dis}
}

func (nes *NES) Load(name string) {
	cart := cartridge.New()
	cart.Load(name)
	nes.mem.AssignMemoryMappedIO(cart)
}

func (nes *NES) Start() {
	// TODO:
}

func (nes *NES) Reset() {
	resetAll(nes.cpu, nes.ppu, nes.apu, nes.dis, nes.mem)
}
