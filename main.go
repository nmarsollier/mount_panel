package main

import (
	"context"
	"embed"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/nmarsollier/nrm_ha_go/internal/alpaca"
	"github.com/nmarsollier/nrm_ha_go/internal/mount"
	"github.com/nmarsollier/nrm_ha_go/internal/tools/env"
	"github.com/nmarsollier/nrm_ha_go/internal/tools/logx"
	"github.com/nmarsollier/nrm_ha_go/internal/ui"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	ctx := logx.CtxWithFields(context.Background(),
		logx.NewFields().Add(logx.CONTROLLER, "Main"))

	// Start the mount model: it connects and reconnects to the ESP in the
	// background, so the app runs even when the mount is not (yet) plugged in.
	mount.Init(env.Get().SerialPort)

	// Serve Alpaca (N.I.N.A./PHD2) in the background alongside the GUI.
	go alpaca.Start()

	app := ui.NewApp()

	err := wails.Run(&options.App{
		Title:  "NRM-HA",
		Width:  600,
		Height: 500,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.Startup,
		OnBeforeClose: func(ctx context.Context) bool {
			choice, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
				Type:          runtime.QuestionDialog,
				Title:         "Quit NRM-HA",
				Message:       "Are you sure you want to quit?",
				Buttons:       []string{"Quit", "Cancel"},
				DefaultButton: "Cancel",
				CancelButton:  "Cancel",
			})
			if err != nil {
				return true // on error, keep the app open
			}
			return choice != "Quit"
		},
		Bind: []any{app},
	})
	if err != nil {
		logx.Error(ctx, err)
		os.Exit(1)
	}
}
