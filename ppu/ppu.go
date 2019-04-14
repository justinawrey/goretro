// Package ppu provides functionality related to the nes PPU.
package ppu

import (
	"github.com/justinawrey/nes/display"
)

// common bit masks
const (
	mask0 = 1 << iota
	mask1
	mask2
	mask3
	mask4
	mask5
	mask6
	mask7
	mask01 = 3   // mask bits 0, 1
	mask57 = 224 // mask bits 5, 6, 7
)

// ctrl1 is the first ppu control register.
// See https://wiki.nesdev.com/w/index.php/PPU_registers.
type ctrl1 struct {
	ntAddr    uint16 // Base nametable address
	addrInc   uint16 // Vram address increment per cpu read/write of ppudata
	sprPtable uint16 // Sprite pattern table address ($0000 or $1000)
	bgPtable  uint16 // Background pattern table address ($0000 or $1000)
	sprSize   int    // Sprite size (8x8 pixels or 8x16 pixels)
	nmi       bool   // Whether or not to generate an NMI at start of vertical blanking interval
}

// write writes to the first ppu control register with a byte of data.
func (c *ctrl1) write(data byte) {
	b01 := data & mask01
	b2 := data&mask2 != 0
	b3 := data&mask3 != 0
	b4 := data&mask4 != 0
	b5 := data&mask5 != 0
	b7 := data&mask7 != 0

	// Defaults (data == 0x00)
	c.ntAddr = 0x2000
	c.addrInc = 1
	c.sprPtable = 0x0000
	c.sprSize = 8
	c.nmi = false

	// See https://wiki.nesdev.com/w/index.php/PPU_registers for
	// detailed information on values below.
	switch b01 {
	case 1:
		c.ntAddr = 0x2400
	case 2:
		c.ntAddr = 0x2800
	case 3:
		c.ntAddr = 0x2C00
	default:
	}

	if b2 {
		c.addrInc = 32
	}
	if b3 {
		c.sprPtable = 0x1000
	}
	if b4 {
		c.bgPtable = 0x1000
	}
	if b5 {
		c.sprSize = 16
	}
	if b7 {
		c.nmi = true
	}
}

// ctrl2 is the first ppu control register.
// See https://wiki.nesdev.com/w/index.php/PPU_registers.
type ctrl2 struct {
	monochrome       bool // Whether or not display should be monochrome
	showBgPixels     bool // Show background in leftmost 8 pixels of screen
	showSpritePixels bool // Show sprites in leftmost 8 pixels of screen
	showBg           bool // Show background
	showSprites      bool // Show sprites
	emphasizeRed     bool // Emphasize red on display
	emphasizeGreen   bool // Emphasize green on display
	emphasizeBlue    bool // Emphasize blue on display
}

// write writes to the second ppu control register with a byte of data.
func (c *ctrl2) write(data byte) {
	b0 := data&mask0 != 0
	b1 := data&mask1 != 0
	b2 := data&mask2 != 0
	b3 := data&mask3 != 0
	b4 := data&mask4 != 0
	b5 := data&mask5 != 0
	b6 := data&mask6 != 0
	b7 := data&mask7 != 0

	// Defaults (data == 0x00)
	c.monochrome = false
	c.showBgPixels = false
	c.showSpritePixels = false
	c.showBg = false
	c.showSprites = false
	c.emphasizeRed = false
	c.emphasizeGreen = false
	c.emphasizeBlue = false

	// See https://wiki.nesdev.com/w/index.php/PPU_registers for
	// detailed information on values below.
	if b0 {
		c.monochrome = true
	}
	if b1 {
		c.showBgPixels = true
	}
	if b2 {
		c.showSpritePixels = true
	}
	if b3 {
		c.showBg = true
	}
	if b4 {
		c.showSprites = true
	}
	if b5 {
		c.emphasizeRed = true
	}
	if b6 {
		c.emphasizeGreen = true
	}
	if b7 {
		c.emphasizeBlue = true
	}
}

