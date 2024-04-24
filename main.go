package main

import "github.com/justinawrey/goretro/internal/nes"

func main() {
	nes := nes.New()
	nes.Load("/home/justin/rom/donkeykong.nes")
}
