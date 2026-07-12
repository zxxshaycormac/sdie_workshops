# assets Roadmap

## Purpose

Stores small source assets and metadata used by generators, documentation, and
branding.

## Route Map

- `logo.svg`, `favicon.svg`: source visual identity assets.
- `emoji.json`, `go-licenses.json`: metadata consumed by build or generation
  tools.

## When To Edit

Edit source assets here when the underlying identity or metadata changes. If an
asset is generated, trace the generator under `build/` first.

## Verification

Run the relevant generator or build target named by the changed asset. For SVG
changes, inspect rendering in a browser or generated output.

## Boundaries

Do not edit generated bundles under `public/assets/` as a substitute for source
asset changes.
