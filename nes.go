package nes

import (
	"github.com/justinawrey/nes/cpu"
	"github.com/justinawrey/nes/memory"
	"github.com/justinawrey/nes/ppu"
)

type nes struct {
	cpu    *cpu.CPU
	ppu    *ppu.PPU
	memory *memory.Memory
}

func main() {
	// Create all modules
	cpu := cpu.New()
	ppu := ppu.New()
	mem := memory.New()

	// Set up memory links
	cpu.UseMemory(mem)
	mem.AssignMemoryMappedIO(ppu)

	// Get all modules to correct start up state
	cpu.Init()
	ppu.Init()

	nes := &nes{cpu, ppu, mem}

	// // TODO: go further
	_ = nes
}
