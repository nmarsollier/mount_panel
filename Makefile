# Mount Panel desktop gateway + Wails UI build targets.
#
# Wails v2 needs the `production` build tag (or `dev`) and, on macOS, the
# UniformTypeIdentifiers framework linked in. The Wails CLI does this
# automatically; plain `go build` does not, so this Makefile wraps it.

.PHONY: build run build-windows winres clean

# build compiles the desktop app + Alpaca server into a single binary (macOS).
build:
	CGO_LDFLAGS="-framework UniformTypeIdentifiers" go build -tags production -o mount_panel .

# run builds (if needed) and runs the app.
run: build
	./mount_panel

# build-windows compiles a Windows GUI executable without a console window.
# Run this on a Windows machine (Wails needs cgo for WebView2). The window and
# .exe icon is embedded via rsrc_windows_amd64.syso (see `winres`).
build-windows:
	go build -ldflags "-H=windowsgui" -tags production -o mount_panel.exe .

# winres regenerates the embedded Windows icon resource (.syso) from
# winres/winres.json (icon at resource ID 3, which Wails expects for the title
# bar). Requires go-winres:
#   go install github.com/tc-hib/go-winres@latest
winres:
	go-winres make --arch amd64

# clean removes the built binaries.
clean:
	rm -f mount_panel mount_panel.exe
