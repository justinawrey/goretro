package mappers

type Mapper struct {
}

func (m *Mapper) ReadRegister(address uint16) (data byte) {
	return 0x00
}

func (m *Mapper) WriteRegister(address uint16, data byte) {

}
