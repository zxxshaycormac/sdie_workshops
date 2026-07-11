## Context

The repository contains roughly 2,500 Go files, 700 Go tests, 150 JavaScript or Vue files, and 480 templates. Runtime behavior crosses routing, service orchestration, database models, Git-facing modules, templates, generated assets, and public API contracts. Existing upstream documentation is valuable but does not provide a compact operating model for AI-assisted changes.

The workspace also contains local Gitea state and currently has no Git metadata. The harness must therefore protect runtime-only paths and must not depend on Git to select verification scope.

## Goals / Non-Goals

**Goals:**

- Make the repository legible to a new human or AI contributor in minutes.
- Turn behavior changes into explicit, reviewable requirements before code edits.
- Route changes to concrete control points and quality gates.
- Provide fast focused checks and truthful evidence reporting.
- Demonstrate one complete OpenSpec workflow using the harness itself.

**Non-Goals:**

- Replace Gitea's existing contributor documentation or Make targets.
- Add a new build system, test framework, or runtime dependency.
- Automatically infer perfect verification scope from an absent Git diff.
- Change Gitea product behavior.

## Decisions

### Use layered, repository-native guidance

`AGENTS.md` will be the concise mandatory entry point. Detailed, slower-changing knowledge will live under `docs/engineering/`, and OpenSpec will own change contracts under `openspec/`.

Alternative considered: put all knowledge in `AGENTS.md`. Rejected because a very large instruction file consumes agent context on every task and becomes difficult for humans to navigate.

### Use OpenSpec for behavior and cross-layer changes

Non-trivial changes will use `proposal -> specs -> design -> tasks -> apply -> validate -> archive`. Trivial text-only changes may bypass artifacts when the contributor explicitly confirms no behavior or contract changed.

Alternative considered: require OpenSpec for every edit. Rejected because ceremony for typo-only changes would reduce adherence without adding useful traceability.

### Require explicit verification targets

The wrapper will expose `spec`, `go`, and `frontend` modes. Focused code modes require caller-supplied packages or test paths. Existing Make targets remain the source for broad lint, build, integration, and end-to-end suites.

Alternative considered: infer changed packages from Git. Rejected for the bootstrap because this workspace currently has no `.git` directory. Git-aware selection can be added later without changing the current interface.

### Treat generated and runtime paths as protected

The instructions and map will identify generated outputs and local state. Agents will work from source inputs and only inspect runtime state when the task explicitly calls for diagnosis.

Alternative considered: rely on `.gitignore`. Rejected because ignore rules do not explain source-of-truth relationships and do not prevent an agent from reading secrets.

## Requirement Traceability

| Requirement | Control point | Verification |
| --- | --- | --- |
| Repository operating instructions | `AGENTS.md` | Manual content review and link checks |
| Code-grounded navigation | `docs/engineering/PROJECT_MAP.md` | Path existence checks |
| Executable verification | `tools/harness/verify.sh` | Exercise usage, `spec`, and focused modes |
| Safety and evidence boundaries | `AGENTS.md`, project map, verification guide | Manual review plus shell guards |
| OpenSpec lifecycle guidance | `docs/engineering/OPEN_SPEC.md`, `openspec/config.yaml` | `openspec validate --all --strict` and `openspec doctor` |

## Risks / Trade-offs

- [Map becomes stale as code moves] -> Keep concrete paths and require map updates when a change introduces a new control point.
- [Contributors run only fast checks] -> Require verification plans in designs and final reports that distinguish completed and omitted checks.
- [OpenSpec becomes documentation theater] -> Specs require testable scenarios, tasks require adjacent checks, and archive is only after implementation evidence.
- [Shell wrapper hides project commands] -> Print every underlying command and keep broad canonical commands documented directly.
- [No Git safety net] -> State the limitation prominently and avoid Git-dependent automation until metadata is restored or initialized.

## Migration Plan

1. Configure OpenSpec and create this bootstrap change.
2. Add operating instructions, project map, workflow guide, and verification entry point.
3. Run strict OpenSpec validation and focused harness checks.
4. Mark tasks complete, archive the change, and let OpenSpec create the main `engineering-harness` specification.

Rollback is file-based: remove the added harness files and OpenSpec directories. No runtime data or schema rollback is required.

## Open Questions

- Whether this extracted workspace should be initialized as a new Git repository or reconnected to its upstream history.
- Which CI system will eventually enforce the harness checks; no CI configuration is changed in this bootstrap.
