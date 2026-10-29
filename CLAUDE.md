# mount_panel Development Rules

Keep this file in a simple format so it can be read and edited quickly.

## Sibling projects (firmware)

This gateway talks to four ESP32-S3 devices over USB CDC-ACM. Each firmware lives
in its own local repo:

- **`~/Dev/bft_firmware`** — mount firmware (the BFT fork mount). Serial protocol
  spec in `main/proto/README.md`, rules in `CLAUDE.md`.
- **`~/Dev/eaf_firmware`** — focuser firmware (`device_id="eaf"`). Serial protocol
  spec in `main/proto/README.md`.
- **`~/Dev/fan_firmware`** — primary-mirror fan + mirror temperature firmware
  (`device_id="fan"`).
- **`~/Dev/temp_firmware`** — ambient temperature/humidity + vibration firmware
  (`device_id="temp"`).

The three accessory firmwares speak `proto=2`, `device_type="accessories"`, and
are told apart by their `device_id`. When firmware references are needed — the
exact message formats or what each device exposes — look there.

## Project definition

- Desktop gateway and UI for the BFT mount (the `bft_firmware` firmware).
- Replaces the HTTP/Alpaca server that currently runs inside the ESP32-S3.
- Talks to the mount and the three accessory controllers (focuser, fan,
  temperature/vibration) over **USB CDC-ACM serial protocols** (length + JSON)
  and exposes **ASCOM Alpaca** (N.I.N.A., PHD2) + a **desktop UI (Wails)**.
- The protocol specs live in the firmware: `bft_firmware/main/proto/README.md`,
  `eaf_firmware/main/proto/README.md`, etc.

## Architecture

Layers, from the outside in:

```
N.I.N.A. / PHD2 ──HTTP──▶ Alpaca REST (GIN) ──┐
                                               ├─▶ Mount ──▶ MountProto ──┐
Desktop UI (Wails) ──bindings──▶ ui ───────────┘                          ├─▶ Proto (serial CDC)
                            ui ──▶ EAF  ──▶ EAFProto ─────────────────────┘
                            ui ──▶ Fan  ──▶ EAFProto ─────────────────────┘
                            ui ──▶ Temp ──▶ EAFProto ─────────────────────┘
```

- **proto/** — shared serial transport: framing (u32 LE length + JSON), port
  opening and a generic request/response correlator by `id`. Knows nothing about
  the mount or accessory message contracts.
- **mountproto/** — the mount message contract: the 5 families (`capabilities`,
  `state`, `config`, `control`, `action`) as typed requests, plus mount port
  detection (`capabilities` handshake).
- **eafproto/** — the accessory message contract, shared by the three accessory
  devices: `proto=2`, `device_type="accessories"`, a `target`/`action`/`op`/
  `params` envelope with `session`/`boot_id`. Port detection matches the
  `device_id` (`eaf`/`fan`/`temp`).
- **mount/** — mount model: manages the connection (detection, **auto-reconnect**
  and hot-unplug detection) and translates high-level operations (slew, move
  axis, guide, track, park…) into `mountproto` calls. No state cache.
- **eaf/** — focuser model: manages its serial link and exposes focuser
  operations to the UI. No state cache.
- **fan/** — fan model: manages its serial link and exposes fan and
  mirror-temperature operations to the UI. No state cache.
- **temp/** — temperature/vibration model: manages its serial link and exposes
  ambient/humidity/dew/vibration readings to the UI. No state cache.
- **alpaca_mount/** — Alpaca REST server for the telescope: ITelescopeV4 endpoints
  (properties + methods) plus `management` and UDP discovery. Translates
  Alpaca → `mount` → `mountproto`. `Start()` builds the router and registers
  both the telescope and the accessory devices.
  - **alpaca_mount/server/** — shared GIN framework (engine singleton, response
    envelope, param parsing, error numbers, server identity).
- **alpaca_eaf/** — Alpaca REST server for the accessory devices: IFocuserV4
  (→ `eaf`), ObservingConditions V2 (→ `temp`) and Switch V3 (→ `fan` + `temp`).
- **ui/** — Wails binding (`App`): exposes `mount`, `eaf`, `fan` and `temp` to the
  desktop UI.
- **frontend/** — web assets (HTML/CSS/JS Alpine.js) served by Wails, which call
  `window.go.ui.App.*`.
- **tools/** — cross-cutting: `logx` (logger) and `env` (config).

## General rules

- **Functional style.** Package-level functions, no methods on structs unless
  there is real state to encapsulate. One file per use case.
- **No dependency injection.** Packages are imported and used directly.
- **No interface overuse.** Don't create interfaces just for testing or
  "best practices". Only when they implement the **strategy** pattern for real.
- **No swagger, no coverage.** There are Alpaca integration tests in
  `test/alpaca` (build tag `integration`), separate from the functional code.
- **GIN** for HTTP, following the pattern of the reference repo
  `nmarsollier/resourcesgo`: `engine` singleton in `server`, an
  `init<Endpoint>` that registers the route + a `get<X>`/`post<X>` handler per
  endpoint.
- **Never commit or push** without explicit authorization.
- Commit only when the feature is verified and compiling.
- The source code (`.go` files) is the single source of truth.

## Code conventions

- **Language**: Go (not C++/Rust).
- **Naming**: `PascalCase` for exported types/structs, `camelCase` for functions
  and variables, `snake_case` for files.
- **One file per use case** (e.g. `find_by_id.go`, `create.go`), as in
  `resourcesgo`.
- **Alpaca errors**: endpoints always respond HTTP 200; the error goes in the
  body's `ErrorNumber` (`1024` = not implemented, `1025` = invalid value). HTTP
  status codes are not used for business errors.
- **Config**: `env.Get()` singleton with values from environment variables and
  defaults.
- **Logging**: `tools/logx` with context and `correlation_id`, no `fmt.Println`.
- **JSON**: use `encoding/json` directly, no wrappers in `tools`.
- Comments inside `.go` files are written in **English**.
- Documentation (`.md` files) is written in **English**.
- Respect: DRY - YAGNI - KISS.
- Header documentation describes the business problem or use case the function
- Only comment code when it is hard for a human to understand.
- Before any change, read the current state of the file to modify.
- Never commit or push to GitHub without explicit authorization.
3$

## Package relationships (dependencies)

- `alpaca_mount` → `mount` → `mountproto` → `proto`; `alpaca_mount` →
  `alpaca_eaf`; `alpaca_eaf` → `eaf` + `fan` + `temp` → `eafproto` → `proto`
  (no reverse dependencies). Both `alpaca_mount` and `alpaca_eaf` use
  `alpaca_mount/server` for the shared HTTP framework.
- `ui` → `mount`, `ui` → `eaf`, `ui` → `fan`, `ui` → `temp`.
- `proto` is self-contained: serial port, framing, request/response with `id`,
  and the disconnection signal.
- `tools/*` depends on no other package; the others may use `tools`.
