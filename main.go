package main

import "github.com/justinawrey/goretro/internal/core"

func main() {
	nes := core.NewNes()
	nes.Load("/home/justin/rom/donkeykong.nes")
}
