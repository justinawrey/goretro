package core

// Component describes a component of the nes which can be deterministically initialized and cleared.
type component interface {
	// Init initializes the module to its correct start up state.
	init()

	// Clear clears the module to its correct power down state.
	clear()
}
