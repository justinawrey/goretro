package nes

import (
	"github.com/justinawrey/nes/cpu"
	"github.com/justinawrey/nes/memory"
	"github.com/justinawrey/nes/ppu"
)

type nes struct {
	cpu       *cpu.CPU
	ppu       *ppu.PPU
	sharedMem *memory.Memory
}

func main() {
	sharedMem := new(memory.Memory)
	cpu := cpu.NewCPU(sharedMem)
	ppu := ppu.NewPPU(sharedMem)

	nes := &nes{cpu, ppu, sharedMem}

	// TODO: go further
	_ = nes
}
