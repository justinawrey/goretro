package nestest

import (
	"os"
	"testing"

	nes "github.com/justinawrey/goretro/internal/core"
)

func TestNestest(t *testing.T) {
	// Set up logging
	nestest := nes.NewNes()
	nestest.Load("nestest.nes")
	fi, err := os.Create("output.log")
	if err != nil {
		t.Error("couldn't create output log file")
		return
	}
	nestest.OutputTo(fi)

	// Start stepping through
	nestest.Start()
}
