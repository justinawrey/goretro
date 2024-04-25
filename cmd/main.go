package main

import (
	"github.com/justinawrey/goretro/internal/audio"
	"github.com/justinawrey/goretro/internal/core"
	"github.com/justinawrey/goretro/internal/display"
	"github.com/justinawrey/goretro/internal/input"
)

func main() {
	nes := core.NewNes()

	nes.UseDisplay(display.NewDisplay())
	nes.UseInput(input.NewInput())
	nes.UseAudio(audio.NewAudio())
	nes.UseCartridge("/home/justin/rom/donkeykong.nes")

	nes.Start()
}
