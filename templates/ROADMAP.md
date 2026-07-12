# templates Roadmap

## Purpose

Stores Go templates for server-rendered pages, emails, feeds, and generated
Swagger template output.

## Route Map

- `repo/`, `user/`, `admin/`, `org/`, and related directories render web pages.
- `mail/`: email bodies.
- `swagger/`: generated Swagger template output.
- Shared partials provide reusable UI structure.

## When To Edit

Edit when changing rendered UI, template data expectations, emails, or feeds.
Trace the handler data in `routers/` and any browser behavior in `web_src/`
before modifying markup.

## Verification

Run the focused integration test for the rendered page or flow. For browser
behavior, add Vitest or E2E coverage as appropriate.

## Boundaries

Do not hand-edit generated Swagger output; update API annotations and generator
sources instead.
