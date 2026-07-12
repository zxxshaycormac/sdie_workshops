## Context

The shared API contract is currently anchored in backend source files:
`routers/api/`, `modules/structs/`, `services/`, `models/`, Swagger annotations,
and generated Swagger templates. Frontend code under `web_src/` consumes those
contracts through browser flows, Vue components, JavaScript modules, and
template-provided initialization data.

The engineering harness needs an explicit frontend document so future changes do
not treat backend API tests or frontend tests as interchangeable evidence.

## Decisions

### 1. Put the frontend-facing guide under `web_src/`

The document belongs where frontend contributors will look first. It points back
to backend sources of truth and engineering docs instead of duplicating API
details.

### 2. Keep backend API tests authoritative for HTTP contracts

Frontend tests prove that browser code consumes the contract correctly. They do
not prove route status codes, response schemas, validation, or persistence side
effects. `docs/engineering/API_TESTS.md` remains the backend API contract test
guide and now links to the frontend guide.

### 3. Do not introduce generated client code

This change is documentation-only. If the project later adds generated frontend
types or clients, that should be a separate OpenSpec change.

## Verification Plan

- Run strict OpenSpec validation.
- Check Markdown and whitespace with `git diff --check`.
- Confirm `AGENTS.md` remains below 200 lines.
