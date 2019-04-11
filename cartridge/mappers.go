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
	address = address - prgROMStart
	return nr.prgROM[address]
}

func (nr *NROM) WriteRegister(address uint16, data byte) {
	address = address - prgROMStart
	nr.prgROM[address] = data
}

func (nr *NROM) Load(bytes []byte) {
	// See https://wiki.nesdev.com/w/index.php/NROM
	nr.prgROM = make([]byte, prgROMBankLen*2)
	nr.chrROM = make([]byte, chrROMBankLen)

	ROMStart := headerLen
	CHRStart := headerLen + (chrROMBankLen * nr.chrROMBanks)
	if nr.hasTrainer {
		ROMStart += trainerLen
		CHRStart += trainerLen
	}

	// Fill prgROM
	if nr.prgROMBanks == 1 {
		// If there is only one bank, it gets mirrored at 0x8000 and 0xC000
		for i := ROMStart; i < ROMStart+prgROMBankLen; i++ {
			nr.prgROM[i-ROMStart] = bytes[i]
			nr.prgROM[(i-ROMStart)+prgROMBankLen] = bytes[i]
		}
	} else if nr.prgROMBanks == 2 {
		// If there are two banks, fill 0x8000 - 0xFFFF
		for i := ROMStart; i < ROMStart+prgROMBankLen*2; i++ {
			nr.prgROM[i-ROMStart] = bytes[i]
		}
	}

	// Fill chrROM
	for i := CHRStart; i < CHRStart+chrROMBankLen; i++ {
		nr.chrROM[i-CHRStart] = bytes[i]
	}
}
