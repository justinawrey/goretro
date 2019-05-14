package main

import "github.com/justinawrey/goretro/nes"

func main() {
	nes := nes.New()
	nes.Load("/home/justin/rom/donkeykong.nes")
}
