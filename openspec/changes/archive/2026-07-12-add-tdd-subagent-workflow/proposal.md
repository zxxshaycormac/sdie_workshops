## Why

The engineering harness already requires OpenSpec, control-point tracing, and
proportionate verification, but it does not make the development sequence
explicit. Future behavior changes can still drift into "implement first, test
after" work, which weakens contract discipline and makes tests confirm an
already-chosen implementation instead of the intended behavior.

## What Changes

- Add a TDD-first workflow for behavior and contract changes.
- Require an independent subagent to write or review focused tests before
  production code when that capability is available.
- Require contributors to record expected failing test evidence before
  implementation when practical, then passing evidence after implementation.
- Document bounded exceptions for documentation-only, formatting-only, and
  mechanical changes.

### Non-goals

- Do not mandate broad full-suite runs for every change.
- Do not require subagents for pure documentation or mechanical edits.
- Do not change product behavior, test code, generated assets, or runtime
  configuration in this change.

## Capabilities

### Modified Capabilities

- `engineering-harness`: Document TDD-first development and subagent test
  ownership as part of the repository operating model.
