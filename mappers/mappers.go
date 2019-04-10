package mappers

const (
	mapperNROM = iota
)

type Mapper interface {
	ReadRegister(uint16) byte
	WriteRegister(uint16, byte)
}

func New(id int) (m Mapper) {
	switch id {
	default:
		fallthrough
	case 0:
		return &NROM{}
	}
}

type NROM struct {
}

func (nr *NROM) ReadRegister(address uint16) (data byte) {
	return 0x00
}

func (nr *NROM) WriteRegister(address uint16, data byte) {}
