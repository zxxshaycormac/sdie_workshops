# .devcontainer Roadmap

## Purpose

Defines development-container defaults for contributors who work in containerized
environments.

## Route Map

- `devcontainer.json`: editor image, extensions, ports, and bootstrap behavior.

## When To Edit

Edit when changing container-only developer experience. Runtime behavior,
production images, and CI should be traced through `docker/`, `.github/`, and
`build/` before changing this file.

## Verification

Validate JSON syntax and, when possible, reopen the repository in the dev
container to confirm the environment starts.

## Boundaries

Do not encode personal paths, local secrets, or machine-specific credentials.
