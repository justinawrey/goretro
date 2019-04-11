package cartridge

func NewMapper(c *Cartridge) (m Mapper) {
	switch c.mapperNum {
	default:
		fallthrough
	case 0:
		return &NROM{c}
	}
}

type NROM struct {
	*Cartridge
}

func (nr *NROM) ReadRegister(address uint16) (data byte) {
	return nr.readPrg(address)
}

func (nr *NROM) WriteRegister(address uint16, data byte) {
	nr.writePrg(address, data)
}

func (nr *NROM) Load(bytes []byte) {
	nr.loadPrg(bytes)
	nr.loadChr(bytes)
}
