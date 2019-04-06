package cpu_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/justinawrey/nes/cpu"
)

const statusBitLen = 8

var validBits = map[string]bool{"0": true, "1": true, "X": true}

var errInvalidStatusString = errors.New("invalid format of status string")

func statusRegisterDiff(want string, got *cpu.Status) string {
	return fmt.Sprintf(`       NV BDIZC
      ---------
want | %s
 got | %v 
`, want, got)
}

func statusEquals(status string, sr *cpu.Status) (success bool, err error) {
	bitEquals := func(bit string, actual bool) (success bool, err error) {
		if _, ok := validBits[bit]; !ok {
			return false, errInvalidStatusString
		}

		if (bit == "X") || (bit == "1" && actual) || (bit == "0" && !actual) {
			return true, nil
		}

		return false, nil
	}

	if len(status) != statusBitLen {
		return false, errInvalidStatusString
	}

	if status[2] != 'X' {
		return false, errInvalidStatusString
	}

	eqn, errn := bitEquals(string(status[0]), sr.N)
	eqv, errv := bitEquals(string(status[1]), sr.V)
	eqb, errb := bitEquals(string(status[3]), sr.B)
	eqd, errd := bitEquals(string(status[4]), sr.D)
	eqi, erri := bitEquals(string(status[5]), sr.I)
	eqz, errz := bitEquals(string(status[6]), sr.Z)
	eqc, errc := bitEquals(string(status[7]), sr.C)

	for _, res := range []struct {
		err error
		eq  bool
	}{
		{errn, eqn},
		{errv, eqv},
		{errb, eqb},
		{errd, eqd},
		{erri, eqi},
		{errz, eqz},
		{errc, eqc},
	} {
		if res.err != nil {
			return false, res.err
		}
		if !res.eq {
			return false, nil
		}
	}
	return true, nil
}

func assertStatus(status string, sr *cpu.Status, t *testing.T) {
	eq, err := statusEquals(status, sr)
	if err != nil {
		t.Error(err)
		return
	}
	if !eq {
		t.Errorf("\n%s\n", statusRegisterDiff(status, sr))
	}
}

func assertRegister(expected, reg byte, t *testing.T) {
	if reg != expected {
		t.Errorf("want: %v, got %v\n", expected, reg)
	}
}

func assertRegister16(expected, reg uint16, t *testing.T) {
	if reg != expected {
		t.Errorf("want: %v, got %v\n", expected, reg)
	}
}

func assertMemory(expected byte, loc uint16, mem *cpu.MemoryMap, t *testing.T) {
	got := mem[loc]
	if got != expected {
		t.Errorf("want: %v, got %v\n", expected, got)
	}
}

