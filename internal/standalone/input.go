package main

import "context"

type Button string

const (
	A      Button = "A"
	B      Button = "B"
	Start  Button = "START"
	Select Button = "SELECT"
	Up     Button = "UP"
	Right  Button = "RIGHT"
	Down   Button = "DOWN"
	Left   Button = "LEFT"
)

type Joypad string

const (
	Primary   Joypad = "PRIMARY"
	Secondary Joypad = "SECONDARY"
)

var buttons = []struct {
	Value  Button
	TSName string
}{
	{A, "A"},
	{B, "B"},
	{Start, "START"},
	{Select, "SELECT"},
	{Up, "UP"},
	{Right, "RIGHT"},
	{Down, "DOWN"},
	{Left, "LEFT"},
}

var joypads = []struct {
	Value  Joypad
	TSName string
}{
	{Primary, "PRIMARY"},
	{Secondary, "SECONDARY"},
}

type WebviewInputDriver struct {
	ctx     context.Context
	joypad1 map[Button]bool
	joypad2 map[Button]bool
}

func newWebviewInputDriver() *WebviewInputDriver {
	return &WebviewInputDriver{
		joypad1: make(map[Button]bool),
		joypad2: make(map[Button]bool),
	}
}

func (w *WebviewInputDriver) setJoypadState(joypad Joypad, button Button, to bool) {
	switch joypad {
	case Primary:
		w.joypad1[button] = to
	case Secondary:
		w.joypad2[button] = to
	}
}

func (w *WebviewInputDriver) getJoypadState(joypad Joypad, button Button) bool {
	inner := func(j map[Button]bool, button Button) bool {
		pressed, ok := j[button]
		if !ok {
			return false
		}

		return pressed
	}

	var pressed bool
	switch joypad {
	case Primary:
		pressed = inner(w.joypad1, button)
	case Secondary:
		pressed = inner(w.joypad2, button)
	}

	return pressed
}

func (w *WebviewInputDriver) SetButton(joypad Joypad, button Button, to bool) {
	w.setJoypadState(joypad, button, to)
}

func (w *WebviewInputDriver) getButton(joypad Joypad, button Button) bool {
	return w.getJoypadState(joypad, button)
}
