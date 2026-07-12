# Frontend API Contract

## Purpose

Frontend code and backend code share one API contract. The backend owns the
authoritative HTTP route, request, response, validation, and Swagger source; the
frontend owns how that contract is consumed in browser flows and how generated
frontend assets are produced.

## Shared Contract Sources

- Route and status behavior: `routers/api/` and adjacent handlers.
- Request and response shapes: `modules/structs/`.
- Business orchestration and side effects: `services/`.
- Persistence behavior: `models/`.
- Swagger source and generated template path: `routers/api/v1/swagger/` and
  `templates/swagger/v1_json.tmpl`.
- Frontend consumers: `web_src/js/`, Vue components, and any template-provided
  initialization data.

When one side changes, trace both sides before implementation. A frontend change
that assumes a new field, status, error shape, permission, or ordering must point
to the backend source of truth and the test that proves it.

## Frontend Route Map

1. Find the browser entry in `web_src/js/index.js` or the feature module.
2. Trace any Vue component, data attribute, DOM hook, or template-provided value.
3. Find the API call helper or fetch wrapper used by the feature.
4. Trace the backend route, handler, `modules/structs` type, service, and model.
5. Check whether Swagger or docs need regeneration.
6. Add or update frontend tests for browser behavior and backend API tests for
   the contract itself.

## Generated Output Boundary

Frontend source lives under `web_src/`. Compiled assets land under
`public/assets/` and must not be hand-edited. If a shared API contract affects a
generated frontend artifact, update the source or generator and then run the
documented build target.

## Verification

For shared API work, combine backend and frontend checks:

- Backend API contract: follow `docs/engineering/API_TESTS.md`.
- Frontend behavior: run the focused Vitest or E2E test that covers the browser
  workflow.
- Generated assets: run the frontend build or generator only when source changes
  require it.

Do not claim the shared API contract is verified by frontend tests alone.
Frontend tests prove browser consumption; backend API tests prove the HTTP
contract.
