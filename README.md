# morning30.com

A mobile-first PWA for a **30-day morning-exercise** habit: a daily floor-and-mat
routine you can run with the screen in your pocket, guided entirely by voice.

- **Repository:** <https://github.com/sfzorin/morning30>
- **License:** [Apache-2.0](LICENSE)
- **Status:** active development

Each day is a single session — **warm-up → main block (×2 rounds) → cool-down** —
that grows from ~20 to ~30 exercises over the cycle, with lighter recovery days
mixed in. Progress, streaks and preferences are stored per user on the server.

Built with the [doors](https://github.com/doors-dev/doors) Go framework and its
GoX templates (server-driven reactive UI), a pure-Go SQLite database, and a
small client-side workout engine.

---

## Features

- **Named workout programs.** Ships with more than one built-in program (e.g.
  *Sergey* — floor/mat, no jumps; *Vlad / Level 2* — plyometric + cardio). Each
  is a 30-day plan with baked-in day-to-day progression. Adding another program
  is a single registry entry.
- **Universal difficulty.** A per-user level from **−3 to +3** scales every
  program by ±10% per step (reps and seconds only; breaths and warm-up/cool-down
  stay fixed).
- **Custom programs & exercises.** A full in-app editor (reorder, add, remove,
  retime) plus per-exercise JSON import/export. A user's custom exercises overlay
  the built-in library.
- **Two exercise types.** Reps (tap *Done*) and timed (countdown). Counts are
  shown as-is; the athlete does both sides of any one-sided movement themselves.
- **Voice-first player (client-side).** Progress bar, pause / skip / back, rest
  periods, screen wake-lock, and a celebratory finish (applause, confetti,
  localized encouragement). The player runs entirely in the browser from a
  server-provided JSON payload and reports completion back through a hook.
- **Spoken guidance** via the browser Web Speech API: a start cue, a rest
  announcement of the next exercise with its count, a 3-2-1 lead-in, a half-time
  marker on timed sets, and running technique narration (how to do it / common
  mistakes). Four levels in settings: off · minimal · normal · detailed.
- **Male / female figures.** Each exercise has a static illustration; a settings
  toggle picks the male or female variant.
- **Accounts.** Email + password (bcrypt). A **guest mode** lets people start
  immediately; signing up later upgrades the guest in place, keeping all progress.
- **Streaks & calendar.** Consecutive-day streaks and a monthly activity calendar.
- **7 languages.** Russian, English, Turkish, German, Spanish, French, Italian —
  UI strings, exercise names, and voice cues.
- **Installable PWA** (web manifest + service worker + icons).

---

## Tech stack

| Layer        | Choice |
|--------------|--------|
| Language     | Go (≥ 1.25) |
| UI framework | [doors](https://github.com/doors-dev/doors) + GoX templates (`.gox` → generated `.x.go`) |
| Database     | SQLite via pure-Go [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite) (no CGO) |
| Auth         | bcrypt (`golang.org/x/crypto`), cookie sessions |
| Workout player | Vanilla JavaScript + Web Speech API |

---

## Architecture

The app is **server-driven** for everything except the workout itself:

- **Pages** (`segments/pages/`) — sign-in/up, home (today + calendar + streak),
  settings, and the router — are doors components rendered and updated on the
  server. A shared reactive `Session` value drives auth-gated routing.
- **The workout player** (`assets/player.js`) is the one client-side piece. The
  server builds a self-contained JSON payload (the day's items, localized labels
  and voice cues, the user's voice mode and figure variant) and hands it to the
  player, which owns all timing, speech, pausing and rest. On finish it calls a
  server hook to record the day and advance the cycle.
- **Layered exercise catalog** (`internal/catalog/`) resolves each exercise
  through three layers, the first match winning:
  `per-user custom → global DB library → built-in code`.
- **Content** (`internal/content/`) holds the exercise library and the named
  standard programs; a program is a shared warm-up/cool-down plus a per-day main
  block, flattened and difficulty-scaled at run time.
- **i18n** (`internal/i18n/`) holds all 7-language UI strings, exercise names and
  spoken cues.

---

## Project layout

```
main.go              entry point: open DB, bootstrap session, routing, static files
app.gox              HTML shell, PWA <head>, service-worker registration, root router
path/                URL path model (/, /login, /register, /day/:n, /settings)
internal/
  content/           exercise library + named 30-day programs + difficulty scaling
  catalog/           per-user exercise catalog (custom over DB over code)
  i18n/              7-language UI strings, exercise names, voice cues
  db/                SQLite: users, sessions, activity (streaks), exercise library
  auth/              bcrypt, session tokens, the shared Session value
  app/               process-wide deps (DB) + per-user program/workout resolution
segments/pages/      router, auth, home, workout, settings, shared header
components/           shared UI (e.g. 404)
assets/
  style.css          mobile-first dark theme
  player.js          client-side workout engine (timer, TTS, rest, finish)
  exupload.js        per-exercise JSON import helper
  static/            manifest, service worker, icons, exercise figures
cmd/genicons         PWA icon generator
deploy/              production docker-compose
Dockerfile           multi-stage container build
```

> **Note:** exercise figures live under `assets/static/ex/<m|f>/<id>.png` and are
> produced by a separate figure-rendering pipeline that is being extracted into
> its own repository; it is intentionally not documented here.

---

## Getting started

**Prerequisites**

- Go ≥ 1.25
- The `gox` CLI — a prebuilt binary pinned to the `github.com/doors-dev/gox`
  module version, kept in `.tooling/bin/` (it cannot be `go install`ed).

**Run**

```bash
./dev.sh run          # regenerate templates, then `go run .`  →  http://localhost:8080
```

`dev.sh` sets `PATH` and `GOFLAGS=-mod=mod` for you. Other commands:

```bash
./dev.sh gen          # regenerate .x.go from .gox (run after editing any .gox)
./dev.sh build        # gen + go build ./...
./dev.sh tidy         # go mod tidy
```

> Edit **`.gox`** files only — the `.x.go` files next to them are generated and
> overwritten by `gox gen`.

**Configuration**

| Env          | Default      | Description |
|--------------|--------------|-------------|
| `MORNING_DB` | `morning.db` | SQLite database file path (the container points this at the persistent `/data` volume) |

The server listens on `:8080`.

---

## Building & deployment

A multi-stage [`Dockerfile`](Dockerfile) produces a small static binary on a slim
Debian runtime; [`deploy/docker-compose.yml`](deploy/docker-compose.yml) runs it
with a persistent data volume. Production deploys are automated via GitHub
Actions (build → ship → restart); the deployment **secrets, volume and DNS are
managed outside this repository** by the operator.

---

## Internationalization

All user-facing text is centralized in `internal/i18n` as `[7]string` tables in
the order `[ru, en, tr, de, es, fr, it]`. Exercise technique content lives in
`internal/content` per language. Adding a language means adding a column to those
tables and a `Lang` entry — there are no translation files to wire up.

---

## Contributing

Issues and pull requests are welcome — see [CONTRIBUTING.md](CONTRIBUTING.md) for
the build/regenerate workflow and conventions.

## License

Licensed under the **Apache License, Version 2.0**. See [LICENSE](LICENSE) and
[NOTICE](NOTICE).
