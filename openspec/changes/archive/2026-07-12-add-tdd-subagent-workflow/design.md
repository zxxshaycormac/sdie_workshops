## Context

The current harness asks contributors to trace code paths, propose non-trivial
changes in OpenSpec, and select verification before coding. That is close to a
testable workflow, but it lacks an explicit red-green-refactor sequence and does
not define how a subagent should participate in test design.

This change keeps the root instructions short and puts the operational detail
under `docs/engineering/`, matching the existing directory roadmap and API
testing structure.

## Decisions

### 1. Add a focused TDD workflow document

`docs/engineering/TDD_WORKFLOW.md` becomes the canonical test-first guide. It
describes the loop from traced contract to OpenSpec scenarios, subagent test
work, expected failure evidence, minimal implementation, green verification, and
refactoring.

### 2. Make the subagent own test perspective

When available, the subagent should write or review tests before production code
changes. The subagent derives tests from OpenSpec scenarios and traced control
points, covers meaningful success/failure/boundary/authorization/side-effect
cases, and does not implement product code unless the user explicitly assigns
that scope.

### 3. Keep exceptions explicit

Pure documentation, comment, formatting, and other mechanical changes can skip
test-first work, but the final report must say why there was no behavioral
contract to test. If the environment cannot run the expected failing test, the
limitation must be reported instead of pretending the TDD step happened.

### 4. Tie API and frontend contract docs into the same loop

HTTP behavior changes use `docs/engineering/API_TESTS.md`. Frontend changes that
consume API behavior also use `web_src/API_CONTRACT.md`. The TDD workflow should
point to both rather than duplicating their contract-specific rules.

## Verification Plan

- Validate the new OpenSpec change strictly.
- Run strict validation for all OpenSpec artifacts.
- Check Markdown and whitespace with `git diff --check`.
- Confirm root `AGENTS.md` remains below 200 lines.
