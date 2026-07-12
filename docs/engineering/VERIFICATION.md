# Verification Guide

Verification is selected from the behavior and risk in the OpenSpec design,
not from how many files changed. Always report the exact command and scope.

## Fast Entry Point

`tools/harness/verify.sh` provides explicit, focused modes:

```sh
./tools/harness/verify.sh harness
./tools/harness/verify.sh spec
./tools/harness/verify.sh go ./services/repository/...
./tools/harness/verify.sh frontend web_src/js/features/repo-release.test.js
```

- `harness` checks required harness files, shell syntax, and OpenSpec health.
- `spec` runs strict validation for all changes and specs plus `doctor`.
- `go` runs Go tests with the repository's default SQLite tags for only the
  packages supplied after the mode.
- `frontend` runs Vitest once for only the paths or filters supplied.

Focused modes reject an empty scope. This prevents a vague command from being
reported as evidence for an unrelated surface.

## Change-To-Check Matrix

For behavior and contract changes, select the focused test before production
code changes, run it once for expected failure when practical, and rerun the
same scope after implementation. See `docs/engineering/TDD_WORKFLOW.md`.

| Change surface | Minimum focused evidence | Broaden when |
| --- | --- | --- |
| OpenSpec or harness docs | `verify.sh harness` | Shared rules or tooling changed |
| Go package behavior | `verify.sh go ./exact/package` | Shared model, context, auth, or infrastructure changed |
| API handler or contract | Handler/service package tests and strict spec validation | Public fields, auth, pagination, or shared converter changed |
| Model or query | Model package tests with SQLite tags | SQL differs by engine or schema/migration changed |
| Migration | Target migration tests | Production supports multiple database engines |
| JavaScript or Vue | `verify.sh frontend <test>` | Shared component or bootstrap changed |
| Template plus browser behavior | Focused Go/Vitest checks | User workflow requires rendered-page or browser evidence |
| CSS or layout | Focused lint plus visual inspection | Shared theme or responsive layout changed |
| Queue, cron, webhook, Actions | Producer and consumer package tests | Retry, timeout, final state, or integration boundary changed |
| Build or generated output | Relevant generator and source check | Generator is shared or release output changes |

## Canonical Project Commands

Broad commands remain visible rather than hidden behind the harness:

```sh
make test-backend       # all non-integration Go packages
make test-frontend      # Vitest watch/default mode
make lint-backend       # golangci-lint, custom vet, editorconfig
make lint-frontend      # ESLint and Stylelint
make lint-md            # Markdown lint
make build              # frontend and backend build
```

Some Make checks use `git diff` or `git status`. They require available `.git`
metadata to provide their intended clean-tree evidence.

## API And Swagger

For public API behavior, verify all of these surfaces:

- Route and auth middleware in `routers/api/v1/api.go`.
- Handler annotations and status behavior.
- Request and response shapes in `modules/structs/`.
- Converters and service/model behavior.
- Swagger references and generated output.

Relevant commands include:

```sh
make generate-swagger
make lint-swagger
make swagger-validate
```

`make swagger-check` compares output with Git and therefore needs available Git
metadata to work as designed.

## Database And Integration Checks

The default Go test tags are `sqlite sqlite_unlock_notify`. Focused package
tests use those tags through the harness.

Integration suites are database-specific:

```sh
make test-sqlite
make test-mysql
make test-pgsql
make test-mssql
```

Migration targets and `#TestName` variants are defined in the Makefile. These
checks can build large test binaries, generate local configuration, and require
database services. Select them in the OpenSpec design before running them.

## End-To-End Checks

Playwright workflows are under `tests/e2e/` and are launched through Go test
wrappers:

```sh
make test-e2e
make 'test-e2e-sqlite#TestSpecificName'
```

They may install browser dependencies, build a test server, and write artifacts
under `tests/e2e/`. Use them for complete user workflows, not as a first-line
replacement for focused unit or package tests.

## Evidence Standard

A completion report must distinguish:

```text
Test-first:
- Added tests/api_tags_search_test.go before product code.
- Pre-implementation run failed as expected: missing tag query filtering.

Passed:
- ./tools/harness/verify.sh spec
- ./tools/harness/verify.sh go ./services/repository/...

Not run:
- Full backend suite; change was isolated to one package.
- E2E suite; no rendered workflow changed.

Limitation:
- Git-based diff and clean-tree verification unavailable because .git is absent.
```

Do not say "all tests pass" unless every relevant suite actually ran. A failed
unrelated check should be reported with enough evidence to distinguish a new
regression from an environment or pre-existing failure.
