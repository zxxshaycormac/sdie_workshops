## Why

This large, cross-layer Gitea codebase has no repository-level operating model for AI-assisted changes. Without explicit context, boundaries, and executable quality gates, an agent can produce plausible edits while missing the real control point, generated artifacts, compatibility contracts, or required tests.

## What Changes

- Add a root instruction surface that defines the required change workflow, safety boundaries, and definition of done.
- Add a code-grounded project map for locating entry points, contracts, persistence, frontend sources, and tests.
- Add a verification guide and executable wrapper for consistent, scoped checks.
- Configure OpenSpec with project context and artifact rules.
- Document the OpenSpec lifecycle so future changes move from intent to requirements, design, tasks, implementation, validation, and archive.

### Non-goals

- Refactor or change Gitea product behavior.
- Replace the existing Makefile, linters, or test frameworks.
- Treat generated files or local runtime state as source code.
- Claim exhaustive knowledge of every package in the repository.

## Capabilities

### New Capabilities

- `engineering-harness`: Repository-level guidance, project navigation, change planning, and proportionate verification for human and AI contributors.

### Modified Capabilities

None.

## Impact

The change adds engineering metadata and tooling at `AGENTS.md`, `docs/engineering/`, `tools/harness/`, and `openspec/`. It does not alter runtime APIs, database schemas, dependencies, generated assets, or Gitea behavior. The primary compatibility surface is contributor workflow: future non-trivial changes are expected to use OpenSpec and report exactly which checks ran.

Observed facts include the Go 1.22 module, the `main.go -> cmd -> routers` startup path, the layered Go directories, and existing Make/Vitest/Playwright checks. A current limitation is that this workspace has no `.git` metadata, so history-aware automation cannot yet be relied upon.
