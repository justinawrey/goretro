package cpu

import "fmt"

// Status holds data for each status flag
// in the 6502 status register.
type Status struct {
	C bool // carry
	Z bool // zero result
	I bool // interrupt disable
	D bool // decimal mode
	B bool // break command
	V bool // overflow
	N bool // zero result
}

// String implements Stringer
func (sr *Status) String() string {
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

// Clear clears all bits of the status register,
// i.e. will set all flags to '0'.
func (sr *Status) Clear() {
	sr.C = false
	sr.Z = false
	sr.I = false
	sr.D = false
	sr.B = false
	sr.V = false
	sr.N = false
}

// Registers holds data for each register
// used by the 6502,
type Registers struct {
	Status      Status
	Accumulator byte
}

// SEC Set Carry Flag
func SEC(r *Registers) {
	r.Status.C = true
}

// CLC Clear Carry Flag
func CLC(r *Registers) {
	r.Status.C = false
}

// SEI Set Interrupt Disable
func SEI(r *Registers) {
	r.Status.I = true
}

// CLI Clear Interrupt Disable
func CLI(r *Registers) {
	r.Status.I = false
}

// SED Set Decimal Mode
func SED(r *Registers) {
	r.Status.D = true
}

// CLD Clear Decimal Mode
func CLD(r *Registers) {
	r.Status.D = false
}

// CLV Clear Overflow Flag
func CLV(r *Registers) {
	r.Status.D = false
}
