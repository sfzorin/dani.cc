# Morning — 30-day morning workout (PWA)

A mobile-first web app for a daily morning routine. 30-day cycle, warm-up →
main set (2 rounds) → cool-down. Floor + mat only, no jumps, no equipment.
Per-user progress and streaks stored on the server, adjustable rest, voice cues,
animated exercise figures, and a 7-language UI.

Built with the [doors](https://github.com/doors-dev/doors) Go framework (GoX
templates), SQLite (pure-Go `modernc.org/sqlite`), and a client-side workout
player.

## Features

- **30-day program** with progression (exercise count + reps/time grow per day).
- **Two exercise types**: reps (tap *Done*) and timed (countdown).
- **Workout player** (client-side): SVG animation, progress bar, pause/skip,
  rest periods, applause + confetti + random encouragement on finish.
- **Voice cues** via the browser Web Speech API (exercise without looking at the
  screen). Toggle in settings. (Designed to later swap to recorded neural-TTS.)
- **Streaks** — consecutive calendar days; the goal is no skipped days.
- **Auth** — email + password (bcrypt). Built to add Google/Apple/Supabase later.
- **Settings** — rest time 10–40 s, language, voice on/off, reset progress.
- **7 languages**: Russian, English, Turkish, German, Spanish, French, Italian.
- **Installable PWA** (manifest + service worker + icons).

## Run

Prerequisites are bundled/handled by `dev.sh`:
- Go ≥ 1.25 (`/opt/homebrew/bin/go`)
- The `gox` CLI binary lives in `.tooling/bin/` (pinned to the module version).

```bash
./dev.sh run          # gox gen + go run .  → http://localhost:8080
```

Other commands: `./dev.sh gen` (regenerate `.x.go` from `.gox`), `./dev.sh build`,
`./dev.sh tidy`. The database file defaults to `morning.db` (override with
`MORNING_DB`).

## Regenerating assets

```bash
go run ./cmd/gensvg     # exercise SVGs  → assets/static/ex/*.svg
go run ./cmd/genicons   # PWA PNG icons  → assets/static/icon-*.png
```

## Project layout

```
main.go               app factory: DB, session bootstrap, static middleware
app.gox               HTML shell, PWA <head>, service-worker registration, router
path/                 URL path model
internal/
  content/            exercise pool + 30-day program builder
  i18n/               7-language UI strings, exercise names, voice cues
  db/                 SQLite: users, progress (streaks), sessions
  auth/               bcrypt, session tokens, shared Session value
  app/                process-wide deps (DB, today's date)
segments/pages/       login/register, home (today + calendar), workout, settings
components/           shared UI (404)
assets/
  style.css           mobile-first dark theme
  player.js           client-side workout engine (timer, TTS, rest, finish)
  static/             manifest, sw.js, icons, ex/<id>.svg
cmd/gensvg, cmd/genicons   asset generators
```

## Notes / next steps

- **Apple sign-in** needs a paid Apple Developer account + verified HTTPS domain;
  deferred. Google / Supabase OAuth can be added via the doors OAuth-callback
  pattern (`internal/auth` is framework-agnostic on purpose).
- **Voice**: browser TTS now; for a recorded "human" voice, pre-generate per-phrase
  MP3s per language and have `player.js` play those instead of `speechSynthesis`.
- Exercise SVGs are schematic animated figures grouped by motion archetype; the
  per-exercise mapping in `cmd/gensvg` can be refined for more distinct poses.
```
