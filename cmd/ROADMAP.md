# cmd Roadmap

## Purpose

Owns process entry commands for the Gitea binary.

## Route Map

- `main.go`, `cmd.go`: command registration and global setup.
- `web.go`, `web_*.go`: installed server startup, listeners, TLS, ACME, and
  graceful behavior.
- `admin_*.go`: administrative CLI commands.
- `hook.go`, `serv.go`, `keys.go`: Git, SSH, and hook entry points.
- `dump*`, `restore_repo.go`, `migrate*`, `doctor*`: maintenance commands.

## When To Edit

Edit when changing CLI behavior, startup order, maintenance commands, or Git/SSH
entry points. Trace downstream services and initialization side effects before
changing command wiring.

## Verification

Run focused `go test ./cmd` for command logic and the specific CLI command with
safe inputs. For server startup changes, also run a local `gitea web` smoke test.

## Boundaries

Business orchestration belongs in `services/`; persistence belongs in `models/`;
`cmd/` should wire commands and report errors.
