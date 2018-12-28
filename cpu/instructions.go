package cpu

// SEC Set Carry Flag
func (cpu *CPU) SEC() {
	cpu.Status.C = true
}

// CLC Clear Carry Flag
func (cpu *CPU) CLC() {
	cpu.Status.C = false
}

// SEI Set Interrupt Disable
func (cpu *CPU) SEI() {
	cpu.Status.I = true
}

// CLI Clear Interrupt Disable
func (cpu *CPU) CLI() {
	cpu.Status.I = false
}

// SED Set Decimal Mode
func (cpu *CPU) SED() {
	cpu.Status.D = true
}

// CLD Clear Decimal Mode
func (cpu *CPU) CLD() {
	cpu.Status.D = false
}

// CLV Clear Overflow Flag
func (cpu *CPU) CLV() {
	cpu.Status.D = false
}
