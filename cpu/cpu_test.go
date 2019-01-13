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

func assertPC(expected, reg uint16, t *testing.T) {
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

func TestCpu(t *testing.T) {
	cpu := cpu.NewCPU()

	// Test status register bit toggling instructions first
	t.Run("test SEC", func(t *testing.T) {
		cpu.SEC(0x0)
		assertStatus("00X00001", cpu.Status, t)
	})

	t.Run("test CLC", func(t *testing.T) {
		cpu.CLC(0x0)
		assertStatus("00X00000", cpu.Status, t)
	})

	t.Run("test SEI", func(t *testing.T) {
		cpu.SEI(0x0)
		assertStatus("00X00100", cpu.Status, t)
	})

	t.Run("test CLI", func(t *testing.T) {
		cpu.CLI(0x0)
		assertStatus("00X00000", cpu.Status, t)
	})

	t.Run("test SED", func(t *testing.T) {
		cpu.SED(0x0)
		assertStatus("00X01000", cpu.Status, t)
	})

	t.Run("test CLD", func(t *testing.T) {
		cpu.CLD(0x0)
		assertStatus("00X00000", cpu.Status, t)
	})

	t.Run("test CLV", func(t *testing.T) {
		cpu.Status.V = true
		cpu.CLV(0x0)
		assertStatus("00X00000", cpu.Status, t)
	})
}
