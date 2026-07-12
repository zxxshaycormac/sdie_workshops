# modules Roadmap

## Purpose

Contains reusable infrastructure, public structs, Git helpers, settings,
rendering helpers, queues, storage, and cross-cutting utilities.

## Route Map

- `structs/`: public API request and response contracts.
- `setting/`: configuration loading and defaults.
- `git/`, `gitrepo/`: Git object, repository, and command helpers.
- `queue/`, `storage/`, `cache/`, `log/`, `markup/`: shared infrastructure.
- `templates/`, `translation/`, `validation/`: support used by routers and
  services.

## When To Edit

Edit when changing shared primitives or public contracts. Search broadly before
changing exported types or helpers because callers often span API, web, service,
model, and tests.

## Verification

Run focused package tests under the changed module. For `structs` or `setting`
changes, also run affected API or startup tests.

## Boundaries

Avoid importing high-level web or service orchestration into low-level modules.
Keep dependency direction stable.
