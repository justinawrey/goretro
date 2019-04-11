package nestest

import (
	"testing"

	"github.com/justinawrey/nes"
)

func TestNestest(t *testing.T) {
	nestest := nes.New()
	nestest.Load("nestest.nes")
}
