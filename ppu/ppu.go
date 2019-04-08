package ppu

type PPU struct {
}

func New() (p *PPU) {
	return &PPU{}
}

// $2000, $2001 - control from cpu -> ppu
// $2002 - ppu status, cpu can read it
// $2003, $2004 - data transfer from cpu -> ppu sprram
// $2005, $2006, $2007 - data transfer from cpu -> ppu vram
// $4014 - data transfer from cpu -> ppu sprram (DMA)

func (p *PPU) ReadRegister(reg uint16) (data byte) {
	return 0x00
}

func (p *PPU) WriteRegister(reg uint16, data byte) {

}

func (p *PPU) Initialize() {

}

func (p *PPU) Clear() {

}
