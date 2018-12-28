package cpu

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
