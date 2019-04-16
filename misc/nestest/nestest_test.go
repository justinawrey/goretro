package nestest

import (
	"os"
	"testing"

	"github.com/justinawrey/nes/nes"
)

func TestNestest(t *testing.T) {
	// Set up logging
	nestest := nes.New()
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
