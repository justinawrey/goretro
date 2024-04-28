package main

import (
	"flag"

	"github.com/justinawrey/goretro/internal/core"
	"github.com/justinawrey/goretro/internal/log"
)

func main() {
	flag.BoolVar(&log.Enabled, "debug", false, "run with debug logging enabled")
	flag.Parse()

	nes := core.NewNes()
	if err := nes.UseCartridge("donkey-kong.nes"); err != nil {
		log.Log(err)
		return
	}

	nes.Start()
}
