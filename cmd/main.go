package main

import (
	"github.com/justinawrey/goretro/internal/core"
	"github.com/justinawrey/goretro/internal/log"
)

func main() {
	nes := core.NewNes()

	if err := nes.UseCartridge("donkey-kong.nes"); err != nil {
		log.Log(err)
		return
	}

	nes.Start()
}
