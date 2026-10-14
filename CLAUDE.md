# mount_panel Development Rules

Keep this file in a simple format so it can be read and edited quickly.

## Sibling project (firmware)

This gateway talks to two ESP32-S3 devices over USB CDC-ACM. Each firmware lives
in its own local repo:

- **`~/Dev/nrm_ha`** — mount firmware (the NRM-HA project). Serial protocol spec
  in `main/proto/README.md`, rules in `CLAUDE.md`.
- **`~/Dev/nrm_bft_eaf`** — accessory firmware (focuser/fan/environment, the
  "NRM-EAF" project). Serial protocol spec in `main/proto/README.md`. It speaks
  `proto=2`, `role="accessories"`; much of it is still mock (fan, sensors).

When firmware references are needed — the exact message formats or what each
device exposes — look there.

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
                                               ├─▶ Mount ──▶ MountProto ──┐
Desktop UI (Wails) ──bindings──▶ ui ───────────┘                          ├─▶ Proto (serial CDC)
                            ui ──▶ EAF ──▶ EAFProto ─────────────────────┘
```

- **proto/** — shared serial transport: framing (u32 LE length + JSON), port
  opening and a generic request/response correlator by `id`. Knows nothing about
  the mount or accessory message contracts.
- **mountproto/** — the mount message contract: the 5 families (`capabilities`,
  `state`, `config`, `control`, `action`) as typed requests, plus mount port
  detection (`capabilities` handshake).
- **eafproto/** — the accessory (EAF) message contract: `proto=2`,
  `role="accessories"`, `target`/`action`/`op`/`params` with a `session`/`boot_id`
  envelope. Focuser move/halt, fan mode and state reads.
- **mount/** — mount model: manages the connection (detection, **auto-reconnect**
  and hot-unplug detection) and translates high-level operations (slew, move
  axis, guide, track, park…) into `mountproto` calls. No state cache.
- **eaf/** — accessory model: manages the second serial link and exposes focuser
  and fan operations to the UI. No state cache.
- **alpaca_mount/** — Alpaca REST server for the telescope: ITelescopeV4 endpoints
  (properties + methods) plus `management` and UDP discovery. Translates
  Alpaca → `mount` → `mountproto`. `Start()` builds the router and registers
  both the telescope and the focuser.
  - **alpaca_mount/server/** — shared GIN framework (engine singleton, response
    envelope, param parsing, error numbers, server identity).
- **alpaca_eaf/** — Alpaca REST server for the focuser: IFocuserV4 endpoints,
  translating Alpaca → `eaf` → `eafproto`.
- **ui/** — Wails binding (`App`): exposes `mount` and `eaf` to the desktop UI.
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
  `alpaca_eaf` and `alpaca_eaf` → `eaf` → `eafproto` → `proto` (no reverse
  dependencies). Both `alpaca_mount` and `alpaca_eaf` use
  `alpaca_mount/server` for the shared HTTP framework.
- `ui` → `mount` and `ui` → `eaf`.
- `proto` is self-contained: serial port, framing, request/response with `id`,
  and the disconnection signal.
- `tools/*` depends on no other package; the others may use `tools`.
