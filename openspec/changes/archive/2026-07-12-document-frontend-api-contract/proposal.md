## Why

The directory roadmaps added a frontend `web_src/ROADMAP.md`, but it only
described the frontend source directory. It did not explicitly document how
frontend code should share the backend API contract or how frontend generated
assets differ from backend-generated API artifacts.

## What Changes

- Add `web_src/API_CONTRACT.md` as the frontend-side guide for shared API
  contracts.
- Link the document from `web_src/ROADMAP.md`.
- Update `docs/engineering/API_TESTS.md` to require both backend API contract
  verification and frontend behavior verification when frontend code consumes an
  API change.

### Non-goals

- Do not generate a frontend API client or change Swagger generation.
- Do not change product behavior, routes, structs, tests, templates, or compiled
  assets.

## Capabilities

### Modified Capabilities

- `engineering-harness`: Document frontend participation in shared API contracts.
