package main

import (
	"context"
	"embed"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/nmarsollier/mount_panel/internal/alpaca_mount"
	"github.com/nmarsollier/mount_panel/internal/eaf"
	"github.com/nmarsollier/mount_panel/internal/fan"
	"github.com/nmarsollier/mount_panel/internal/mount"
	"github.com/nmarsollier/mount_panel/internal/temp"
	"github.com/nmarsollier/mount_panel/internal/tools/env"
	"github.com/nmarsollier/mount_panel/internal/tools/logx"
	"github.com/nmarsollier/mount_panel/internal/ui"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	ctx := logx.CtxWithFields(context.Background(),
		logx.NewFields().Add(logx.CONTROLLER, "Main"))

	// Start the mount model: a single connection attempt at startup; the user
	// re-connects manually from the UI when it is unplugged.
	mount.Init(env.Get().SerialPort)

	// Start the accessory models (focuser, fan, temperature/vibration) each on
	// their own serial link.
	eaf.Init(env.Get().EAFSerialPort)
	fan.Init(env.Get().FanSerialPort)
	temp.Init(env.Get().TempSerialPort)

	// Serve Alpaca (N.I.N.A./PHD2) in the background alongside the GUI.
	go alpaca_mount.Start()

	app := ui.NewApp()

	err := wails.Run(&options.App{
		Title:  "Mount Panel",
		Width:  600,
		Height: 500,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.Startup,
		OnBeforeClose: func(ctx context.Context) bool {
			choice, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
				Type:          runtime.QuestionDialog,
				Title:         "Quit Mount Panel",
				Message:       "Are you sure you want to quit?",
				Buttons:       []string{"Quit", "Cancel"},
				DefaultButton: "Cancel",
				CancelButton:  "Cancel",
			})
			if err != nil {
				return true // on error, keep the app open
			}
			// On Windows the dialog ignores Buttons and shows a Yes/No box,
			// so the confirm answer comes back as "Yes" instead of "Quit".
			return choice != "Quit" && choice != "Yes"
		},
		Bind: []any{app},
	})
	if err != nil {
		logx.Error(ctx, err)
		os.Exit(1)
	}
}
