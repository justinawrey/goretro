package nes

import (
	"github.com/justinawrey/nes/apu"
	"github.com/justinawrey/nes/cartridge"
	"github.com/justinawrey/nes/cpu"
	"github.com/justinawrey/nes/display"
	"github.com/justinawrey/nes/memory"
	"github.com/justinawrey/nes/ppu"
)

type nes struct {
	cpu *cpu.CPU
	ppu *ppu.PPU
	apu *apu.APU
	mem *memory.Memory
	dis *display.Display
}

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

func main() {
	// Create all modules
	cpu := cpu.New()
	ppu := ppu.New()
	apu := apu.New()
	mem := memory.New()
	dis := display.New()

	// Load a .nes file
	cart := cartridge.New()
	cart.Load("donkeykong.nes")

	// Set up memory mapped IO
	cpu.UseMemory(mem)
	mem.AssignMemoryMappedIO(ppu, cart, apu)

	// Use correct display
	ppu.UseDisplay(dis)

	// Get all modules to correct start up state
	initAll(cpu, ppu, apu, dis)

	// TODO: go further
	nes := &nes{cpu, ppu, apu, mem, dis}
	_ = nes
}
