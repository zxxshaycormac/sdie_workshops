# web_src Roadmap

## Purpose

Contains frontend source code for JavaScript, Vue components, CSS, themes, and
browser behavior that are compiled into public assets.

## Route Map

- `js/`: frontend initialization, feature modules, Vue components, and tests.
- `css/`: source styles and themes.
- `API_CONTRACT.md`: frontend rules for shared backend/frontend API contracts.
- Build output lands under `public/assets/`.

## When To Edit

Edit when changing browser-side behavior, styling, or frontend components. Trace
server-rendered template data in `templates/` and handler data in `routers/`
before changing initialization assumptions. For shared API behavior, read
`API_CONTRACT.md` before editing frontend consumers.

## Verification

Run focused Vitest for JS behavior and frontend build/lint when source bundles
change. Use E2E or visual checks for user-visible workflows. Shared API changes
also require the backend API checks described in `docs/engineering/API_TESTS.md`.

## Boundaries

Do not edit `public/assets/` directly for frontend source changes.
