// Package input provides an interface through which the caller
// can access input messages from the OS.
package input

// Input holds input information.
type Input struct{}

// NewInput creates a new Input.
func NewInput() (i *Input) {
	return &Input{}
}

// clear implements core.Component.
func (i *Input) clear() {
	panic("unimplemented")
}

// init implements core.Component.
func (i *Input) init() {
	panic("unimplemented")
}
