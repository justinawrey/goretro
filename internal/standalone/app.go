package main

import (
	"context"
	"math/rand"
)

const displayBufSize = 256 * 256 * 4

var stream [displayBufSize]int

func init() {
	for i := 0; i < len(stream); i += 4 {
		stream[i] = rand.Intn(256)
		stream[i+1] = rand.Intn(256)
		stream[i+2] = rand.Intn(256)
		stream[i+3] = 255
	}
}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
	a.InitializeInputListeners()
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) bool {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

func (a *App) RequestFrame() [displayBufSize]int {
	return stream
}
