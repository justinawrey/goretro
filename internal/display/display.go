// Package display provides functionality related to the UI of the emulator.
package display

// Display holds all information regarding the UI of the emulator.
// It is not emulation specific.
type Display struct{}

// NewDisplay creates a new Display.
func NewDisplay() (d *Display) {
	return &Display{}
}
