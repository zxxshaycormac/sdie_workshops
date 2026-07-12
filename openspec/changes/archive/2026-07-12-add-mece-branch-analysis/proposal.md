## Why

The harness now requires TDD-first implementation, but it should also make the
test design thinking explicit. Tests written first are only useful when they
come from a disciplined model of the behavior space rather than from the most
obvious happy path.

This change adds MECE-style branch analysis to the engineering harness: split
possibilities into mutually exclusive buckets, make the set as exhaustive as the
current contract reasonably allows, and explicitly name excluded or unreachable
branches.

## What Changes

- Add MECE branch analysis guidance to the TDD workflow.
- Require OpenSpec designs and scenarios to consider relevant state, permission,
  input, data, downstream, timing, and UI/API output branches.
- Tie MECE output to subagent test design and final evidence.
- Keep root `AGENTS.md` concise by linking to detailed guidance.

### Non-goals

- Do not require speculative tests for impossible or out-of-scope states.
- Do not replace risk-based verification with exhaustive full-suite execution.
- Do not change product behavior, tests, generated assets, or runtime
  configuration in this change.

## Capabilities

### Modified Capabilities

- `engineering-harness`: Document MECE-style branch analysis as part of
  OpenSpec, TDD, and verification planning.
