# NRM-HA desktop gateway + Wails UI build targets.
#
# Wails v2 needs the `production` build tag (or `dev`) and, on macOS, the
# UniformTypeIdentifiers framework linked in. The Wails CLI does this
# automatically; plain `go build` does not, so this Makefile wraps it.

.PHONY: build run build-windows clean

# build compiles the desktop app + Alpaca server into a single binary (macOS).
build:
	CGO_LDFLAGS="-framework UniformTypeIdentifiers" go build -tags production -o nrm_ha_go .

# run builds (if needed) and runs the app.
run: build
	./nrm_ha_go

# build-windows compiles a Windows GUI executable without a console window.
# Run this on a Windows machine (Wails needs cgo for WebView2).
build-windows:
	go build -ldflags "-H=windowsgui" -tags production -o nrm_ha_go.exe .

# clean removes the built binaries.
clean:
	rm -f nrm_ha_go nrm_ha_go.exe
