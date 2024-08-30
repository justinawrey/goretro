package main

import (
	"flag"

	"github.com/justinawrey/goretro/internal/core/audio"
	"github.com/justinawrey/goretro/internal/core/display"
	"github.com/justinawrey/goretro/internal/core/input"
	"github.com/justinawrey/goretro/internal/core/log"
	"github.com/justinawrey/goretro/internal/core/nes"
)

func main() {
	var scale int

	flag.BoolVar(&log.Enabled, "debug", false, "run with debug logging enabled")
	flag.IntVar(&scale, "render-scale", 4, "display rendering scale")
	flag.Parse()

	disp := display.NewDisplay(scale)
	input := input.NewInput()
	audio := audio.NewAudio()

	nes := nes.NewNes(disp, input, audio)
	if err := nes.UseCartridge("donkey-kong.nes"); err != nil {
		log.Log(err)
		return
	}
}
