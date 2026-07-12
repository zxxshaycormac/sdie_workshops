# tests Roadmap

## Purpose

Contains integration, end-to-end, fixtures, and test infrastructure that verify
real routed behavior and repository workflows.

## Route Map

- `integration/`: HTTP and database-backed integration tests.
- `e2e/`: Playwright-based browser tests.
- `fixtures/` and helper packages: shared test data and setup.
- Repository fixtures mirror Git behavior used by tests.

## When To Edit

Edit when product behavior changes need routed or browser-level coverage, or
when test infrastructure itself changes. Read adjacent tests before adding a new
pattern.

## Verification

Run the narrow Make target or package test named by the changed flow. API changes
must follow `docs/engineering/API_TESTS.md`.

## Boundaries

Do not reduce assertions or skip tests only to pass a suite. Fix the product
defect or update obsolete expectations explicitly.
