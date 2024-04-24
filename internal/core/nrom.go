package core

// NROM - iNES mapper #0.
// See https://wiki.nesdev.com/w/index.php/NROM.
type nrom struct {
	*cartridge
}

// readRegister implements mmio.MemoryMappedIO.
func (nr *nrom) readRegister(address uint16) (data byte) {
	return nr.readPrg(address)
}

// writeRegister implements mmio.MemoryMappedIO.
func (nr *nrom) writeRegister(address uint16, data byte) {
	nr.writePrg(address, data)
}

// load implements Mapper.
func (nr *nrom) load(bytes []byte) {
	nr.loadPrg(bytes)
	nr.loadChr(bytes)
}
