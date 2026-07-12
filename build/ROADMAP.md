# build Roadmap

## Purpose

Contains build, generation, formatting, locale, license, and test-environment
helpers used by Make targets and CI.

## Route Map

- `generate-*.go`: source generators for bindata, emoji, licenses, gitignores,
  and related outputs.
- `codeformat/`: formatting tools and tests.
- `test-env-*.sh`: local and CI environment preparation.
- `update-locales.sh`: localization update helper.

## When To Edit

Edit when changing how generated files, checks, or build evidence are produced.
Trace the Make target and generated outputs before modifying a helper.

## Verification

Run the narrow generator or formatter test first, then the Make target that calls
it. Inspect diffs for generated outputs.

## Boundaries

Do not use build helpers to hide product test failures. Generated files should
come from the documented source and generator path.
