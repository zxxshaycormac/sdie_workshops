# tools Roadmap

## Purpose

Stores repository-local operational tools that support development, verification,
and the engineering harness.

## Route Map

- `harness/verify.sh`: focused verification wrapper used by AGENTS and
  OpenSpec workflows.
- Other tools should be small, documented, and callable from the repository
  root.

## When To Edit

Edit when changing local verification entry points or adding repeatable
repository tooling. Prefer scripts that wrap existing Make or CLI targets rather
than inventing parallel behavior.

## Verification

Run the tool's help path and at least one successful focused mode. For harness
changes, run `./tools/harness/verify.sh spec`.

## Boundaries

Do not store task-specific scratch scripts here unless they are promoted into
durable project tooling.
