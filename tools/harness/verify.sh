#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

usage() {
  cat <<'EOF'
Usage:
  ./tools/harness/verify.sh harness
  ./tools/harness/verify.sh spec
  ./tools/harness/verify.sh go <package> [package...]
  ./tools/harness/verify.sh frontend <test-path-or-filter> [more...]

Modes:
  harness   Check harness files, shell syntax, and OpenSpec health.
  spec      Strictly validate all OpenSpec changes/specs and run doctor.
  go        Run focused Go tests with the repository's default SQLite tags.
  frontend  Run focused Vitest tests once.
EOF
}

say_run() {
  printf '+ '
  printf '%q ' "$@"
  printf '\n'
  "$@"
}

verify_spec() {
  command -v openspec >/dev/null 2>&1 || {
    echo "error: openspec is not installed or not on PATH" >&2
    exit 1
  }
  say_run openspec validate --all --strict --no-interactive
  say_run openspec doctor
}

verify_harness() {
  local required=(
    AGENTS.md
    docs/engineering/PROJECT_MAP.md
    docs/engineering/OPEN_SPEC.md
    docs/engineering/VERIFICATION.md
    openspec/config.yaml
    tools/harness/verify.sh
  )
  local path

  for path in "${required[@]}"; do
    if [[ ! -s "$path" ]]; then
      echo "error: required harness file is missing or empty: $path" >&2
      exit 1
    fi
  done

  say_run bash -n tools/harness/verify.sh
  verify_spec
}

mode="${1:-}"
if [[ -z "$mode" ]]; then
  usage >&2
  exit 2
fi
shift

case "$mode" in
  harness)
    if (( $# != 0 )); then
      echo "error: harness mode does not accept targets" >&2
      usage >&2
      exit 2
    fi
    verify_harness
    ;;
  spec)
    if (( $# != 0 )); then
      echo "error: spec mode does not accept targets" >&2
      usage >&2
      exit 2
    fi
    verify_spec
    ;;
  go)
    if (( $# == 0 )); then
      echo "error: go mode requires at least one package" >&2
      usage >&2
      exit 2
    fi
    command -v go >/dev/null 2>&1 || {
      echo "error: go is not installed or not on PATH" >&2
      exit 1
    }
    say_run go test -tags="sqlite sqlite_unlock_notify" "$@"
    ;;
  frontend)
    if (( $# == 0 )); then
      echo "error: frontend mode requires at least one test path or filter" >&2
      usage >&2
      exit 2
    fi
    command -v npx >/dev/null 2>&1 || {
      echo "error: npx is not installed or not on PATH" >&2
      exit 1
    }
    say_run npx vitest run "$@"
    ;;
  -h|--help|help)
    usage
    ;;
  *)
    echo "error: unknown mode: $mode" >&2
    usage >&2
    exit 2
    ;;
esac
