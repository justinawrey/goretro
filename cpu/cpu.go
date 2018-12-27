package cpu

import "fmt"

// StatusRegister holds data for each status flag
// in the 6502 status register.
type StatusRegister struct {
	C bool // carry
	Z bool // zero result
	I bool // interrupt disable
	D bool // decimal mode
	B bool // break command
	V bool // overflow
	N bool // zero result
}

// String implements Stringer
func (sr *StatusRegister) String() string {
	convert := func(bit bool) string {
		if bit {
			return "1"
		}
		return "0"
	}

	return fmt.Sprintf("%s%s0%s%s%s%s%s",
		convert(sr.N),
		convert(sr.V),
		convert(sr.B),
		convert(sr.D),
		convert(sr.I),
		convert(sr.Z),
		convert(sr.C),
	)
}

// SEC Set Carry Flag
func SEC(sr *StatusRegister) {
	sr.C = true
}

// CLC Clear Carry Flag
func CLC(sr *StatusRegister) {
	sr.C = false
}
