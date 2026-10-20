# mount_panel

Desktop gateway and UI for the [bft_firmware](https://github.com/nmarsollier/bft_firmware) mount.

It is the **desktop app** that replaces the HTTP/Alpaca server currently running
inside the ESP32-S3. It talks to the mount over a **USB CDC-ACM serial protocol**
and exposes **ASCOM Alpaca** (for N.I.N.A. and PHD2) plus a **desktop UI built
with Wails**.

---

## The idea

The ESP32-S3 chokes on the network (TCP/HTTP over USB-NCM): it is slow and hangs
on pointless operations. The solution is to take the network off the ESP and move
it to this app:

```
N.I.N.A. / PHD2 ──HTTP──▶ Alpaca REST (GIN) ──┐
                                               ├─▶ Mount ──▶ Proto ──▶ ESP32-S3
Desktop UI (Wails) ──bindings──▶ ui ───────────┘
```

- **The ESP** is the authority over the mount: it publishes its state and executes
  a handful of operations. It no longer knows anything about Alpaca, HTTP or
  ASCOM.
- **This app** is the integration: it translates that into Alpaca and the UI, and
  handles everything that is client-only.

The full protocol spec lives in the firmware:
[`main/proto/README.md`](https://github.com/nmarsollier/bft_firmware/blob/main/main/proto/README.md).

---

## The protocol it consumes

5 message families, frames `[u32 little-endian length][JSON UTF-8]`:

| Family | Role |
| --- | --- |
| `capabilities` | identity + static capabilities (read once on connect) |
| `state` | a coherent snapshot of the whole mount |
| `config` | read/edit persistent configuration (site, time, guide rates) |
| `control` | desired continuous-motion state (tracking + manual) |
| `action` | an operation with a beginning and an end (`stop`, `goto`, `guide`, `home`, `park`, `unpark`, `sync`†, `reset`, `move`, `limits`) |

Each request carries an `id` (echoed in the response) and `boot_id` identifies
reboots. `CONTROL` honours "omitted field = keep; speed 0 = stop".

---

## What this app exposes

### Alpaca (ITelescopeV4)

The app translates each Alpaca call into the protocol:

| Alpaca request | Resolution |
| --- | --- |
| RA, DEC, tracking, slewing, park, guiding | `state` fields |
| Capabilities (`can*`) | `capabilities` |
| `MoveAxis` | axis `control` |
| Change tracking | `control` |
| `PulseGuide` | `action: guide` |
| GoTo, home, park, unpark | the matching `action` |
| Target coordinates | app state until the GoTo runs |
| Client/transaction identifiers | managed only by this app |
| Alpaca errors | translation of internal errors |

### Desktop UI (Wails)

Manual mount control: a native window that reuses the firmware's Alpine.js SPA,
calling `mount` through **Wails bindings** (no HTTP). See the
[spec](#desktop-ui-spec-wails) below.

### App responsibilities

- **Connect to the ESP** (detection + reconnect): enumerates serial ports and does
  a `capabilities` handshake to verify it is a mount; if it disconnects, it
  reconnects automatically. Works on macOS and Windows.
- **Provide the time (SNTP)**: the ESP no longer syncs time (it has no network).
  This app sends the time to the ESP via `config.utc` on connect, on a new
  `boot_id`, when `state.time_valid` is false, and periodically (~1 h). Details in
  the migration section.
- **No state cache**: every Alpaca property read fetches a fresh `state` snapshot
  from the ESP over the serial link (like the firmware reads the motors
  directly).
- **UDP discovery** for Alpaca (port 32227), so clients can find the server.

---

## Migration: what leaves the ESP and who does it here

Inventory of what the firmware currently initialises/exposes and what falls to the
Go gateway once `usb_net` (network) and `www` (SPA) are removed. The rule:
**everything network/client/UI migrates to Go; everything hardware/motors stays on
the ESP**.

### Firmware init today (`runtime_setup.c` + `main.c`)

| Current init | Stays or migrates? | Who does it afterwards |
| --- | --- | --- |
| `nvs_flash_init` (settings persistence) | stays | ESP (NVS) |
| `led_init`, `buzzer_init` | stays | ESP (GPIO/LEDC) |
| `mount_init` (loads settings) | stays | ESP |
| `power_init` (12 V ADC) | stays | ESP |
| `motors_init`, `motors_update_power` | stays | ESP |
| `usb_cdc_init` + `proto_init` | stays | ESP |
| `esp_netif_init` + `esp_event_loop` | **removed** | — (no network) |
| `usb_net_init` (NCM) | **removed** | — (replaced by CDC) |
| `sntp_start` | **migrates to Go** | Go sends `config.utc` |
| `rest_server_start` (REST/UI) | **migrates to Go** | Go serves the UI |
| `rest_alpaca_server_start` (Alpaca) | **migrates to Go** | Go serves Alpaca |
| `udp_alpaca_start` (discovery) | **migrates to Go** | Go answers discovery |

### SNTP / time (detail)

Today `sntp/sntp_start.c` syncs the ESP clock via SNTP (NTP) over the network:
host DNS + `pool.ntp.org`, re-sync every 1 h, which enables `time_valid` (GOTO
requires a valid time).

Without the network, **the ESP has no time source**. The Go gateway provides it:

1. **On connect** (and on a new `boot_id`), it sends `config.utc` with the current
   time.
2. **Periodically** (every ~1 h, like the firmware did) it re-sends `config.utc`
   to correct the ESP RTC drift.

The time source in Go is the **host system clock** (macOS/Linux keep it synced via
NTP). No dedicated SNTP client is needed unless the host clock is untrusted; in
that case Go could run its own SNTP, but normally `time.Now().UTC()` is enough.

### New gateway responsibilities (things the firmware did that must be recreated)

| What | Status |
| --- | --- |
| Time (`config.utc`) | ✅ implemented (`mount` provides time + re-sync) |
| Alpaca (ITelescopeV4) | ✅ implemented in `internal/alpaca` |
| Desktop UI | ✅ implemented in `internal/ui` (Wails) |
| UDP discovery (32227) | ✅ implemented in `internal/alpaca` |
| Serial detection/reconnect | ✅ implemented in `internal/proto` + `internal/mount` |
| **Reset/reboot** the ESP | ✅ `action.reset` |
| **Set current position as** (home/limits) | ✅ `action.limits` |
| **Move N degrees relative** (Manual Slew) | ✅ `action.move` |

Every desktop-UI command has a protocol equivalent. See the protocol at
[`bft_firmware/main/proto/README.md`](https://github.com/nmarsollier/bft_firmware/blob/main/main/proto/README.md).

## Internal architecture (packages)

```
main.go                     # entry point: init mount + Alpaca (goroutine) + Wails
frontend/                   # embedded desktop UI (HTML/CSS/JS Alpine.js)
internal/
  proto/                    # serial protocol client (framing + detection)
  mount/                    # mount model (connection + operations, no cache)
  alpaca/                   # Alpaca REST server (GIN): ITelescopeV4 endpoints
    server/                 # GIN engine + logging middleware
  ui/                       # Wails binding (App) → mount
  tools/                    # logx (logger) + env (config), cross-cutting
test/alpaca/                # Alpaca integration tests (build tag `integration`)
```

**Functional** style: package-level functions, one file per use case, no
dependency injection. See [`CLAUDE.md`](CLAUDE.md).

---

## Prerequisites (fresh install)

Everything needed to compile and run from a clean machine. The app is a single
Go module (`go 1.25.0`) using Wails v2; every Go dependency is fetched by
`go build` from `go.mod`, so beyond a Go toolchain you only need the platform
pieces below.

### Common

- **Go 1.25+** — required by `go.mod`.
  - macOS: `brew install go`
  - Windows: `winget install GoLang.Go` (or the installer from <https://go.dev/dl/>)
- **Node.js + npm** — only to rebuild `frontend/dist/index.html` from
  `frontend/` (`cd frontend && npm run build`). That file is committed, so a
  plain build does **not** need Node.
  - macOS: `brew install node`
  - Windows: `winget install OpenJS.NodeJS.LTS` (or <https://nodejs.org/>)

### macOS

Wails and the serial-port enumerator use **cgo** and link the macOS SDK, so you
need **Xcode Command Line Tools** (clang + the SDK frameworks: WebKit, Cocoa,
UniformTypeIdentifiers, IOKit):

```sh
xcode-select --install
```

Requires macOS 11+ (for the `UniformTypeIdentifiers` framework). The Makefile
already passes the extra `-framework UniformTypeIdentifiers` linker flag, so
nothing else is needed.

### Windows

**No C compiler is needed**: the WebView2, COM and serial-port dependencies are
all pure Go (Wails v2 loads WebView2 through the pure-Go `go-winloader`), so
`go build` works straight from a Go install.

The app needs the **WebView2 runtime** at *runtime* (not to compile):

- **Windows 11**: already installed.
- **Windows 10**: usually installed via Windows Update; otherwise install the
  "Evergreen" runtime from
  <https://developer.microsoft.com/microsoft-edge/webview2/>.

### Optional

- **Wails CLI** — not required (the Makefile builds with plain `go build`), but
  handy to check the environment:
  ```sh
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  wails doctor
  ```
- **go-winres** — only to regenerate the embedded Windows icon resource
  (`make winres`):
  ```sh
  go install github.com/tc-hib/go-winres@latest
  ```

---

## How to run

```sh
make build   # builds the desktop UI (Wails) + the Alpaca server into one binary
./mount_panel  # or `make run`
```

Wails needs the `production` build tag and, on macOS, the
`UniformTypeIdentifiers` framework — the `Makefile` wraps the exact command:

```sh
CGO_LDFLAGS="-framework UniformTypeIdentifiers" go build -tags production -o mount_panel .
```

On Windows, build a GUI executable without a console window (`make build-windows`):

```sh
go build -ldflags "-H=windowsgui" -tags production -o mount_panel.exe .
```

The Windows `.exe`, window and taskbar icon comes from `rsrc_windows_amd64.syso`
(a Windows resource embedded by `go build`). It is configured in
[`winres/winres.json`](winres/winres.json), whose source image is
[`winres/icon.png`](winres/icon.png) — a square transparent-background PNG, set at
resource ID 3 (the ID Wails loads for the title bar). To change the icon, replace
that PNG (keep it square so it is not distorted) and regenerate the resource with
`make winres` (requires `go install github.com/tc-hib/go-winres@latest`).

The app **starts even when the mount is not connected**: it detects the serial
port (enumerates the ports + `capabilities` handshake to verify it is a mount) and
reconnects automatically once it is plugged in. Works on macOS and Windows.

Configuration via environment variables (see `internal/tools/env`):

- `PORT` — Alpaca server port (default `11111`).
- `SERIAL_PORT` — pins the CDC port (empty = auto-detection).

For development with hot-reload you can use `wails dev` (requires installing the
CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`).
