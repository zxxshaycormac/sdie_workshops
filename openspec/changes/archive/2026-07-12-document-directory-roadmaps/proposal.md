## Why

The engineering harness currently gives contributors a root entry point and a
project-wide map, but it does not give a local orientation document inside each
top-level source directory. Contributors still have to infer whether a directory
is source, generated output, packaging, tests, runtime-adjacent configuration, or
OpenSpec state before they can make a safe change.

The root `AGENTS.md` is also close to the requested 200-line ceiling. Keeping
every detailed rule in that file makes it harder to scan and risks turning the
entry point into another long manual.

## What Changes

- Add a concise `ROADMAP.md` to every Git-tracked first-level directory.
- Use each directory roadmap to state the directory purpose, primary control
  paths, expected edits, verification signals, and boundaries.
- Compress `AGENTS.md` into a short operating entry point and route-map index.
- Move detailed HTTP API testing rules into `docs/engineering/API_TESTS.md` and
  link to them from `AGENTS.md`.
- Keep runtime and dependency directories out of the roadmap requirement unless
  they are tracked project source directories.

### Non-goals

- Do not document local runtime contents under `data/`, `log/`, `node_modules/`,
  `.make_evidence/`, `.git/`, or other untracked local state.
- Do not change product behavior, build output, test fixtures, database state, or
  generated assets.
- Do not replace `docs/engineering/PROJECT_MAP.md`; directory roadmaps should
  point contributors to concrete local orientation and then back to the project
  map for cross-layer tracing.

## Capabilities

### New Capabilities

- `directory-roadmap-guidance`: Directory-local roadmaps exist for tracked
  first-level directories and are indexed from the repository entry point.

### Modified Capabilities

- `engineering-harness`: The root operating instructions remain concise and link
  to detailed supporting documents rather than embedding every rule inline.

## Impact

### Observed control points

- Root instructions: `AGENTS.md`.
- Engineering docs: `docs/engineering/`.
- OpenSpec state: `openspec/changes/document-directory-roadmaps/`.
- Tracked first-level directories from `git ls-tree -d --name-only HEAD`.

### Compatibility surface

This is documentation-only. It does not change routes, APIs, persistence,
templates, assets, tests, or runtime configuration.
