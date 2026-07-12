# routers Roadmap

## Purpose

Owns HTTP route registration, request binding, context use, permissions, and web
or API response handling.

## Route Map

- `web/`: browser routes and server-rendered page handlers.
- `api/v1/`: public API routes, Swagger annotations, and handlers.
- `private/`: internal routes used by hooks, runners, and server components.
- `common/`: shared router initialization helpers.

## When To Edit

Edit when changing route shape, authorization, binding, status codes, response
fields, template data, redirects, or web/API behavior. Trace services, models,
modules, templates, and tests before changing a handler.

## Verification

Run focused handler tests and routed integration tests. HTTP API changes must
also follow `docs/engineering/API_TESTS.md`.

## Boundaries

Keep business orchestration in `services/` and persistence in `models/`; routers
should coordinate request/response behavior.
