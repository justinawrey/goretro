package cpu_test

import (
	"errors"

	"github.com/justinawrey/nes/cpu"
)

const statusBitLen = 8

var validBits = map[byte]bool{'0': true, '1': true, 'X': true}

var errInvalidStatusString = errors.New("invalid format of status string")

func statusEquals(status string, sr *cpu.StatusRegister) (success bool, err error) {
	bitEquals := func(bitRepr byte, actual bool) (success bool, err error) {
		if _, ok := validBits[bitRepr]; !ok {
			return false, errInvalidStatusString
		}

		if (bitRepr == 'X') || (bitRepr == '1' && actual) || (bitRepr == '0' && !actual) {
			return true, nil
		}

		return false, nil
	}

	if len(status) != statusBitLen {
		return false, errInvalidStatusString
	}

	eqc, errc := bitEquals(status[0], sr.C)
	eqz, errz := bitEquals(status[1], sr.Z)
	eqi, erri := bitEquals(status[2], sr.I)
	eqd, errd := bitEquals(status[3], sr.D)
	eqb, errb := bitEquals(status[4], sr.B)
	eqv, errv := bitEquals(status[6], sr.V)
	eqn, errn := bitEquals(status[7], sr.N)

	for _, res := range []struct {
		err error
		eq  bool
	}{
		{errc, eqc},
		{errz, eqz},
		{erri, eqi},
		{errd, eqd},
		{errb, eqb},
		{errv, eqv},
		{errn, eqn},
	} {
		if res.err != nil {
			return false, err
		}
		if !res.eq {
			return false, nil
		}
	}
	return true, nil
}
