package nes

import (
	"bytes"
	"fmt"
	"os"
)

// iNES related memory sizes
const (
	iNesHeaderLen  = 0x10
	iNesTrainerLen = 0x200
	prgROMBankLen  = 0x4000
	chrROMBankLen  = 0x2000
)

// errINesFileInvalid is an error related to a given iNES file being invalid.
// This means that the file has an invalid header, e.g. the file is < 16 bytes long or
// is >= 16 bytes long but the first 16 bytes have malformed contents.
// See http://wiki.nesdev.com/w/index.php/INES.
type errINesFileInvalid string

func newErrINesFileInvalid(message string) errINesFileInvalid {
	return errINesFileInvalid(message)
}

// Error implements error.
func (err errINesFileInvalid) Error() string {
	return fmt.Sprintf("iNES file invalid: %v", string(err))
}

// cartridge represents a nes cartridge.
type cartridge struct {
	path string // the path at which the backing iNES file resides on disk

	mapperNum   int // iNES mapper #
	prgROMBanks int // Number of prgROM banks
	chrROMBanks int // Number of chrROM (VROM) banks
	ramBanks    int // Number of RAM banks

	hasSRAM             bool // whether or not this cart supports SRAM
	hasTrainer          bool // whether or not there is a 512kB trainer preceding rom
	vertMirroring       bool // whether to use vertical mirroring (or horizontal mirroring)
	fourScreenMirroring bool // whether or not to ignore above flag and use four screen mirroring
}

// String implements Stringer.
func (c *cartridge) String() string {
	return fmt.Sprintf("[%v] mapper: %v, prg ROM banks: %v, chr ROM banks: %v, RAM banks: %v", c.path, c.mapperNum, c.prgROMBanks, c.chrROMBanks, c.ramBanks)
}

// newCartridge creates a new catridge from the file specified at relative path path.
// Only supports the iNES file type.  If the file type is detected to be iNES, parse out
// and store all relevant information.  If the file is found but does not satisfy the iNES format,
// returns an error of type errINesFileInvalid.
func newCartridge(path string) (*cartridge, error) {
	decode := func(file []byte, c *cartridge) error {
		if len(file) < iNesHeaderLen {
			return newErrINesFileInvalid("file less than 16 bytes long")
		}

		if !bytes.Equal(file[:4], []byte{0x4E, 0x45, 0x53, 0x1A}) {
			return newErrINesFileInvalid("invalid first 4 bytes, should be 'NES' followed by MS-DOS EOF")
		}

		c.prgROMBanks = int(file[4])
		c.chrROMBanks = int(file[5])

		control := file[6]
		control2 := file[7]

		c.vertMirroring = control&mask0 != 0
		c.hasSRAM = control&mask1 != 0
		c.hasTrainer = control&mask2 != 0
		c.fourScreenMirroring = control&mask3 != 0
		c.mapperNum = int(byte(control&0b11110000>>4) | byte(control2&0b11110000))

		// for compatibility with old iNES versions, when 0 is indicated
		// in byte 8 we should consider it as 1 ram bank
		c.ramBanks = int(file[8])
		if c.ramBanks == 0 {
			c.ramBanks = 1
		}

		return nil
	}

	c := &cartridge{path: path}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Decode useful information from ROM header
	if err := decode(bytes[:iNesHeaderLen], c); err != nil {
		return nil, err
	}

	return c, nil
}
