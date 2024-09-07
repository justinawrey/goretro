package app

import "context"

// App struct
type App struct {
	ContextHolder
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// domReady is called after front-end resources have been loaded
func (a App) DomReady(ctx context.Context) {
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) BeforeClose(ctx context.Context) bool {
	return false
}

// shutdown is called at application termination
func (a *App) Shutdown(ctx context.Context) {
	// Perform your teardown here
}
