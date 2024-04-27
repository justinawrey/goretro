package main

import (
	"log"

	"github.com/justinawrey/goretro/internal/core"
)

func main() {
	nes := core.NewNes()

	if err := nes.UseCartridge("donkey-kong.nes"); err != nil {
		log.Fatalln(err)
	}

	nes.Start()
}