func TestInstructions(t *testing.T) {
	newCPUAndTest := func(inner func(c *cpu.CPU, t *testing.T)) func(*testing.T) {
		// Pre-computation closure
		return func(t *testing.T) {
			// All tests should be ran with a new CPU, in parallel
			t.Parallel()
			cpu := cpu.NewCPU()

			// Ignore any start-up state for these CPU tests
			cpu.Clear()

			// Run the actual test
			inner(cpu, t)
		}
	}

	// Test status register bit toggling instructions first
	t.Run("test SEC", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		cpu.SEC(0x0)
		assertStatus("00X00001", cpu.Status, t)
	}))

	t.Run("test CLC", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		cpu.CLC(0x0)
		assertStatus("00X00000", cpu.Status, t)
	}))

	t.Run("test SEI", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		cpu.SEI(0x0)
		assertStatus("00X00100", cpu.Status, t)
	}))

	t.Run("test CLI", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		cpu.CLI(0x0)
		assertStatus("00X00000", cpu.Status, t)
	}))

	t.Run("test SED", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		cpu.SED(0x0)
		assertStatus("00X01000", cpu.Status, t)
	}))

	t.Run("test CLD", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		cpu.CLD(0x0)
		assertStatus("00X00000", cpu.Status, t)
	}))

	t.Run("test CLV", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		cpu.Status.V = true
		cpu.CLV(0x0)
		assertStatus("00X00000", cpu.Status, t)
	}))

	t.Run("test INX", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// Should set negative flag, no zero flag
		cpu.X = 127
		cpu.INX(0x0)
		assertStatus("10X00000", cpu.Status, t)
		assertRegister(128, cpu.X, t)

		// Should set zero flag, no negative flag
		cpu.X = 255
		cpu.INX(0x0)
		assertStatus("00X00010", cpu.Status, t)
		assertRegister(0, cpu.X, t)
	}))

	t.Run("test INY", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// Should set negative flag, no zero flag
		cpu.Y = 127
		cpu.INY(0x0)
		assertStatus("10X00000", cpu.Status, t)
		assertRegister(128, cpu.Y, t)

		// Should set zero flag, no negative flag
		cpu.Y = 255
		cpu.INY(0x0)
		assertStatus("00X00010", cpu.Status, t)
		assertRegister(0, cpu.Y, t)
	}))

	t.Run("test INC", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// Should set negative flag, no zero flag
		var addr uint16 = 0x000A
		cpu.Write(addr, 127)
		cpu.INC(addr)
		assertStatus("10X00000", cpu.Status, t)
		assertMemory(128, addr, cpu.MemoryMap, t)

		// Should set zero flag, no negative flag
		cpu.Write(addr, 255)
		cpu.INC(addr)
		assertStatus("00X00010", cpu.Status, t)
		assertMemory(0, addr, cpu.MemoryMap, t)
	}))

	t.Run("test LDA", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// Should set negative flag, no zero flag
		var addr uint16 = 0x000A
		cpu.Write(addr, 0xAA)
		cpu.LDA(addr)
		assertStatus("10X00000", cpu.Status, t)
		assertRegister(0xAA, cpu.A, t)

		// Should set zero flag, no negative flag
		cpu.Write(addr, 0x00)
		cpu.LDA(addr)
		assertStatus("00X00010", cpu.Status, t)
		assertRegister(0x00, cpu.A, t)
	}))

	t.Run("test LDX", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// Should set negative flag, no zero flag
		var addr uint16 = 0x000A
		cpu.Write(addr, 0xAA)
		cpu.LDX(addr)
		assertStatus("10X00000", cpu.Status, t)
		assertRegister(0xAA, cpu.X, t)

		// Should set zero flag, no negative flag
		cpu.Write(addr, 0x00)
		cpu.LDX(addr)
		assertStatus("00X00010", cpu.Status, t)
		assertRegister(0x00, cpu.X, t)
	}))

	t.Run("test LDY", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// Should set negative flag, no zero flag
		var addr uint16 = 0x000A
		cpu.Write(addr, 0xAA)
		cpu.LDY(addr)
		assertStatus("10X00000", cpu.Status, t)
		assertRegister(0xAA, cpu.Y, t)

		// Should set zero flag, no negative flag
		cpu.Write(addr, 0x00)
		cpu.LDY(addr)
		assertStatus("00X00010", cpu.Status, t)
		assertRegister(0x00, cpu.Y, t)
	}))

	t.Run("test NOP", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		cpu.Status.C = true
		cpu.X = 0xAA

		// unused addr, NOP should not affect any of above
		cpu.NOP(0x00)
		assertStatus("00X00001", cpu.Status, t)
		assertRegister(0xAA, cpu.X, t)
	}))

	t.Run("test LSRA", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// Should set carry flag, zero flag
		// Accumulator should result in having data 0x00
		cpu.A = 0x01
		cpu.LSRA(0x00)
		assertStatus("00X00011", cpu.Status, t)
		assertRegister(0x00, cpu.A, t)

		// Should unset carry flag, zero flag
		// Accumulator should result in having data 0x01
		cpu.A = 0x02
		cpu.LSRA(0x00)
		assertStatus("00X00000", cpu.Status, t)
		assertRegister(0x01, cpu.A, t)

		// Carry flag should NOT rotate into bit 7 of accumulator
		cpu.Status.C = true
		cpu.A = 0x02
		cpu.LSRA(0x00)
		assertStatus("00X00000", cpu.Status, t)
		assertRegister(0x01, cpu.A, t)
	}))

	t.Run("test LSRM", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// Should set carry flag, zero flag
		// Memory at 0x000A should result in having data 0x00
		var addr uint16 = 0x000A
		cpu.Write(addr, 0x01)
		cpu.LSRM(addr)
		assertStatus("00X00011", cpu.Status, t)
		assertMemory(0x00, addr, cpu.MemoryMap, t)

		// Should unset carry flag, zero flag
		// Memory at 0x000A should result in having data 0x01
		cpu.Write(addr, 0x02)
		cpu.LSRM(addr)
		assertStatus("00X00000", cpu.Status, t)
		assertMemory(0x01, addr, cpu.MemoryMap, t)

		// Carry flag should NOT rotate into bit 7 of memory
		cpu.Status.C = true
		cpu.Write(addr, 0x02)
		cpu.LSRM(addr)
		assertStatus("00X00000", cpu.Status, t)
		assertMemory(0x01, addr, cpu.MemoryMap, t)
	}))

	t.Run("test ORA", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// Should set negative flag, no zero flag
		// Assert that 0x46 | 0xAA = 0xEE
		var addr uint16 = 0x000A
		cpu.Write(addr, 0x46)
		cpu.A = 0xAA
		cpu.ORA(addr)
		assertStatus("10X00000", cpu.Status, t)
		assertRegister(0xEE, cpu.A, t)

		// Should set zero flag, no negative flag
		// Assert that 0x00 | 0x00 = 0x00
		cpu.Write(addr, 0x00)
		cpu.A = 0x00
		cpu.ORA(addr)
		assertStatus("00X00010", cpu.Status, t)
		assertRegister(0x00, cpu.A, t)
	}))

	t.Run("test ASLA", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// Should set carry flag, unset zero flag
		// Accumulator should result in having data 0xFE
		cpu.A = 0xFF
		cpu.ASLA(0x00)
		assertStatus("10X00001", cpu.Status, t)
		assertRegister(0xFE, cpu.A, t)

		// Should unset carry flag, zero flag
		// Accumulator should result in having data 0x02
		cpu.A = 0x01
		cpu.ASLA(0x00)
		assertStatus("00X00000", cpu.Status, t)
		assertRegister(0x02, cpu.A, t)

		// Should set carry flag, zero flag
		cpu.A = 0x80
		cpu.ASLA(0x00)
		assertStatus("00X00011", cpu.Status, t)
		assertRegister(0x00, cpu.A, t)
	}))

	t.Run("test ASLM", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// Should set carry flag, unset zero flag
		// Memory at addr should result in having data 0xFE
		var addr uint16 = 0x000A
		cpu.Write(addr, 0xFF)
		cpu.ASLM(addr)
		assertStatus("10X00001", cpu.Status, t)
		assertMemory(0xFE, addr, cpu.MemoryMap, t)

		// Should unset carry flag, zero flag
		// Accumulator should result in having data 0x02
		cpu.Write(addr, 0x01)
		cpu.ASLM(addr)
		assertStatus("00X00000", cpu.Status, t)
		assertMemory(0x02, addr, cpu.MemoryMap, t)

		// Should set carry flag, zero flag
		cpu.Write(addr, 0x80)
		cpu.ASLM(addr)
		assertStatus("00X00011", cpu.Status, t)
		assertMemory(0x00, addr, cpu.MemoryMap, t)
	}))

	t.Run("test ROLA", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// Should set carry flag, unset zero flag
		// Accumulator should result in having data 0xFE
		cpu.A = 0xFF
		cpu.ROLA(0x00)
		assertStatus("00X00001", cpu.Status, t)
		assertRegister(0xFE, cpu.A, t)

		// Should unset carry flag, zero flag
		// Accumulator should result in having data 0x02
		cpu.A = 0x01
		cpu.ROLA(0x00)
		assertStatus("00X00000", cpu.Status, t)
		assertRegister(0x02, cpu.A, t)

		// Should set carry flag, unset zero flag, contents of old carry
		// should rotate left into bit 0 of accumulator value
		cpu.Status.C = true
		cpu.A = 0x80
		cpu.ROLA(0x00)
		assertStatus("00X00001", cpu.Status, t)
		assertRegister(0x01, cpu.A, t)
	}))

	t.Run("test ROLM", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// Should set carry flag, unset zero flag
		// Memory should result in having data 0xFE
		var addr uint16 = 0x000A
		cpu.Write(addr, 0xFF)
		cpu.ROLM(addr)
		assertStatus("00X00001", cpu.Status, t)
		assertMemory(0xFE, addr, cpu.MemoryMap, t)

		// Should unset carry flag, zero flag
		// Memory should result in having data 0x02
		cpu.Write(addr, 0x01)
		cpu.ROLM(addr)
		assertStatus("00X00000", cpu.Status, t)
		assertMemory(0x02, addr, cpu.MemoryMap, t)

		// Should set carry flag, unset zero flag, contents of old carry
		// should rotate left into bit 0 of memory value
		cpu.Status.C = true
		cpu.Write(addr, 0x80)
		cpu.ROLM(addr)
		assertStatus("00X00001", cpu.Status, t)
		assertMemory(0x01, addr, cpu.MemoryMap, t)
	}))

	t.Run("test RORA", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// Should set carry flag, unset zero flag
		// Accumulator should result in having data 0x7F
		cpu.A = 0xFF
		cpu.RORA(0x00)
		assertStatus("00X00001", cpu.Status, t)
		assertRegister(0x7F, cpu.A, t)

		// Should unset carry flag, zero flag
		// Accumulator should result in having data 0x01
		cpu.A = 0x02
		cpu.RORA(0x00)
		assertStatus("00X00000", cpu.Status, t)
		assertRegister(0x01, cpu.A, t)

		// Should set carry flag, unset zero flag, contents of old carry
		// should rotate right into bit 7 of accumulator value
		cpu.Status.C = true
		cpu.A = 0x01
		cpu.RORA(0x00)
		assertStatus("00X00001", cpu.Status, t)
		assertRegister(0x80, cpu.A, t)
	}))

	t.Run("test RORM", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// Should set carry flag, unset zero flag
		// Memory should result in having data 0x7F
		var addr uint16 = 0x000A
		cpu.Write(addr, 0xFF)
		cpu.RORM(addr)
		assertStatus("00X00001", cpu.Status, t)
		assertMemory(0x7F, addr, cpu.MemoryMap, t)

		// Should unset carry flag, zero flag
		// Memory should result in having data 0x01
		cpu.Write(addr, 0x02)
		cpu.RORM(addr)
		assertStatus("00X00000", cpu.Status, t)
		assertMemory(0x01, addr, cpu.MemoryMap, t)

		// Should set carry flag, unset zero flag, contents of old carry
		// should rotate right into bit 7 of memory value
		cpu.Status.C = true
		cpu.Write(addr, 0x01)
		cpu.RORM(addr)
		assertStatus("00X00001", cpu.Status, t)
		assertMemory(0x80, addr, cpu.MemoryMap, t)
	}))

	t.Run("test STA", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		cpu.Status.C = true
		cpu.Status.Z = true
		cpu.Status.N = true
		var addr uint16 = 0x000A

		// Status flags should be left unchanged, memory at addr
		// should have contents of accumulator
		cpu.A = 0xAA
		cpu.STA(addr)
		assertStatus("10X00011", cpu.Status, t)
		assertMemory(0xAA, addr, cpu.MemoryMap, t)
	}))

	t.Run("test STX", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		cpu.Status.C = true
		cpu.Status.Z = true
		cpu.Status.N = true
		var addr uint16 = 0x000A

		// Status flags should be left unchanged, memory at addr
		// should have contents of X register
		cpu.X = 0xBB
		cpu.STX(addr)
		assertStatus("10X00011", cpu.Status, t)
		assertMemory(0xBB, addr, cpu.MemoryMap, t)
	}))

	t.Run("test STY", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		cpu.Status.C = true
		cpu.Status.Z = true
		cpu.Status.N = true
		var addr uint16 = 0x000A

		// Status flags should be left unchanged, memory at addr
		// should have contents of Y register
		cpu.Y = 0xCC
		cpu.STY(addr)
		assertStatus("10X00011", cpu.Status, t)
		assertMemory(0xCC, addr, cpu.MemoryMap, t)
	}))

	t.Run("test TAX", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// Zero flag should be set
		cpu.A = 0x00
		cpu.TAX(0x00)
		assertStatus("00X00010", cpu.Status, t)
		assertRegister(0x00, cpu.X, t)

		// Zero flag should be unset, negative flag set
		cpu.A = 0xF0
		cpu.TAX(0x00)
		assertStatus("10X00000", cpu.Status, t)
		assertRegister(0xF0, cpu.X, t)
	}))

	t.Run("test TAY", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// Zero flag should be set
		cpu.A = 0x00
		cpu.TAY(0x00)
		assertStatus("00X00010", cpu.Status, t)
		assertRegister(0x00, cpu.Y, t)

		// Zero flag should be unset, negative flag set
		cpu.A = 0xF0
		cpu.TAY(0x00)
		assertStatus("10X00000", cpu.Status, t)
		assertRegister(0xF0, cpu.Y, t)
	}))

	t.Run("test TSX", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// Zero flag should be set
		cpu.SP = 0x00
		cpu.TSX(0x00)
		assertStatus("00X00010", cpu.Status, t)
		assertRegister(0x00, cpu.X, t)

		// Zero flag should be unset, negative flag set
		cpu.SP = 0xF0
		cpu.TSX(0x00)
		assertStatus("10X00000", cpu.Status, t)
		assertRegister(0xF0, cpu.X, t)
	}))

	t.Run("test TXA", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// Zero flag should be set
		cpu.X = 0x00
		cpu.TXA(0x00)
		assertStatus("00X00010", cpu.Status, t)
		assertRegister(0x00, cpu.A, t)

		// Zero flag should be unset, negative flag set
		cpu.X = 0xF0
		cpu.TXA(0x00)
		assertStatus("10X00000", cpu.Status, t)
		assertRegister(0xF0, cpu.A, t)
	}))

	t.Run("test TXS", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// Zero flag should be set
		cpu.X = 0x00
		cpu.TXS(0x00)
		assertStatus("00X00010", cpu.Status, t)
		assertRegister(0x00, cpu.SP, t)

		// Zero flag should be unset, negative flag set
		cpu.X = 0xF0
		cpu.TXS(0x00)
		assertStatus("10X00000", cpu.Status, t)
		assertRegister(0xF0, cpu.SP, t)
	}))

	t.Run("test TYA", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// Zero flag should be set
		cpu.Y = 0x00
		cpu.TYA(0x00)
		assertStatus("00X00010", cpu.Status, t)
		assertRegister(0x00, cpu.A, t)

		// Zero flag should be unset, negative flag set
		cpu.Y = 0xF0
		cpu.TYA(0x00)
		assertStatus("10X00000", cpu.Status, t)
		assertRegister(0xF0, cpu.A, t)
	}))

	t.Run("test ADC", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		t.SkipNow()
	}))

	t.Run("test AND", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// Should set zero flag, unset negative flag
		// Accumulator should result in having data 0x00 because 0xAA & 0x55 == 0x00
		var addr uint16 = 0x0A
		cpu.A = 0xAA
		cpu.Write(addr, 0x55)
		assertStatus("00X00010", cpu.Status, t)
		assertRegister(0x00, cpu.A, t)

		// Should unset zero flag, set negative flag
		// Accumulator should result in having data
		cpu.A = 0x88
		cpu.Write(addr, 0xFF)
		assertStatus("10X00000", cpu.Status, t)
		assertRegister(0x88, cpu.A, t)
	}))

	t.Run("test BCC", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// flags should not change throughout
		cpu.Status.N = true
		cpu.Status.C = true
		var start uint16 = 0x01
		var dest uint16 = 0x02

		// If carry set, don't branch
		// Should not affect status register
		cpu.PC = start
		cpu.BCC(dest)
		assertStatus("10X00001", cpu.Status, t)
		assertRegister16(start, cpu.PC, t)

		// Otherwise, branch
		// Should not affect status register
		cpu.Status.C = false
		cpu.PC = start
		cpu.BCC(dest)
		assertStatus("10X00000", cpu.Status, t)
		assertRegister16(dest, cpu.PC, t)
	}))

	t.Run("test BCS", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// flags should not change throughout
		cpu.Status.N = true
		cpu.Status.C = true
		var start uint16 = 0x01
		var dest uint16 = 0x02

		// If carry set, branch
		// Should not affect status register
		cpu.PC = start
		cpu.BCS(dest)
		assertStatus("10X00001", cpu.Status, t)
		assertRegister16(dest, cpu.PC, t)

		// Otherwise, don't branch
		// Should not affect status register
		cpu.Status.C = false
		cpu.PC = start
		cpu.BCS(dest)
		assertStatus("10X00000", cpu.Status, t)
		assertRegister16(start, cpu.PC, t)
	}))

	t.Run("test BEQ", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// flags should not change throughout
		cpu.Status.N = true
		cpu.Status.Z = true
		var start uint16 = 0x01
		var dest uint16 = 0x02

		// If zero flag set, branch
		// Should not affect status register
		cpu.PC = start
		cpu.BEQ(dest)
		assertStatus("10X00010", cpu.Status, t)
		assertRegister16(dest, cpu.PC, t)

		// Otherwise, don't branch
		// Should not affect status register
		cpu.Status.Z = false
		cpu.PC = start
		cpu.BEQ(dest)
		assertStatus("10X00000", cpu.Status, t)
		assertRegister16(start, cpu.PC, t)
	}))

	t.Run("test BIT", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		t.SkipNow()
	}))

	t.Run("test BMI", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// flags should not change throughout
		cpu.Status.N = true
		cpu.Status.Z = true
		var start uint16 = 0x01
		var dest uint16 = 0x02

		// If negative flag set, branch
		// Should not affect status register
		cpu.PC = start
		cpu.BMI(dest)
		assertStatus("10X00010", cpu.Status, t)
		assertRegister16(dest, cpu.PC, t)

		// Otherwise, don't branch
		// Should not affect status register
		cpu.Status.N = false
		cpu.PC = start
		cpu.BMI(dest)
		assertStatus("00X00010", cpu.Status, t)
		assertRegister16(start, cpu.PC, t)
	}))

	t.Run("test BNE", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		// flags should not change throughout
		cpu.Status.N = true
		cpu.Status.Z = true
		var start uint16 = 0x01
		var dest uint16 = 0x02

		// If zero flag set, don't branch
		// Should not affect status register
		cpu.PC = start
		cpu.BNE(dest)
		assertStatus("10X00010", cpu.Status, t)
		assertRegister16(start, cpu.PC, t)

		// Otherwise, branch
		// Should not affect status register
		cpu.Status.Z = false
		cpu.PC = start
		cpu.BNE(dest)
		assertStatus("10X00000", cpu.Status, t)
		assertRegister16(dest, cpu.PC, t)
	}))

	t.Run("test BPL", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		t.SkipNow()
	}))

	t.Run("test BRK", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		t.SkipNow()
	}))

	t.Run("test BVC", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		t.SkipNow()
	}))

	t.Run("test BVS", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		t.SkipNow()
	}))

	t.Run("test CMP", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		t.SkipNow()
	}))

	t.Run("test CPX", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		t.SkipNow()
	}))

	t.Run("test CPY", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {

	}))

	t.Run("test DEC", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		t.SkipNow()
	}))

	t.Run("test DEX", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		t.SkipNow()
	}))

	t.Run("test DEY", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		t.SkipNow()
	}))

	t.Run("test EOR", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		t.SkipNow()
	}))

	t.Run("test JMP", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		t.SkipNow()
	}))

	t.Run("test JSR", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		t.SkipNow()
	}))

	t.Run("test PHA", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		t.SkipNow()
	}))

	t.Run("test PHP", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		t.SkipNow()
	}))

	t.Run("test PLA", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		t.SkipNow()
	}))

	t.Run("test PLP", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		t.SkipNow()
	}))

	t.Run("test RTI", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		t.SkipNow()
	}))

	t.Run("test RTS", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		t.SkipNow()
	}))

	t.Run("test SBC", newCPUAndTest(func(cpu *cpu.CPU, t *testing.T) {
		t.SkipNow()
	}))
}

func TestAddressingModes(t *testing.T) {
	t.SkipNow()
}
