package ppu

import "github.com/justinawrey/nes/memory"

type PPU struct {
	cpuMem *memory.Memory
}

func NewPPU(m *memory.Memory) (p *PPU) {
	ppu := &PPU{
		cpuMem: m,
	}
	return ppu
}
