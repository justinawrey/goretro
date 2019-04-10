package cartridge

import "github.com/justinawrey/nes/mappers"

type Cartridge struct {
	mappers.Mapper
}

func New() (c *Cartridge) {
	return &Cartridge{}
}

func (c *Cartridge) Load(name string) {

}
