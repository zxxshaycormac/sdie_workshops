# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Purpose

Training workshop materials for SDIE (Software Development in Existing) practices, using the Gitea codebase (parent directory) as the real-world training target. Content is bilingual — English and Chinese.

## Workshop Structure

| Workshop | Duration | SDIE Direction | Maintenance Type | Core Question |
|---|---|---|---|---|
| **W0** (Optional) | 2.5h | Pure reverse: code → understanding → spec | None (understanding only) | How to build a mental model before modifying code? |
| **W1** | 2.5h | Forward S→D→I→E | Add feature | How to add a feature using the SDIE process? |
| **W2** | 2.5h | Reverse E→S | Fix defect | How does Harness help safely fix bugs? |
| **W3** | 2.5h | Reverse I→S | System optimization (no behavior change) | How to guarantee refactoring doesn't change behavior? |
| **W4** | 2.5h | Reverse D→S | Adaptation/migration | How to encapsulate recurring patterns as Skills? |

Each workshop follows the same structure: concept review (30m) → instructor demo (15m) → hands-on practice (3 rounds with break) → review (10m).

## Candidate Analysis

Located in `candidates/`. All sourced from Gitea v1.23.0 release notes with linked GitHub issues:

- `feature_candidates.md` — 19 features, 4 analyzed in depth. **W1 exercise**: Tag search (#31998)
- `bugfix_candidates.md` — 36 bugfixes, 5 analyzed in depth. **W2 exercise**: Actions status aggregation (#32857)
- `refactor_candidates.md` — 5 refactors, 1 analyzed in depth. **W3 exercise**: Split mail sender sub-package (#18664)

`basic/code_summary.md` — Gitea codebase architecture overview with LOC statistics (~301K Go LOC).

## Workshop Exercises (Real Gitea Issues)

- **W1**: Tag search on repo tags page (Issue #31998) — forward SDIE with `brainstorming` + `writing-plans` + `TDD`
- **W2**: Actions job status aggregation bug (Issue #32857) — `systematic-debugging` + TDD + Hook creation (PreToolUse/PostToolUse)
- **W3**: Split mail sender from mailer service (Issue #18664) — MECE analysis + PEV cycle + Edit Lock Hook + Characterization Test
- **W4**: Issue time estimation (Issue #23112) — comprehensive D→S→I→E + `subagent-driven-development` + custom Skill creation

## Harness Progression

Workshops build Harness assets incrementally:

| Created In | Asset | Type |
|---|---|---|
| W0/W1 | `CLAUDE.md` | Reasoning feedforward |
| W2 | PreToolUse Hook (safety) + PostToolUse Hook (lint) | Computational feedback |
| W3 | Edit Lock Hook + Characterization Test | Computational feedback |
| W4 | Custom Skill (via `skill-creator`) | Reasoning feedforward |

## Key Concepts

- **SDIE**: Specification → Design → Implementation → Evaluation
- **MECE**: Mutually Exclusive, Collectively Exhaustive — used in W0 for reverse engineering self-checks, in W3 for blast radius control
- **PEV**: Plan-Execute-Verify cycle for safe refactoring (W3)
- **Harness 2x2 matrix**: feedforward/feedback × computational/reasoning (W2)
- **Agent = Model + Harness** (W2)
- **Progressive disclosure**: CLAUDE.md (always loaded) → Skills (on-demand) → External docs (deep reference)

## Workflow

Workshops reference the Gitea source at `../` (parent directory). All PR/issue links point to `go-gitea/gitea` on GitHub. **Students must not look at the actual fix PRs** — each workshop is designed for guided discovery via Explore Agent and SDIE process.

## Excluded: `NOLOOKINSIDE/`

The `NOLOOKINSIDE/` folder (under `workshops/`) contains training reference samples. **Do NOT read, reference, search, or use any content from this folder** when working on workshop exercises or any task governed by this file. Treat it as if it does not exist. This ensures exercises are solved via the SDIE process rather than by copying provided samples.
