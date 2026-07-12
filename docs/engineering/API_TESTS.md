# HTTP API Testing Contract

Any change that affects HTTP API behavior MUST update the corresponding API
tests in the same change. This includes routes, handlers, request and response
types in `modules/structs`, binding or validation, status codes, response
fields, authentication or middleware, business state transitions, database
reads/writes, and calls to downstream services or infrastructure.

API behavior changes follow the repository TDD workflow in
`docs/engineering/TDD_WORKFLOW.md`: write or review the focused API tests before
production code changes, run them to observe the expected failure when
practical, then implement the smallest change that makes them pass. Use MEMC
branch analysis to decide which route, input, permission, state, failure, and
side-effect branches belong in focused API coverage.

Frontend and backend changes share this same API contract. Backend source owns
the authoritative route, request/response structs, validation, status behavior,
and Swagger source. Frontend source under `web_src/` owns browser consumption of
that contract. When a frontend change depends on an API field, error shape,
permission, ordering, or side effect, also read `web_src/API_CONTRACT.md` and
verify both the backend API contract and the browser behavior.

Before adding or changing HTTP integration tests, read
`tests/integration/README.md` and inspect adjacent
`tests/integration/api_*_test.go` files. Reuse the existing test server, fixture
database, login sessions, access tokens, `MakeRequest` helpers,
`tests.PrepareTestEnv(t)`, and cleanup patterns. For focused handler tests under
`routers/api/`, follow neighboring tests that use `contexttest.MockAPIContext`,
`unittest.PrepareTestEnv`, and related helpers. Reuse existing queue, storage,
mail, or external-service test abstractions; do not introduce ad hoc monkeypatch
or mock conventions.

Before changing an interface, search all existing coverage by route path,
handler name, `modules/structs` type, service or model symbol, and business
keyword:

- focused handler and package tests under `routers/api/**/*_test.go`;
- routed HTTP contract and flow tests under `tests/integration/api_*_test.go`;
- related business, persistence, conversion, and infrastructure tests under
  `services/`, `models/`, and `modules/`.

If the affected interface has no coverage, add it. If coverage exists, update
requests, assertions, fixtures, and dependencies so they represent the intended
behavior. API tests should cover the relevant success path, failure path, boundary
case, authorization behavior, and important side effects. If a business order is
part of the public contract, verify it through an integration flow rather than
only an isolated handler call.

When business logic and API tests both change, use an independent subagent for
API test implementation or review before product code changes when that
capability is available. If it is not available, the final report MUST include
an independent review checklist for route reachability, status codes, response
body assertions, authorization, fixtures, side effects, and cleanup.

After implementation, run the affected API tests and diagnose failures. Fix
product or compatibility defects when tests expose them. Only update old
expectations when they no longer represent the correct behavior. Do not delete
assertions, skip tests, or reduce coverage merely to make the suite pass.

The final report MUST list every added or updated API test, the important MEMC
branches covered, the expected failing evidence before implementation or why it
was not practical, the exact commands run after implementation, their results,
and any related tests not run with the reason and residual risk.
