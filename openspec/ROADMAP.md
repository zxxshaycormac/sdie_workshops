# openspec Roadmap

## Purpose

Stores change-contract artifacts and durable specifications for the engineering
harness.

## Route Map

- `config.yaml`: OpenSpec store configuration.
- `changes/`: active and archived proposals, designs, tasks, and delta specs.
- `specs/`: current accepted system capabilities after archive.

## When To Edit

Edit before or alongside non-trivial behavior, contract, data, architecture, or
cross-layer changes. Keep artifacts synchronized with implementation evidence.

## Verification

Run `openspec validate --all --strict --no-interactive` and `openspec doctor`.
Use `openspec status --change <id>` when lifecycle order is unclear.

## Boundaries

OpenSpec records intent and requirements; it does not replace tests, code review,
or Git history.
