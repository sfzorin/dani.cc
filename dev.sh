#!/usr/bin/env bash
# Dev helper for the "morning" doors app.
# Usage: ./dev.sh gen | build | run | tidy | all
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export PATH="$ROOT/.tooling/bin:/opt/homebrew/bin:$PATH"
export GOFLAGS=-mod=mod
cd "$ROOT"

gen() {
  # Regenerate all .x.go from .gox (gox gen recurses from the given dir).
  gox gen .
}

cmd="${1:-all}"
case "$cmd" in
  gen)   gen ;;
  tidy)  go mod tidy ;;
  build) gen && go build ./... ;;
  run)   gen && go run . ;;
  all)   gen && go build ./... && echo "BUILD OK" ;;
  *) echo "unknown command: $cmd"; exit 1 ;;
esac
