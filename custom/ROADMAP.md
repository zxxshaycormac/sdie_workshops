# custom Roadmap

## Purpose

Provides source examples for site customization and local configuration shape.

## Route Map

- `conf/app.example.ini`: documented configuration example.
- `conf/app.ini`: local runtime configuration in this workspace.

## When To Edit

Edit examples when changing documented configuration defaults or options. Only
touch local runtime config when the task explicitly concerns this local Gitea
service.

## Verification

For example changes, compare with `modules/setting/` defaults and documentation.
For local runtime diagnosis, restart Gitea and verify the specific setting.

## Boundaries

Treat `custom/conf/app.ini` as sensitive local state. Do not quote secrets,
tokens, internal URLs, or signing keys in commits, OpenSpec artifacts, or final
reports.
