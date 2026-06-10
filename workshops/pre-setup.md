# Pre-Setup Guide

Before the workshop begins, ensure your machine meets all requirements below.
Steps are provided for macOS, Windows (WSL2), and Linux.

## 1. Prerequisites

| Tool | Minimum Version | Purpose |
|---|---|---|
| Git | 2.x | Source control |
| Go | 1.22+ | Backend compilation |
| Node.js | 18+ | Frontend build |
| Make | any | Build orchestration |

### Install by Platform

**macOS (Homebrew):**

```bash
xcode-select --install    # Git + Make
brew install go node
```

**Windows (WSL2 — Ubuntu):**

Open a WSL2 terminal (Ubuntu recommended), then:

```bash
sudo apt update
sudo apt install -y git make gcc
```

Install Go from https://go.dev/dl/ — download the `linux-amd64` tarball:

```bash
curl -LO https://go.dev/dl/go1.22.10.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.10.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

Install Node.js via NodeSource:

```bash
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs
```

> **WSL2 tip:** If `make build` fails with CGO errors, install the C compiler:
> `sudo apt install -y gcc`

**Linux (Debian/Ubuntu):**

Same as WSL2 steps above (skip the WSL2 setup itself).

**Linux (Fedora/RHEL):**

```bash
sudo dnf install -y git make gcc
```

Go and Node.js — follow the same tarball/NodeSource instructions above, or use:

```bash
sudo dnf install -y golang nodejs
```

### Verify (all platforms)

```bash
go version      # expect: go1.22.x or higher
node --version  # expect: v18.x.x or higher
git --version
make --version
```

## 2. Clone & Build

```bash
# Clone the repo (if you haven't already)
git clone $WORKSHOP
cd gitea

# Build with SQLite support
TAGS="bindata sqlite sqlite_unlock_notify" make build
```

First build takes 2-5 minutes (frontend + backend). Subsequent builds are faster.

Verify the binary:

```bash
./gitea --version
# expect: Gitea version 1.22.x built with ... sqlite ...
```

## 3. Run the Dev Server

```bash
GITEA_RUN_MODE=dev ./gitea web
```

Open **http://localhost:3000** in your browser.

- First visit shows the **install page** — keep defaults (SQLite3), set an admin username/password, click "Install Gitea".
- This creates `custom/conf/app.ini` and the SQLite database.
- Stop the server with `Ctrl+C`.

### macOS / Linux

The command above works as-is in any terminal.

### Windows (WSL2)

Run the same command inside your WSL2 terminal. Then access the server from your Windows browser at **http://localhost:3000** — WSL2 forwards localhost ports automatically.

If localhost forwarding is disabled, find the WSL2 IP:

```bash
hostname -I | awk '{print $1}'
# open http://<that-ip>:3000 in your Windows browser
```

## 4. Verify Tests

```bash
# Backend unit tests (should pass)
make test-backend

# Single test (example)
go test -run TestBufWriter ./modules/git/pkg/

# SQLite integration tests
make test-sqlite
```

Integration tests take longer. At minimum, confirm `make test-backend` passes.

## 5. Install Claude Code

The workshops use Claude Code as the AI coding harness.

### macOS / Linux / WSL2

```bash
npm install -g @anthropic-ai/claude-code
```

### Windows (native)

Claude Code runs inside WSL2 — install it there, not in PowerShell.

Verify:

```bash
claude --version
```

## 6. Recommended Editor Setup

- **VS Code** with the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.Go) for jump-to-definition and refactoring
- **Claude Code VS Code extension** (optional) for in-editor AI assistance
- **WSL2 users:** open the repo folder via VS Code's "Remote - WSL" extension for full IDE support

## Troubleshooting

### `make build` fails with Go version error

Ensure Go >= 1.22:

```bash
go version
# macOS:  brew upgrade go
# Linux:  reinstall from go.dev/dl or upgrade via package manager
```

### `make build` fails on frontend (webpack)

Ensure Node >= 18:

```bash
node --version
# macOS:  brew upgrade node
# Linux:  sudo apt install -y nodejs  (via NodeSource)
```

### CGO / C compiler errors (Linux / WSL2)

SQLite requires CGO, which needs a C compiler:

```bash
# Debian/Ubuntu
sudo apt install -y gcc

# Fedora/RHEL
sudo dnf install -y gcc
```

### Port 3000 already in use

```bash
# macOS / Linux / WSL2
lsof -i :3000       # find the process
kill <PID>           # stop it
```

Or use a different port:

```bash
./gitea web -p 3001
```

### SQLite build tag missing

The binary must include `sqlite` in its build tags. Verify:

```bash
./gitea --version
# look for "sqlite" in the output
```

If missing, rebuild with:

```bash
TAGS="bindata sqlite sqlite_unlock_notify" make build
```

### WSL2: `make build` fails with "cannot find package"

Ensure Go and Node are installed **inside WSL2**, not in Windows. Verify:

```bash
which go    # should be /usr/local/go/bin/go or similar Linux path
which node  # should be /usr/bin/node or similar Linux path
```

If they point to a Windows path (e.g. `/mnt/c/...`), you're using the wrong shell.

## Quick Reference

```bash
# Build
TAGS="bindata sqlite sqlite_unlock_notify" make build

# Run
GITEA_RUN_MODE=dev ./gitea web

# Test
make test-backend
go test -run TestFunctionName ./path/to/package

# Lint
make lint

# Hot-reload (alternative to manual rebuild)
make watch
```
