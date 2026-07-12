# .codex Roadmap

## Purpose

Stores repository-local Codex skills that adapt the OpenSpec workflow to this
codebase.

## Route Map

- `skills/openspec-*`: conversational entry points for explore, propose, apply,
  update, sync, and archive flows.
- Skills should point back to `docs/engineering/` and `openspec/`, not duplicate
  full project knowledge.

## When To Edit

Edit this directory only when changing local agent workflow instructions.
Product behavior changes belong in source directories and OpenSpec artifacts.

## Verification

Run `openspec validate --all --strict --no-interactive` and exercise the relevant
skill path manually when behavior changes.

## Boundaries

Do not store credentials, local browser state, terminal transcripts, or one-off
task notes here.
