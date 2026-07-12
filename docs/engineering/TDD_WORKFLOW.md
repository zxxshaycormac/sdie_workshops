# TDD-First Workflow

Use this workflow for changes that alter behavior, public contracts,
persistence, frontend consumption, asynchronous flows, or cross-layer
orchestration. Documentation-only, comment-only, formatting-only, and mechanical
changes can skip test-first work when the final report states that no behavioral
contract changed.

## Loop

1. Trace the control path before writing tests: route, context, handler, service,
   model, module, template, frontend entry, queue, or external boundary.
2. Capture the intended behavior in OpenSpec scenarios for non-trivial changes.
3. Ask an independent subagent to write or review the focused tests first when
   subagent support is available.
4. Run the focused test before implementation when practical and record the
   expected failure or the reason it cannot be run yet.
5. Implement the smallest production change that makes the focused test pass.
6. Run the same focused test again, then add broader checks only when the changed
   surface or risk requires them.
7. Refactor only while the focused tests stay green, and update OpenSpec tasks
   with the evidence.

## Subagent Role

The test subagent owns the test perspective, not the product implementation. It
should derive cases from OpenSpec scenarios and traced control points, inspect
neighboring test patterns, and cover meaningful success, failure, boundary,
authorization, ordering, and side-effect behavior for the change.

The subagent should not change production code unless the user explicitly gives
that scope. If no subagent is available, the main contributor must write or
review the failing test first and report that the subagent step was unavailable.

## Evidence

Final reports for behavior changes must include:

- the test or review that was created before production code;
- the pre-implementation failing command and failure reason, when practical;
- the post-implementation passing command;
- any broader verification that was intentionally skipped, with residual risk.

Do not describe a focused package or single Vitest file as "all tests pass".
Name the exact scope that ran.

## API And Frontend Contracts

HTTP API behavior follows `docs/engineering/API_TESTS.md`. Frontend work that
consumes API behavior also follows `web_src/API_CONTRACT.md`. In those changes,
write or review backend API contract tests and frontend consumption tests before
implementation whenever the contract or browser behavior changes.
