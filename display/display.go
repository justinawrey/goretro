// Package display provides functionality related to the UI of the emulator.
package display

// Display holds all information regarding the UI of the emulator.
// It is not emulation specific.
type Display struct{}

// New creates a new Display.
func New() (d *Display) { return &Display{} }

// Init implements nes.Module.
func (d *Display) Init() {}

// Clear implements nes.Module.
func (d *Display) Clear() {}
