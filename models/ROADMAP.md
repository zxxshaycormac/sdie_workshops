# models Roadmap

## Purpose

Owns persistence models, query construction, migrations, fixtures, and database
contracts.

## Route Map

- Domain packages such as `repo`, `user`, `issues`, `git`, `actions`, and
  `packages` define persisted state and queries.
- `db/`: database engine helpers, transactions, iteration, collation, and search
  helpers.
- `migrations/`: schema and data migrations.
- `fixtures/`: test data loaded by unit and integration tests.

## When To Edit

Edit when changing stored data, query semantics, migrations, or model-level
business invariants. Trace callers in `services/`, `routers/`, `modules/`, and
templates before changing field shape.

## Verification

Run focused package tests such as `go test ./models/<pkg>/...`. For schema or
cross-database behavior, run migration tests and the relevant database-specific
integration target.

## Boundaries

Do not put HTTP concerns here. Keep request binding in `routers/` or
`services/forms/`; keep orchestration in `services/`.
