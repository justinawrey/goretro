package app

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type RenderEvent string

const Render RenderEvent = "RENDER"

var DisplayEvents = []struct {
	Value  RenderEvent
	TSName string
}{
	{Render, "RENDER"},
}

var frameBuffer [256 * 256 * 4]int

type WebviewDisplayDriver struct {
	ContextHolder
}

func NewWebviewDisplayDriver() *WebviewDisplayDriver {
	return &WebviewDisplayDriver{}
}

func (w *WebviewDisplayDriver) renderFrame() {
	runtime.EventsEmit(w.ctx, string(Render), frameBuffer)
}
