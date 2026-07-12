# services Roadmap

## Purpose

Contains use-case orchestration across models, modules, notifications, queues,
repository operations, forms, and background workflows.

## Route Map

- `forms/`: web form contracts and validation helpers.
- `context/`: request context and repository/user loading.
- Domain services such as `repository`, `issue`, `pull`, `release`, `actions`,
  `webhook`, `mirror`, `task`, and `cron`.
- `convert/`: conversion between internal models and public structs.

## When To Edit

Edit when a change coordinates multiple models/modules or owns business flow.
Trace callers from `routers/` and side effects into queues, notifications,
storage, Git, or persistence.

## Verification

Run focused service tests and the nearest routed integration flow for user-facing
behavior. For async work, verify producer, payload, consumer, retry, and final
state.

## Boundaries

Do not put raw HTTP binding here unless it belongs in `forms/`. Do not bypass
model transaction boundaries.
