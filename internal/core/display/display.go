// Package display provides functionality related to the UI of the emulator.
package display

import (
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"
)

const (
	width  = 256
	height = 240
)

var (
	white = pixel.RGB(1, 1, 1)
	black = pixel.RGB(0, 0, 0)
)

// Display holds all information regarding the UI of the emulator.
// It is not emulation specific.
type Display struct {
	*pixelgl.Window
	*imdraw.IMDraw

	scale int // the rendering scale of the display
}

func NewDisplay(scale int) (disp *Display) {
	return &Display{scale: scale}
}

// NewDisplay creates a new Display.
// func NewDisplay(scale int) {
// 	pixelgl.Run(func() {
// 		drawPx := func(x, y, scale int, d *Display) {
// 			lowerX := float64((x % width) * scale)
// 			lowerY := float64((height - y - 1) * d.scale)
//
// 			// TODO: use real colors
// 			if x%2 == 0 {
// 				d.IMDraw.Color = black
// 			} else {
// 				d.IMDraw.Color = white
// 			}
//
// 			d.IMDraw.Push(pixel.V(lowerX, lowerY))
// 			d.IMDraw.Push(pixel.V(lowerX+float64(d.scale), lowerY+float64(d.scale)))
// 			d.IMDraw.Rectangle(0)
//
// 		}
//
// 		cfg := pixelgl.WindowConfig{
// 			Title:  "goretro",
// 			Bounds: pixel.R(0, 0, float64(width*scale), float64(height*scale)),
// 			VSync:  true,
// 		}
//
// 		window, err := pixelgl.NewWindow(cfg)
// 		if err != nil {
// 			panic(err)
// 		}
//
// 		imd := imdraw.New(nil)
// 		d := &Display{window, imd, scale}
//
// 		for !d.Window.Closed() {
// 			d.IMDraw.Clear()
// 			d.IMDraw.Reset()
//
// 			// TODO: get real pixel data
// 			for y := 0; y < height; y++ {
// 				for x := 0; x < width; x++ {
// 					drawPx(x, y, scale, d)
// 				}
// 			}
//
// 			d.Window.Clear(white)
// 			d.IMDraw.Draw(d.Window)
// 			d.Window.Update()
// 		}
// 	})
// }
