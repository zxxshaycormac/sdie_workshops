# CLAUDE.md SDIE Restructure

**Date:** 2026-06-18
**Topic:** Restructure root `CLAUDE.md` around the SDIE methodology
**Status:** Approved

## Goal

Restructure the root `CLAUDE.md` so its top-level sections map directly to the four SDIE phases — Specification, Design, Implementation, Evaluation — with SDIE principles as a reference section. Primary audience is AI assistants (Claude/Codex), so the structure must be machine-parseable with explicit phase boundaries and pointers.

## Why

The current `CLAUDE.md` is organized around tooling commands (Build & Run, Tests, Lint, Code Generation) and an Architecture block. From an SDIE perspective:

- The (S)pecification layer is a single one-line pointer to `openspec/specs/spec.md`.
- The (D)esign layer is split awkwardly across "Architecture" and "Design Conventions."
- The (I)mplementation layer is implicit (Build/Lint/CodeGen/Conventions).
- The (E)valuation layer is missing entirely — no DoD, no Validation/Verification split, no cross-stage validation matrix.

A full SDIE restructure surfaces the missing Evaluation layer, makes the Specification phase actionable (delta-spec workflow, acceptance criteria), and gives AI agents explicit phase boundaries to navigate.

## Approach

Approach A: Full SDIE restructure. Four phase sections as the top-level skeleton, an SDIE Workflow Map framing the doc, and a brief SDIE Principles section at the end.

Rejected alternatives:
- **Approach B** (lighter SDIE framing): preserves current section order, only adds Evaluation section. Doesn't actually restructure; SDIE feels bolted on.
- **Approach C** (SDIE + workflow playbooks): adds concrete per-scenario playbooks. Too maintenance-heavy for now; could be added later.

## Scope

In scope:
- Rewrite root `/Users/genewu/github/gitea/CLAUDE.md` with SDIE-phase top-level structure.
- Preserve all existing content verbatim (commands, architecture blocks, conventions) — only reorganized.
- Add new (S)pecification workflow, (E)valuation layer with DoD, cross-stage validation matrix, SDIE Principles section.

Out of scope:
- Module-level `CLAUDE.md` files (cmd/, models/, etc.) — no changes.
- `design.md` — no changes.
- `openspec/specs/spec.md` — no changes.

## Design

### Top-level Structure

```
# CLAUDE.md
## SDIE Workflow Map
## (S)pecification — Specify + Business Context
## (D)esign — Architecture + Planning
## (I)mplementation — Development + Testing
## (E)valuation — Validation + Verification
## SDIE Principles
```

### Content Mapping

| Current section | New location | Notes |
|---|---|---|
| (intro) | `## SDIE Workflow Map` | Replaces one-line intro with phase map |
| `## Feature Specs` | `## (S)pecification` | Expanded: discovery, delta-spec workflow, acceptance criteria |
| `## Design Conventions (@design.md)` | `## (D)esign > Cross-cutting Conventions` | First subsection of Design |
| `## Architecture` (+ subsections) | `## (D)esign > Architecture subsections` | All subsections preserved verbatim |
| `## Build & Run` | `## (I)mplementation > Build & Run` | Preserved verbatim |
| `## Lint & Format` | `## (I)mplementation > Lint & Format` | Preserved verbatim |
| `## Code Generation` | `## (I)mplementation > Code Generation` | Preserved verbatim |
| `## Conventions` | `## (I)mplementation > Conventions` | Preserved verbatim |
| `## Tests` | `## (E)valuation > Tests` | Preserved verbatim |
| — (none) | `## (E)valuation > Contract Verification` | New |
| — (none) | `## (E)valuation > Definition of Done` | New, concise checklist |
| — (none) | `## (E)valuation > Cross-stage Validation` | New, S×I / D×I / S×E matrix |
| — (none) | `## SDIE Principles` | New, brief |

### Module-Level CLAUDE.md Reference

Per user feedback, do NOT enumerate the per-directory CLAUDE.md files. A single sentence points to them:

> Most directories ship their own `CLAUDE.md` layered on top of `design.md`. Read the one in the directory you're editing before writing code.

### New Subsections

**SDIE Workflow Map (top of doc):**
Brief intro to SDIE, the four phases, and the instruction to work top-down (read Spec → consult Design → implement → verify in Eval).

**(S)pecification expansion:**
- "When Starting Work" — three-step discovery: identify category, read feature entry, note cross-cutting impact.
- "When Proposing a Change" — OpenSpec delta spec workflow with ADDED/MODIFIED/REMOVED markers, pointer to `openspec/changes/` examples.
- "Acceptance Criteria" — BDD-style scenarios become the contract for (E)val.

**(D)esign additions:**
- `@design.md` becomes the first subsection (Cross-cutting Conventions).
- Brief Module-Level CLAUDE.md pointer (one sentence, no enumeration).
- "Planning" subsection: pointer to plan mode / Superpowers writing-plans skill for non-trivial changes.

**(I)mplementation additions:**
- "Local Rules Per Directory" pointer (one sentence).
- "TDD Discipline" pointer to Superpowers test-driven-development skill.

**(E)valuation expansion (all new):**
- Validation/Verification definitions at top.
- Contract Verification subsection (Swagger regen, migration + fixtures).
- Definition of Done — concise checklist:
  - [ ] Relevant spec category in `openspec/specs/` consulted
  - [ ] `design.md` principles honored
  - [ ] New code covered by tests
  - [ ] `make lint` passes
  - [ ] `make generate-swagger` regenerated if API changed
  - [ ] Migration + fixtures updated if DB schema changed
  - [ ] Acceptance criteria from (S)pec verified
- Cross-stage Validation matrix (D×I, S×I, S×E).

**SDIE Principles (brief, one line each):**
- MECE
- Token economics
- Harness Engineering + Steering Loop
- Cross-stage validation
- Real-time discarding

## Verification

- All current commands, code blocks, and bullet points appear unchanged in the new file (diff should show only reorganization + new sections).
- File renders as valid Markdown (no broken headings, code fences balanced).
- Each SDIE phase is reachable as a top-level `##` heading.

## Rollout

Single commit: `docs: restructure CLAUDE.md around SDIE phases`. No migrations, no breaking changes — the file is documentation only.
