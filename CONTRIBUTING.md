# Contributing to morning30.com

Thanks for your interest in improving the project! This guide covers the local
workflow and a few conventions.

## Prerequisites

- Go ≥ 1.25
- The pinned `gox` CLI binary in `.tooling/bin/` (it ships with the repo; it
  cannot be `go install`ed because its version must match the
  `github.com/doors-dev/gox` module in `go.mod`).

## Workflow

```bash
./dev.sh run      # regenerate templates + run on http://localhost:8080
./dev.sh build    # regenerate + go build ./...
./dev.sh gen      # regenerate .x.go from .gox only
./dev.sh tidy     # go mod tidy
go test ./...     # run the test suite
```

### GoX templates

UI is written in `.gox` files. The `.x.go` files beside them are **generated** by
`gox gen` — never edit a `.x.go` by hand. After changing any `.gox`, run
`./dev.sh gen` (or `build`/`run`, which do it for you) and commit both the `.gox`
and its regenerated `.x.go`.

## Conventions

- Match the style of the surrounding code (naming, comment density, error
  handling).
- Keep user-facing strings in `internal/i18n` — every key is a `[7]string` in the
  order `[ru, en, tr, de, es, fr, it]`. Add all seven languages.
- The workout player (`assets/player.js`) is plain browser JavaScript; run
  `node --check assets/player.js` before committing changes to it.
- Run `go build ./...` and `go test ./...` before opening a pull request.

## Licensing

By contributing, you agree that your contributions are licensed under the
project's [Apache License 2.0](LICENSE).
