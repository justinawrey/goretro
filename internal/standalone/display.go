package main

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type RenderEvent string

const Render RenderEvent = "RENDER"

var displayEvents = []struct {
	Value  RenderEvent
	TSName string
}{
	{Render, "RENDER"},
}

var frameBuffer [256 * 256 * 4]int

type WebviewDisplayDriver struct {
	ctx context.Context
}

func newWebviewDisplayDriver() *WebviewDisplayDriver {
	return &WebviewDisplayDriver{}
}

func (w *WebviewDisplayDriver) renderFrame() {
	runtime.EventsEmit(w.ctx, string(Render), frameBuffer)
}
