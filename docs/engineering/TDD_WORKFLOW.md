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
3. Build a MECE branch list for the relevant behavior space before selecting
   tests.
4. Ask an independent subagent to write or review the focused tests first when
   subagent support is available.
5. Run the focused test before implementation when practical and record the
   expected failure or the reason it cannot be run yet.
6. Implement the smallest production change that makes the focused test pass.
7. Run the same focused test again, then add broader checks only when the changed
   surface or risk requires them.
8. Refactor only while the focused tests stay green, and update OpenSpec tasks
   with the evidence.

## MECE Branch Analysis

In this repository, MECE means splitting the relevant behavior space into
mutually exclusive categories and making the list maximally complete for the
current contract. It is bounded by the OpenSpec scope: impossible, unreachable,
or non-goal branches should be named as exclusions instead of silently ignored.

Use MECE analysis before asking for tests. Consider at least these branch
families when they apply:

- actor, permission, ownership, and authentication state;
- input shape, empty values, invalid values, ordering, pagination, and limits;
- persisted data state, migrations, transactions, and supported databases;
- feature flags, configuration, environment, locale, and template-provided data;
- downstream service success, failure, timeout, retry, and partial result;
- concurrency, time, idempotency, ordering, and cache behavior;
- API status, response shape, UI state, generated output, and side effects.

The output does not need to be a large table for every small change. It must be
clear enough that the chosen tests can be traced back to the important branches,
and that intentionally skipped branches have a reason and residual risk.

## Subagent Role

The test subagent owns the test perspective, not the product implementation. It
should derive cases from OpenSpec scenarios, MECE branch analysis, and traced
control points; inspect neighboring test patterns; and cover meaningful success,
failure, boundary, authorization, ordering, and side-effect behavior for the
change.

The subagent should not change production code unless the user explicitly gives
that scope. If no subagent is available, the main contributor must write or
review the failing test first and report that the subagent step was unavailable.

## Evidence

Final reports for behavior changes must include:

- the test or review that was created before production code;
- the important MECE branches covered by that test or review;
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
