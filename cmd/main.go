package main

import (
	"flag"

	"github.com/justinawrey/goretro/internal/core"
	"github.com/justinawrey/goretro/internal/display"
	"github.com/justinawrey/goretro/internal/log"
)

func main() {
	var scale int

	flag.BoolVar(&log.Enabled, "debug", false, "run with debug logging enabled")
	flag.IntVar(&scale, "render-scale", 4, "display rendering scale")
	flag.Parse()

	nes := core.NewNes()
	if err := nes.UseCartridge("donkey-kong.nes"); err != nil {
		log.Log(err)
		return
	}

	display.NewDisplay(scale)
}
