# docker Roadmap

## Purpose

Contains Docker image documentation and manifest templates.

## Route Map

- `README.md`: Docker usage notes.
- `manifest.tmpl`, `manifest.rootless.tmpl`: image manifest templates.
- Runtime image behavior may also involve `Dockerfile`, `Makefile`, `build/`,
  and `.github/workflows/`.

## When To Edit

Edit when changing container image metadata, rootless behavior, or Docker-facing
documentation. Trace CI and release workflows before changing manifests.

## Verification

Run the relevant Docker dry-run, build, or workflow-equivalent target when
available. At minimum, validate template syntax and affected image names.

## Boundaries

Do not place local registry credentials or deployment-specific secrets here.
