# docs Roadmap

## Purpose

Stores user documentation, development documentation, static doc assets, and the
engineering harness guides.

## Route Map

- `content/`: localized product documentation.
- `static/`: images and static files referenced by docs.
- `engineering/`: repository operating model, OpenSpec guidance, project map,
  verification, and API testing contract.
- `README*.md`: docs build and language entry points.

## When To Edit

Edit when behavior, configuration, setup, API usage, or engineering workflow
changes. User-visible behavior changes should update docs in the same OpenSpec
change.

## Verification

Run docs build or link checks when available. For engineering docs, run strict
OpenSpec validation and verify cross-links.

## Boundaries

Do not document secrets from local runtime files. Keep durable project guidance
here; one-off task notes belong in OpenSpec change artifacts or final reports.
