# Module-Level CLAUDE.md Review Checklist — Design Spec

**Date**: 2026-06-14
**Branch**: training/base_line
**Status**: Approved (application in progress)
**Source**: `/superpowers:brainstorming` session
**Reviews**: the doc set created by [2026-06-13-design-conventions-doc-design.md](./2026-06-13-design-conventions-doc-design.md)

## Goal

A reusable 6-criterion checklist for reviewing the 8 module-level `CLAUDE.md` files in this repo. The lens is **utility for AI assistant**: will an AI auto-loading these docs do the right thing?

## Background

The 8 module-level `CLAUDE.md` files exist for **progressive disclosure**. Claude Code auto-loads them when the working directory is inside that module. They sit on top of the always-loaded root `CLAUDE.md` + `design.md` (root references `design.md` via `@design.md`).

Their job is not to enumerate tasks; it is to add directory-specific depth that `design.md` does not carry. A section earns its slot only by teaching an AI something it would not otherwise know when working in that directory.

The doc set shipped in commit `9fbfdfd015` (2026-06-13). This checklist is the first review pass.

## The 6 Criteria

Each criterion: one-sentence definition + concrete failure modes. The failure modes anchor subjective criteria so two reviewers converge.

### 1. Actionable
The Rule line tells the AI what to do (imperative), not just describes what exists (indicative).

- Failure: "X is the pattern" / "X is used" without an action verb.
- Failure: passive voice that hides who acts ("transactions are wrapped" — by whom?).
- Failure: Rule line that only becomes imperative when read alongside the Why.

### 2. Unambiguous
Two reasonable readers would do the same thing after reading the Rule.

- Failure: vague quantifier without definition (`common`, `appropriate`, `pragmatic`) — unless surrounding text pins it down.
- Failure: conditional without the condition specified (`if needed`, `where appropriate`).
- Failure: unresolved conflict between rules — "match local style" vs "prefer new pattern" — without precedence.

### 3. Calibrated
`universal` / `common` / `aspirational` matches actual code frequency.

- Failure: `universal` claim with no grep evidence, or evidence that contradicts it.
- Failure: aspirational rule presented as established.
- **Operational definition (baked in):** `universal` ≈ ≥95% of relevant cases, no real exceptions in new code; `common` ≈ majority with documented exceptions; `aspirational` = stated goal not yet achieved. The `web_src/CLAUDE.md` "honest note up front" is the model.

### 4. Honest exceptions
Exceptions say "don't copy this" or "match style only when extending," never read as endorsement of the violating pattern.

- Failure: exception framed as an equally-valid alternative.
- Failure: "accepted technical debt" without saying what would be unaccepted.
- Failure: violating pattern listed in Exceptions without an explicit "do not copy in new code" signal.

### 5. Progressive-disclosure value
Section adds directory-specific context that `design.md` doesn't, and would help an AI auto-loading this file while working in this directory.

- Failure: duplicates `design.md` content — wastes the auto-load slot.
- Failure: covers material not relevant to working in this directory.
- Failure: misses directory-specific patterns (signature shapes, local helpers, file-layout conventions) that an AI working here would need but `design.md` doesn't carry.
- **Heuristic:** "If I deleted this section, would an AI working in this directory — already holding `design.md` — actually do worse work?"

### 6. Chaseable references
Every `see design.md Section X` / `see routers/web/CLAUDE.md §Y` resolves to content that actually covers the topic.

- Failure: reference to a section that doesn't exist or covers a different scope.
- Failure: reference too vague ("see design.md") to be actionable.
- **Note:** `design.md` references are realistic (auto-loaded); sibling-module `CLAUDE.md` references are cold-chases — flag them, because the AI may not have that file loaded.

## Walk-through Procedure

**Unit of review**: each numbered section within each file. Within a section, check the four standard fields (Rule / Why / Frequency / Exceptions) against each criterion that applies. Sections bundling multiple rules under one number (e.g. `models/CLAUDE.md` §7 "Additional Conventions") are accepted as-is — the review applies to each rule within.

**Cold-read first**: read the whole file once before judging any section. A section that looks weak in isolation may be supported by content elsewhere in the same file.

**Per-section procedure** — six passes, in order:

1. **Actionable** — read the Rule line as a standalone instruction. Does it tell you what to do? If it needs the Why to become imperative, flag.
2. **Unambiguous** — try to construct two different readings that lead to different code. Trip wires: `reasonable`, `appropriate`, `pragmatic`, `if needed`, `where appropriate`, `should consider`. Each is a flag unless nearby text pins it down.
3. **Calibrated** — Frequency: (a) is a frequency term present? (b) is grep evidence cited? (c) does the evidence support the term? Any mismatch → flag. Missing term or evidence → flag.
4. **Honest exceptions** — read Exceptions *as if you'd skipped the Rule*. Does it read as endorsement of the violating pattern? Specifically hunt for "accepted" / "fine" / "okay" without an explicit "do not copy" signal nearby.
5. **Progressive-disclosure value** — apply the deletion test: "If I removed this section, would an AI auto-loading this file do worse work in this directory?" Then cross-check against `design.md`: is any of this content already there? If yes → flag as redundancy.
6. **Chaseable references** — for every `see X Section Y`, resolve it. `design.md` references are realistic; sibling-module references are cold-chases — flag.

**Out of scope of this review** (don't do during this pass):
- Factual accuracy of grep counts against the current tree (separate review).
- Prose quality / stylistic nits.
- Proposing new sections or restructuring. Review produces findings only; fixes are local edits that preserve structure.

## Output Format

Per-file findings list. Each finding = section number + criterion + one-line description + optional one-line fix. No formal severity in v1; reviewer may informally group major/minor at write-up time.

```
## routers/web/CLAUDE.md
- §3 — Criterion 2 (Unambiguous): "the middleware cannot be forgotten when routes are added" reads as hyperbolic; obligation unclear. Tighten to a concrete failure mode.
- §7 — Criterion 5 (Progressive-disclosure): restates design.md §1 layering rule. Shorten to a directory-specific pointer.
```

## Scope

**In scope** (8 files):
- `cmd/CLAUDE.md`
- `models/CLAUDE.md`
- `modules/CLAUDE.md`
- `services/CLAUDE.md`
- `routers/web/CLAUDE.md`
- `routers/api/v1/CLAUDE.md`
- `templates/CLAUDE.md`
- `web_src/CLAUDE.md`

**Out of scope**:
- Root `CLAUDE.md` (different purpose — build/test commands + `@design.md` pointer).
- `design.md` itself (structurally different — no directory-specific depth to provide).
- `workshops/CLAUDE.md` (different lineage — training materials, not part of the design conventions doc set).

## Batching

Per-file batching. Each session: 1-2 files completely (cold-read + 6 criteria + fixes). Findings accumulate in a single report file across sessions. When resuming: re-read `design.md` + this checklist doc before continuing (~2 min).

Estimated cost: ~7 sections × 8 files × 6 criteria ≈ 336 micro-passes. Cold-read adds ~5 min/file. Total: ~60-90 min for one careful walk-through.

## What's Next

Apply this checklist to all 8 files in scope; fix issues inline; commit changes. Findings report location: `docs/superpowers/reviews/2026-06-14-claude-md-review-findings.md` (created if any findings need to outlive the fix commit).