// status is the ppu status register.
// See https://wiki.nesdev.com/w/index.php/PPU_registers for more info.
type status struct {
	vramWriteIgnore     bool // Whether or not to ignore vram writes
	highScanlineSprites bool // Sprite overflow
	spriteHit           bool // Sprite 0 hit (nonzero pixel of sprite 0 overlaps a nonzero bg pixel)
	vBlank              bool // Whether or not vblank has started
}

// read reads data from the ppu status register and returns data as a byte.
func (s *status) read() (data byte) {
	data = 0x00

	if s.vramWriteIgnore {
		data |= mask4
	}
	if s.highScanlineSprites {
		data |= mask5
	}
	if s.spriteHit {
		data |= mask6
	}
	if s.vBlank {
		data |= mask7
	}
	return data
}

// doubleWriter is a register which has different functionality based on
// whether or not it has been written to an even or odd number of times.
type doubleWriter struct {
	toggle bool
	data1  byte
	data2  byte
}

// write writes to dw with a byte of data.
func (dw *doubleWriter) write(data byte) {
	if !dw.toggle {
		dw.data1 = data
	} else {
		dw.data2 = data
	}
	dw.toggle = !dw.toggle
}

// read16 reads a word of data from dw.
func (dw *doubleWriter) read16() (word uint16) {
	lo := uint16(dw.data2)
	hi := uint16(dw.data1)
	return hi<<8 | lo
}

// Memory sizes
const (
	sprRAMSize = 0x100
	vRAMSize   = 0x1000
)

// IO registers
const (
	ctrlReg1      = 0x2000
	ctrlReg2      = 0x2001
	statusReg     = 0x2002
	sprRAMAddrReg = 0x2003
	sprRAMDataReg = 0x2004
	scrollAddrReg = 0x2005
	vRAMAddrReg   = 0x2006
	vRAMDataReg   = 0x2007
	sprDMAReg     = 0x4014
)

// PPU represents the picture processing unit of the nes.
type PPU struct {
	ctrl1                    // PPU control reg 1
	ctrl2                    // PPU control reg 2
	status                   // PPU status reg
	sprRAMAddr byte          // SPR-RAM read/write address
	scrollAddr *doubleWriter // Fine scroll position (two writes: X scroll, Y scroll)
	vRAMAddr   *doubleWriter // PPU read/write address (two writes: hi byte, lo byte)
	//TODO: DMA

	sprRAM [sprRAMSize]byte // PPU SPR-RAM
	vRAM   [vRAMSize]byte   // PPU VRAM

	*display.Display // Output display driver
}

// New creates a new PPU.
func New() (p *PPU) {
	return &PPU{
		scrollAddr: &doubleWriter{},
		vRAMAddr:   &doubleWriter{},
	}
}

// UseDisplay sets the ppu p to display picture information using the
// display driver d.
func (p *PPU) UseDisplay(d *display.Display) {
	p.Display = d
}

// ReadRegister implements mmio.MemoryMappedIO.
func (p *PPU) ReadRegister(reg uint16) (data byte) {
	switch reg {
	case statusReg:
		return p.status.read()
	case sprRAMAddrReg:
		// TODO: read data from sprRam
		fallthrough
	case vRAMDataReg:
		// TODO: read data from vram
		fallthrough
	default:
		return 0x00
	}
}

// WriteRegister implements mmio.MemoryMappedIO.
func (p *PPU) WriteRegister(reg uint16, data byte) {
	switch reg {
	case ctrlReg1:
		p.ctrl1.write(data)
	case ctrlReg2:
		p.ctrl2.write(data)
	case sprRAMAddrReg:
		p.sprRAMAddr = data
	case sprRAMDataReg:
		p.sprRAM[p.sprRAMAddr] = data
	case scrollAddrReg:
		p.scrollAddr.write(data)
	case vRAMAddrReg:
		p.vRAMAddr.write(data)
	case vRAMDataReg:
		address := p.vRAMAddr.read16()
		p.vRAM[address] = data
	case sprDMAReg:
		//TODO: perform DMA
	default:
	}
}

// Init implements nes.Module.
func (p *PPU) Init() {
	// TODO: PPU start up state.
}

// Clear implements nes.Module.
func (p *PPU) Clear() {
	*p = PPU{
		scrollAddr: &doubleWriter{},
		vRAMAddr:   &doubleWriter{},
	}
}
