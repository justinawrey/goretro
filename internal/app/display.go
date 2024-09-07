package app

import (
	"context"

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
	Ctx context.Context
}

func NewWebviewDisplayDriver() *WebviewDisplayDriver {
	return &WebviewDisplayDriver{}
}

func (w *WebviewDisplayDriver) renderFrame() {
	runtime.EventsEmit(w.Ctx, string(Render), frameBuffer)
}
