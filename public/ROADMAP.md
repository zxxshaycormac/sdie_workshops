# public Roadmap

## Purpose

Contains public web assets served by the application, including generated
frontend bundles.

## Route Map

- `assets/`: generated frontend output from `web_src/` and build tooling.
- Static public files may be served directly by the web server.

## When To Edit

Prefer editing source files under `web_src/`, `templates/`, `options/`, or
`assets/` and then running the generator. Direct edits are only appropriate for
true source files in this directory.

## Verification

Run the frontend build or specific asset generator and inspect generated diffs.
For UI behavior, use focused Vitest or E2E checks.

## Boundaries

Do not patch generated bundles by hand to make a test pass.
