# .gitea Roadmap

## Purpose

Stores Gitea-native repository metadata such as issue templates used when this
project is hosted on a Gitea instance.

## Route Map

- `issue_template.md`: default issue text for Gitea-hosted collaboration.

## When To Edit

Edit when changing repository collaboration templates for Gitea. GitHub-specific
automation belongs under `.github/`.

## Verification

Check Markdown rendering in the target Gitea instance when the template changes.

## Boundaries

Do not place runtime server configuration here; local server state lives under
`custom/`, `data/`, and `log/` and is normally outside source edits.
