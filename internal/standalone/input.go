package main

func (a *App) setJoypadState(button string, to bool) {
	a.joypadState[button] = to
}

func (a *App) getJoypadState(button string) bool {
	pressed, ok := a.joypadState[button]
	if !ok {
		return false
	}

	return pressed
}

func (a *App) SetButtonA(to bool) {
	a.setJoypadState(buttonA, to)
}

func (a *App) SetButtonB(to bool) {
	a.setJoypadState(buttonB, to)
}

func (a *App) SetButtonStart(to bool) {
	a.setJoypadState(buttonStart, to)
}

func (a *App) SetButtonSelect(to bool) {
	a.setJoypadState(buttonSelect, to)
}

func (a *App) SetButtonUp(to bool) {
	a.setJoypadState(buttonUp, to)
}

func (a *App) SetButtonRight(to bool) {
	a.setJoypadState(buttonRight, to)
}

func (a *App) SetButtonDown(to bool) {
	a.setJoypadState(buttonDown, to)
}

func (a *App) SetButtonLeft(to bool) {
	a.setJoypadState(buttonLeft, to)
}

func (a *App) getButtonA(to bool) bool {
	return a.getJoypadState(buttonA)
}

func (a *App) getButtonB(to bool) bool {
	return a.getJoypadState(buttonB)
}

func (a *App) getButtonStart(to bool) bool {
	return a.getJoypadState(buttonStart)
}

func (a *App) getButtonSelect(to bool) bool {
	return a.getJoypadState(buttonSelect)
}

func (a *App) getButtonUp(to bool) bool {
	return a.getJoypadState(buttonUp)
}

func (a *App) getButtonRight(to bool) bool {
	return a.getJoypadState(buttonRight)
}

func (a *App) getButtonDown(to bool) bool {
	return a.getJoypadState(buttonDown)
}

func (a *App) getButtonLeft(to bool) bool {
	return a.getJoypadState(buttonLeft)
}
