package ppu

const (
	mask0 = 1 << iota
	mask1
	mask2
	mask3
	mask4
	mask5
	mask6
	mask7
	mask01 = 3
	mask57 = 224
)

type ctrl1 struct {
	ntAddr    uint16
	addrInc   uint16
	sprPtable uint16
	bgPtable  uint16
	sprSize   int
	nmi       bool
}

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

type ctrl2 struct {
	monochrome       bool
	showBgPixels     bool
	showSpritePixels bool
	showBg           bool
	showSprites      bool
	color            byte
}

func (c *ctrl2) write(data byte) {
	b0 := data&mask0 != 0
	b1 := data&mask1 != 0
	b2 := data&mask2 != 0
	b3 := data&mask3 != 0
	b4 := data&mask4 != 0
	b57 := (data & mask57) >> 5

	// Defaults (data == 0x00)
	c.monochrome = false
	c.showBgPixels = false
	c.showSpritePixels = false
	c.showBg = false
	c.showSprites = false
	c.color = b57

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
}

type status struct {
	vramWriteIgnore     bool
	highScanlineSprites bool
	spriteHit           bool
	vBlank              bool
}

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

type doubleWriter struct {
	toggle bool
	data1  byte
	data2  byte
}

func (dw *doubleWriter) write(data byte) {
	if !dw.toggle {
		dw.data1 = data
	} else {
		dw.data2 = data
	}
	dw.toggle = !dw.toggle
}

func (dw *doubleWriter) read16() (word uint16) {
	lo := uint16(dw.data2)
	hi := uint16(dw.data1)
	return hi<<8 | lo
}

const (
	// Memory
	sprRAMSize = 0x100
	vRAMSize   = 0x10000

	// IO registers
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

type PPU struct {
	ctrl1
	ctrl2
	status
	sprRAMAddr byte
	scrollAddr *doubleWriter
	vRAMAddr   *doubleWriter
	sprRAM     [sprRAMSize]byte //TODO: flesh out memories
	vRAM       [vRAMSize]byte   //TODO: flesh out memories

	//TODO: DMA
}

func New() (p *PPU) {
	return &PPU{
		scrollAddr: &doubleWriter{},
		vRAMAddr:   &doubleWriter{},
	}
}

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

func (p *PPU) Init() {
	// TODO:
}

func (p *PPU) Clear() {
	*p = PPU{
		scrollAddr: &doubleWriter{},
		vRAMAddr:   &doubleWriter{},
	}
}

func (p *PPU) Reset() {
	p.Clear()
	p.Init()
}
