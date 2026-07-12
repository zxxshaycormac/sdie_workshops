## Context

The repository already asks contributors to trace real control points, write
OpenSpec scenarios, and use TDD-first verification. The missing step is a
documented way to avoid shallow scenario selection. Without branch analysis,
contributors may write a failing test for the happy path and still miss
authorization, state, boundary, side-effect, or failure branches.

## Decisions

### 1. Define MEMC locally as practical branch analysis

In this harness, MEMC means splitting the relevant behavior space into
mutually exclusive categories and making the list maximally complete for the
current contract. This is intentionally bounded by OpenSpec scope and explicit
non-goals, so it does not demand tests for impossible or irrelevant states.

### 2. Put operational detail in the TDD workflow

`docs/engineering/TDD_WORKFLOW.md` is where future agents decide which tests a
subagent should write before implementation. The MEMC checklist belongs there
because it directly shapes test design.

### 3. Link OpenSpec, API, and verification guidance

OpenSpec scenarios should be reviewed against the MEMC branch list. API tests
should use it to select success, failure, boundary, authorization, and
side-effect coverage. Final reports should state which important branches were
covered and which were intentionally left out.

## Verification Plan

- Validate the new OpenSpec change strictly.
- Run strict validation for all OpenSpec artifacts.
- Check Markdown and whitespace with `git diff --check`.
- Confirm root `AGENTS.md` remains below 200 lines.
