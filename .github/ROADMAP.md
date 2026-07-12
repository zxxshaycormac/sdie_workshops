# .github Roadmap

## Purpose

Defines GitHub collaboration metadata and GitHub Actions workflows.

## Route Map

- `ISSUE_TEMPLATE/` and `pull_request_template.md`: contributor-facing forms.
- `workflows/`: CI, release, translation, label, database, Docker, and E2E jobs.
- `labeler.yml`, `actionlint.yaml`, `FUNDING.yml`: repository automation and
  metadata.

## When To Edit

Edit when changing GitHub-hosted review, CI, release, or automation behavior.
Trace Make targets, scripts, environment variables, and required secrets before
changing a workflow.

## Verification

Run YAML validation or `actionlint` when available. For workflow changes, include
the affected job name and expected trigger in the final report.

## Boundaries

Never commit tokens, signing keys, or local credentials. GitHub secrets should be
referenced by name only.
