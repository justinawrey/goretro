package main

import (
	"context"
	"embed"
	"log"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	// Create an instance of the app structure
	app := newApp()
	inputDriver := newWebviewInputDriver()
	displayDriver := newWebviewDisplayDriver()
	audioDriver := newWebviewAudioDriver()

	_, isDev := os.LookupEnv("WAILS_DEV")

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "goretro",
		Width:             512,
		Height:            512,
		MinWidth:          512,
		MinHeight:         512,
		MaxWidth:          512,
		MaxHeight:         512,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		Assets:            assets,
		Menu:              nil,
		Logger:            nil,
		LogLevel:          logger.DEBUG,
		OnStartup: func(ctx context.Context) {
			app.ctx = ctx
			inputDriver.ctx = ctx
			displayDriver.ctx = ctx
			audioDriver.ctx = ctx
		},
		OnDomReady:       app.domReady,
		OnBeforeClose:    app.beforeClose,
		OnShutdown:       app.shutdown,
		WindowStartState: options.Normal,
		Bind: []any{
			app,
			inputDriver,
			displayDriver,
			audioDriver,
		},
		EnumBind: []interface{}{
			joypads,
			buttons,
			displayEvents,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			// DisableFramelessWindowDecorations: false,
			WebviewUserDataPath: "",
		},
		// Mac platform specific options
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 true,
				HideToolbarSeparator:       true,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "goretro",
				Message: "",
				Icon:    icon,
			},
		},
		Debug: options.Debug{
			OpenInspectorOnStartup: isDev,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
