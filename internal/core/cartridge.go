package core

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

// iNES mapper numbers for supported mappers
const (
	mapperNROM = iota
)

// iNES related memory sizes
const (
	headerLen     = 0x10
	trainerLen    = 0x200
	prgROMBankLen = 0x4000
	chrROMBankLen = 0x2000
)

// errInvalidHeader is an error related to the header of an iNES file
// being invalid.  This means it is either the incorrect size (16 bytes) or has
// malformed contents. See http://wiki.nesdev.com/w/index.php/INES.
type errInvalidHeader string

// newErrInvalidHeader creates a new invalidHeaderErr with a message.
func newErrInvalidHeader(message string) (err errInvalidHeader) {
	return errInvalidHeader(message)
}

// Error implements error.
func (err errInvalidHeader) Error() (message string) {
	return fmt.Sprintf("invalid header: %s", string(err))
}

// Mapper is a memory mapper as used by NES cartridges.
// For a list of mappers and their iNES numbers see http://wiki.nesdev.com/w/index.php/Mapper.
type mapper interface {
	// Mapper implements mmio.MemoryMappedIO.
	memoryMappedIO

	// Load specifies how memory should be loaded on initial mapper load.
	// A slice of bytes is provided and represents binary data as loaded directly
	// from a .nes file.
	load([]byte)
}

// newMapper creates a new Mapper from the provided iNES mapper id.
// The cartridge to which the Mapper is associated to must also be provided.
func newMapper(id int, c *cartridge) (m mapper) {
	switch id {
	default:
		fallthrough
	case 0:
		return &nrom{c}
	}
}

// cartridge represents a nes Cartridge.
// It implements mmio.MemoryMappedIO and thus can be directly
// written to / read from via registers.
type cartridge struct {
	mapper

	prgROM []byte
	chrROM []byte

	mapperNum   int // iNES mapper #
	prgROMBanks int // Number of prgROM banks
	chrROMBanks int // Number of chrROM (VROM) banks

	romStart int // byte number at which actual rom starts
	chrStart int // byte number at which actual chrROM starts

	hasSRAM             bool // whether or not this cart supports SRAM
	hasTrainer          bool // whether or not there is a 512kB trainer preceding rom
	vertMirroring       bool // whether to use vertical mirroring (or horizontal mirroring)
	fourScreenMirroring bool // whether or not to ignore above flag and use four screen mirroring
}

// New creates a new cartridge.
func newcartridge() (c *cartridge) {
	return &cartridge{}
}

// Load loads a .nes file with name name into the cartridge c.
func (c *cartridge) Load(name string) {
	// Only use iNES format
	if !strings.HasSuffix(name, ".nes") {
		log.Fatalln("only .nes format is supported")
	}

	// Read .nes file
	bytes, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatalln("error loading cartridge: " + err.Error())
	}

	// Decode useful information from ROM header
	if err := c.decodeHeader(bytes[:headerLen]); err != nil {
		log.Fatalln("error loading cartridge: " + err.Error())
	}

	// Use correct mapper as according to file header
	c.mapper = newMapper(c.mapperNum, c)

	// Load data into prgRom and chrRom
	c.mapper.load(bytes)
}

// writePrg writes to prgROM at address, with data.
// ROM is not meant to be written to (as the name implies),
// but this is here for completeness.
func (c *cartridge) writePrg(address uint16, data byte) {
	c.prgROM[address-prgROMStart] = data
}

// readPrg reads from prgROM at address and returns a byte of data.
func (c *cartridge) readPrg(address uint16) (data byte) {
	return c.prgROM[address-prgROMStart]
}

// loadPrg loads data into prgROM.  It performs the default action of loading either:
//
//	a) if there is only one bank of prgROM, it gets mirrored at 0x8000 and 0xC000.
//	b) if there is two, it gets fully loaded from 0x8000 to 0xFFFF.
//
// More complicated mappers may require their own loading
// logic which can be defined in their respective Load methods.
func (c *cartridge) loadPrg(bytes []byte) {
	c.prgROM = make([]byte, prgROMBankLen*2)

	if c.prgROMBanks == 1 {
		// If there is only one bank, it gets mirrored at 0x8000 and 0xC000.
		for i := c.romStart; i < c.romStart+prgROMBankLen; i++ {
			c.prgROM[i-c.romStart] = bytes[i]
			c.prgROM[(i-c.romStart)+prgROMBankLen] = bytes[i]
		}
	} else if c.prgROMBanks == 2 {
		// If there are two banks, fill 0x8000 - 0xFFFF.  No mirroring required.
		for i := c.romStart; i < c.romStart+prgROMBankLen*2; i++ {
			c.prgROM[i-c.romStart] = bytes[i]
		}
	}

}

// loadChr loads data into chrROM.
func (c *cartridge) loadChr(bytes []byte) {
	c.chrROM = make([]byte, chrROMBankLen)

	// Fill chrROM.
	for i := c.chrStart; i < c.chrStart+chrROMBankLen; i++ {
		c.chrROM[i-c.chrStart] = bytes[i]
	}
}

// decodeHeader decodes the raw header bytes header and populates
// c with pertinent information.  Returns an errInvalidHeader if
// the supplied header bytes are invalid.
func (c *cartridge) decodeHeader(header []byte) (err error) {
	// Header should be 16 bytes long.
	if len(header) != headerLen {
		return newErrInvalidHeader("invalid header length")
	}

	// First 4 bytes: 'NES' followed by MS-DOS EOF.
	if !bytes.Equal(header[:4], []byte{0x4E, 0x45, 0x53, 0x1A}) {
		return newErrInvalidHeader("invalid constant")
	}

	c.prgROMBanks = int(header[4])
	c.chrROMBanks = int(header[5])

	// 6th byte specifies flags, see https://wiki.nesdev.com/w/index.php/INES.
	flags6 := header[6]
	c.vertMirroring = flags6&mask0 != 0
	c.hasSRAM = flags6&mask1 != 0
	c.hasTrainer = flags6&mask2 != 0
	c.fourScreenMirroring = flags6&mask3 != 0

	// 7th byte also specifies flags, see link.
	flags7 := header[7]
	c.mapperNum = int(byte(flags6&0xF0>>4) | byte(flags7&0xF0))

	// Compute prg and chr start bytes.
	c.romStart = headerLen
	c.chrStart = headerLen + (chrROMBankLen * c.chrROMBanks)
	if c.hasTrainer {
		c.romStart += trainerLen
		c.chrStart += trainerLen
	}

	return nil
}
