# options Roadmap

## Purpose

Stores option data, locale source files, default templates, labels, licenses,
and other configurable resources bundled with the application.

## Route Map

- `locale/`: source localization strings.
- `gitignore/`, `label/`, `license/`, `readme/`: repository creation defaults.
- Other option bundles are loaded by settings, templates, or bindata generation.

## When To Edit

Edit when changing bundled defaults or source locale text. Trace consumers in
`modules/options`, `modules/translation`, templates, and bindata generation.

## Verification

Run locale or bindata generation when required, plus the focused UI/API test that
uses the changed option. Do not hand-edit generated bindata as the primary fix.

## Boundaries

Non-English locale files may be managed by translation tooling. Follow the
project localization workflow before editing translated strings directly.
