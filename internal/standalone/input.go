package main

import (
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	RECEIVE_INPUT_A      = "receive-input-a"
	RECEIVE_INPUT_B      = "receive-input-b"
	RECEIVE_INPUT_SELECT = "receive-input-select"
	RECEIVE_INPUT_START  = "receive-input-start"
	RECEIVE_INPUT_UP     = "receive-input-up"
	RECEIVE_INPUT_RIGHT  = "receive-input-right"
	RECEIVE_INPUT_DOWN   = "receive-input-down"
	RECEIVE_INPUT_LEFT   = "receive-input-left"

	REQUEST_INPUT_A      = "request-input-a"
	REQUEST_INPUT_B      = "request-input-b"
	REQUEST_INPUT_SELECT = "request-input-select"
	REQUEST_INPUT_START  = "request-input-start"
	REQUEST_INPUT_UP     = "request-input-up"
	REQUEST_INPUT_RIGHT  = "request-input-right"
	REQUEST_INPUT_DOWN   = "request-input-down"
	REQUEST_INPUT_LEFT   = "request-input-left"
)

func (a *App) InitializeInputListeners() {
	runtime.EventsOn(a.ctx, RECEIVE_INPUT_A, LogInput)
	runtime.EventsOn(a.ctx, RECEIVE_INPUT_B, LogInput)
	runtime.EventsOn(a.ctx, RECEIVE_INPUT_SELECT, LogInput)
	runtime.EventsOn(a.ctx, RECEIVE_INPUT_START, LogInput)
	runtime.EventsOn(a.ctx, RECEIVE_INPUT_UP, LogInput)
	runtime.EventsOn(a.ctx, RECEIVE_INPUT_RIGHT, LogInput)
	runtime.EventsOn(a.ctx, RECEIVE_INPUT_DOWN, LogInput)
	runtime.EventsOn(a.ctx, RECEIVE_INPUT_LEFT, LogInput)
}

func (a *App) ReadA() {
	runtime.EventsEmit(a.ctx, REQUEST_INPUT_A)
}

func (a *App) ReadB() {
	runtime.EventsEmit(a.ctx, REQUEST_INPUT_B)
}

func (a *App) ReadSelect() {
	runtime.EventsEmit(a.ctx, REQUEST_INPUT_SELECT)
}

func (a *App) ReadStart() {
	runtime.EventsEmit(a.ctx, REQUEST_INPUT_START)
}

func (a *App) ReadUp() {
	runtime.EventsEmit(a.ctx, REQUEST_INPUT_UP)
}

func (a *App) ReadRight() {
	runtime.EventsEmit(a.ctx, REQUEST_INPUT_RIGHT)
}

func (a *App) ReadDown() {
	runtime.EventsEmit(a.ctx, REQUEST_INPUT_DOWN)
}

func (a *App) ReadLeft() {
	runtime.EventsEmit(a.ctx, REQUEST_INPUT_LEFT)
}

func LogInput(args ...any) {
	fmt.Println(args)
}
