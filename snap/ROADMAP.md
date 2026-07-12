# snap Roadmap

## Purpose

Contains Snap packaging metadata for distributing Gitea as a snap package.

## Route Map

- Packaging manifests and hooks describe how the application is built and run in
  Snap environments.

## When To Edit

Edit when changing Snap packaging behavior, exposed commands, confinement, or
packaged assets. Cross-check `docker/`, `contrib/`, and release workflows for
parallel packaging assumptions.

## Verification

Validate Snap metadata and run the package build in a Snap-capable environment
when possible.

## Boundaries

Do not use Snap metadata to change application runtime behavior outside the
packaging surface.
