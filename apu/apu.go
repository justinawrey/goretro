package apu

type APU struct {
}

func New() (a *APU) {
	return &APU{}
}

func (a *APU) ReadRegister(address uint16) (data byte) { return 0x00 }
func (a *APU) WriteRegister(address uint16, data byte) {}
func (a *APU) Init()                                   {}
func (a *APU) Clear()                                  {}
