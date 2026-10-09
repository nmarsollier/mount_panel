# nrm_ha_go Development Rules

Keep this file in a simple format so it can be read and edited quickly.

## Sibling project (firmware)

This gateway talks to the mount's firmware over USB CDC-ACM. The firmware lives
in **another local repo**:

- **`~/Dev/nrm_ha`** — ESP32-S3 firmware (the NRM-HA project). The serial
  protocol spec is in `main/proto/README.md`, and its rules in `CLAUDE.md`.

When firmware references are needed — the exact message formats or what the mount
exposes — look there.

## Project definition

- Desktop gateway and UI for the NRM-HA mount (the `nrm_ha` firmware).
- Replaces the HTTP/Alpaca server that currently runs inside the ESP32-S3.
- Talks to the ESP over a **USB CDC-ACM serial protocol** (length + JSON) and
  exposes **ASCOM Alpaca** (N.I.N.A., PHD2) + a **desktop UI (Wails)**.
- The protocol spec lives in the firmware: `nrm_ha/main/proto/README.md`.

## Architecture

Layers, from the outside in:

```
N.I.N.A. / PHD2 ──HTTP──▶ Alpaca REST (GIN) ──┐
                                               ├─▶ Mount ──▶ Proto (serial CDC)
Desktop UI (Wails) ──bindings──▶ ui ───────────┘
```

- **proto/** — serial protocol client: framing (u32 LE length + JSON), the 5
  message families (`capabilities`, `state`, `config`, `control`, `action`), and
  port detection (enumerates + `capabilities` handshake to verify it is a mount).
  This is the layer that talks over the serial port.
- **mount/** — mount model: manages the connection (detection, **auto-reconnect**
  and hot-unplug detection) and translates high-level operations (slew, move
  axis, guide, track, park…) into `proto` calls. No state cache: every read
  fetches a fresh snapshot over the serial link.
- **alpaca/** — Alpaca REST server with GIN: ITelescopeV4 endpoints (properties +
  methods) and `management`. Translates Alpaca → `mount` → `proto`.
  - **alpaca/server/** — GIN engine (singleton) + logging middleware.
- **ui/** — Wails binding (`App`): exposes `mount` to the desktop frontend.
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

## Package relationships (dependencies)

- `alpaca` → `mount` → `proto`, and `ui` → `mount` (no reverse dependencies).
- `proto` is self-contained: it handles the serial port, framing, the
  request/response with `id`, port detection and the disconnection signal.
- `tools/*` depends on no other package; the others may use `tools`.
