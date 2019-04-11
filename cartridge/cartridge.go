package cartridge

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/justinawrey/nes/mmio"
)

const (
	mapperNROM = iota
)

const (
	mask0 = 1 << iota
	mask1
	mask2
	mask3
)

const (
	headerLen     = 0x10
	trainerLen    = 0x200
	prgROMBankLen = 0x4000
	chrROMBankLen = 0x2000

	prgROMStart = 0x8000
)

type invalidHeaderErr string

func newInvalidHeaderErr(message string) (err invalidHeaderErr) {
	return invalidHeaderErr(message)
}

func (err invalidHeaderErr) Error() (message string) {
	return fmt.Sprintf("invalid header: %s", string(err))
}

type Mapper interface {
	mmio.MemoryMappedIO
	Load([]byte)
}

type Cartridge struct {
	Mapper

	prgROM []byte
	chrROM []byte

	mapperNum   int
	prgROMBanks int
	chrROMBanks int

	hasSRAM             bool
	hasTrainer          bool
	vertMirroring       bool
	fourScreenMirroring bool
}

func New() (c *Cartridge) {
	return &Cartridge{}
}

func (c *Cartridge) decodeHeader(header []byte) (err error) {
	// Header should be 16 bytes long
	if len(header) != headerLen {
		return newInvalidHeaderErr("invalid header length")
	}

	// First 4 bytes: 'NES' followed by MS-DOS EOF
	if !bytes.Equal(header[:4], []byte{0x4E, 0x45, 0x53, 0x1A}) {
		return newInvalidHeaderErr("invalid constant")
	}

	c.prgROMBanks = int(header[4])
	c.chrROMBanks = int(header[5])

	// 6th byte specifies flags, see https://wiki.nesdev.com/w/index.php/INES
	flags6 := header[6]
	c.vertMirroring = flags6&mask0 != 0
	c.hasSRAM = flags6&mask1 != 0
	c.hasTrainer = flags6&mask2 != 0
	c.fourScreenMirroring = flags6&mask3 != 0

	// 7th byte also specifies flags, see link
	flags7 := header[7]
	c.mapperNum = int(byte(flags6&0xF0>>4) | byte(flags7&0xF0))

	return nil
}

func (c *Cartridge) Load(name string) {
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
	c.Mapper = NewMapper(c)

	// Load data into prgRom and chrRom
	c.Mapper.Load(bytes)
}
