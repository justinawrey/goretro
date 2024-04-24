package nes

// NROM - iNES mapper #0.
// See https://wiki.nesdev.com/w/index.php/NROM.
type NROM struct {
	*Cartridge
}

// ReadRegister implements mmio.MemoryMappedIO.
func (nr *NROM) ReadRegister(address uint16) (data byte) {
	return nr.readPrg(address)
}

// WriteRegister implements mmio.MemoryMappedIO.
func (nr *NROM) WriteRegister(address uint16, data byte) {
	nr.writePrg(address, data)
}

// Load implements Mapper.
func (nr *NROM) Load(bytes []byte) {
	nr.loadPrg(bytes)
	nr.loadChr(bytes)
}
