## Context

The repository already has a root `AGENTS.md`, an OpenSpec workflow guide, a
project map, and a verification guide. These documents describe how to trace a
change across the codebase, but contributors still need a local entry point when
they land in a top-level directory.

The current root `AGENTS.md` is 164 lines before this change. Adding a full
directory roadmap table without moving detail out would push the file close to
the requested limit and make the entry point harder to scan.

The Git-tracked first-level directories are:

```text
.codex .devcontainer .gitea .github assets build cmd contrib custom docker docs
models modules openspec options public routers services snap templates tests
tools web_src
```

Local runtime/dependency directories such as `.git`, `data`, `log`,
`node_modules`, `.make_evidence`, and `sqlite-log` are intentionally excluded
from the per-directory roadmap requirement because they are not tracked source
directories and can contain mutable or sensitive local state.

## Decisions

### 1. Use `ROADMAP.md` as the common file name

Each tracked first-level directory gets a `ROADMAP.md`. The common name makes the
documents easy to discover, link, and audit.

### 2. Keep roadmaps short and operational

Each roadmap uses the same compact shape:

- purpose;
- route map;
- when to edit;
- verification signals;
- boundaries.

This keeps the files useful for navigation without turning them into duplicated
package documentation.

### 3. Make root `AGENTS.md` an index, not a manual

The root instructions remain the mandatory starting point, but long secondary
guidance is moved into dedicated files. In particular, HTTP API testing rules
move to `docs/engineering/API_TESTS.md`. `AGENTS.md` links to that file and to
every directory roadmap.

### 4. Preserve the harness's source/runtime boundary

Roadmaps are added only to tracked project directories. Runtime directories
remain documented as boundaries in `AGENTS.md` and `PROJECT_MAP.md`, not as
places where contributors should add source-like documentation.

## Verification Plan

- Confirm `AGENTS.md` stays below 200 lines.
- Confirm every Git-tracked first-level directory has `ROADMAP.md`.
- Run strict OpenSpec validation.
- Run a documentation-focused diff check.

## Risks / Trade-offs

- [Risk] A roadmap can become stale if it tries to document every package.
  -> Keep each roadmap at the directory boundary level and point to tracing.
- [Risk] Hidden tracked directories could be missed.
  -> Generate the checklist from `git ls-tree -d --name-only HEAD`.
- [Trade-off] Runtime directories do not receive local roadmaps.
  -> This preserves the existing harness rule that runtime state is not source.
