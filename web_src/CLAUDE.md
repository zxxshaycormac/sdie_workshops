# `web_src/` — Design Conventions

`web_src/` holds the JavaScript, CSS, Vue components, Fomantic UI build, and SVG assets that webpack bundles into the static frontend loaded by every browser request. This file goes deeper than `design.md` on frontend-specific rules; the cross-cutting rules remain authoritative and are not repeated here.

**Honest note up front.** The Gitea frontend has accumulated layers (jQuery-era vanilla JS, Fomantic UI, Vue, Tailwind) and the conventions are less codified than the Go backend. Rules below describe the current dominant pattern with grep evidence; aspirational guidance is flagged as such.

---

## 1. JavaScript entry — `web_src/js/index.js` is the single bootstrap

**Rule**: All page-side JS lives under `web_src/js/` and is wired through `web_src/js/index.js`, which is the `index` webpack entry declared in `webpack.config.js:77`. Each module exports an `init*` function (e.g. `initRepoEditor`, `initGlobalDropzone`); `index.js` imports and calls them at the bottom of the file. Do not add new entry points or self-invoking scripts.
**Why**: A single bootstrap means page-specific behavior is grep-discoverable in one file, every init function goes through the same lifecycle (after `bootstrap.js` which sets `__webpack_public_path__` and the global error handler), and lazy-loaded chunks stay predictable. The `entry` block at `webpack.config.js:77-96` lists only five production entries (`index`, `webcomponents`, `swagger`, `eventsource.sharedworker`, plus per-theme CSS) — adding more fragments the asset graph.
**Frequency**: dominant. `grep -rn "^export function init\|^export const init" web_src/js/ | wc -l` returns 118 init exports; `index.js` contains ~166 init calls. Sample: `web_src/js/index.js:4` imports `initRepoActivityTopAuthorsChart` from `./components/RepoActivityTopAuthors.vue`.
**Exceptions**: `web_src/js/webcomponents/` has its own `index.js` and webpack entry (`webcomponents`) because custom elements must be registered before any template renders. `web_src/js/standalone/swagger.js` is a separate entry because the Swagger page is standalone. `web_src/js/features/eventsource.sharedworker.js` is a Web Worker entry.

---

## 2. Feature/page-based module layout under `web_src/js/`

**Rule**: Organize new JS by feature, not by type. Put feature/page-specific code in `web_src/js/features/<feature>.js` (or `web_src/js/features/admin/<area>.js`); shared helpers in `web_src/js/utils.js` or `web_src/js/utils/`; cross-feature modules (sortable, tippy, toast) in `web_src/js/modules/`; markup rendering helpers in `web_src/js/markup/`; standalone entry-only code in `web_src/js/standalone/`.
**Why**: A flat `features/` directory (~70 files, e.g. `repo-editor.js`, `repo-projects.js`, `notification.js`, `admin/users.js`) means `grep -rln "initRepo" web_src/js/features/` finds every repo-related initializer. A type-based layout (`controllers/`, `services/`) would not match how Gitea actually dispatches by page context.
**Frequency**: dominant. Directories verified: `web_src/js/features/` (per-feature files), `web_src/js/features/admin/`, `web_src/js/markup/` (8 files: `anchors.js`, `content.js`, `mermaid.js`, etc.), `web_src/js/render/` (`ansi.js`, `pdf.js`), `web_src/js/modules/` (`fetch.js`, `fomantic.js`, `sortable.js`, `tippy.js`, `toast.js`, `stores.js`), `web_src/js/utils/`.
**Exceptions**: `web_src/js/vendor/` holds vendored shims — do not add new vendor code without discussion.

---

## 3. Avoid inline `<script>` in templates — write an `init*` instead

**Rule**: New behavior must not be added as inline `<script>` inside `templates/**/*.tmpl`. Put the code in a JS module, export an `init*`, and call it from `web_src/js/index.js`. Inline scripts break CSP-friendly bundling, escape webpack's chunk loading, and cannot reuse the CSRF/fetch wrapper.
**Why**: webpack only processes what is imported; inline scripts silently bypass module resolution, lazy-loading, and the global error handler set in `web_src/js/bootstrap.js:1`.
**Frequency**: dominant as a rule. `grep -rln "<script>" templates/ | wc -l` returns 4 — all legacy: `templates/base/head_script.tmpl` (bootstraps `window.config` before JS loads, must stay inline), `templates/repo/clone_script.tmpl` (synchronous clone-button init to avoid flicker; the file itself has a TODO to remove it), `templates/repo/diff/box.tmpl`, `templates/shared/combomarkdowneditor.tmpl`. These are exceptions, not precedents.
**Exceptions**: Passing server-rendered config into `window.config` (as `head_script.tmpl` does) is the sanctioned pattern — it must run before the JS bundle. Small per-template initializers that read DOM before paint may be tolerated if they are flagged with a TODO and have a structural reason (e.g. avoiding FOUC).

---

## 4. Vue components are single-file `.vue`, Options API, used sparingly

**Rule**: New interactive UI goes in `web_src/js/components/<Name>.vue` (single-file component), exports an `init<Name>` function (alongside the default-exported component) that mounts to a DOM element by selector, and is called from `index.js`. Use the Vue **Options API** (`export default { components, props, data, mounted, methods }`) — not `<script setup>`.
**Why**: Vue is used only for genuinely stateful/re-render-heavy widgets (charts, file trees, autocomplete selectors), not for full pages. The Options API is what every existing component uses (`__VUE_OPTIONS_API__: true` is forced in `webpack.config.js:202`); `<script setup>` would diverge from the 16 existing components. The `init*` mount pattern keeps Vue opt-in per page rather than a global root.
**Frequency**: limited. Only 16 `.vue` files exist (`find web_src -name "*.vue" | wc -l` = 16), all under `web_src/js/components/`. Zero use `<script setup>` (`grep -l "<script setup" web_src/js/components/*.vue | wc -l` = 0). All 16 `.vue` files use `export default {`. Sample: `web_src/js/components/RepoRecentCommits.vue` defines `props: { locale }`, fetches via `GET`, renders chart.js; `web_src/js/components/DashboardRepoList.vue` exports `initDashboardRepoList`.
**Exceptions**: None — do not introduce `<script setup>`, Composition API, or a global Vue app mount. If a UI does not need reactivity, write plain DOM JS in a `features/*.js` file instead.

---

## 5. CSS — Tailwind `tw-` utility classes in templates; component CSS in `web_src/css/modules/`

**Rule**: For new styling prefer Tailwind utility classes prefixed with `tw-` (the project's tailwind config sets `prefix: 'tw-'` and `important: true` at `tailwind.config.js:27-28`). For structural/component CSS that cannot be expressed as utilities, add a scoped file under `web_src/css/modules/<component>.css` and import it from `web_src/css/index.css`. Do **not** invent a new class prefix and do **not** pollute global selectors.
**Why**: Gitea runs Tailwind alongside Fomantic UI; the `tw-` prefix and `important: true` let utility classes override Fomantic without specificity wars. There is **no `GT-` prefix** in this codebase (`grep -rln "GT-" web_src/ templates/ | grep -v CLAUDE.md` returns nothing). The Fomantic overrides live as `web_src/css/modules/{button,card,input,menu,...}.css` (imported at `web_src/css/index.css:5-30`) so each overridden component is one file.
**Frequency**: dominant in templates — `grep -rn "tw-" templates/ | wc -l` returns 1030 occurrences across templates. CSS files in `web_src/css/modules/` (~25 files) are Fomantic overrides; `web_src/css/features/` (~10 files: `heatmap.css`, `gitgraph.css`, `codeeditor.css`) holds feature-specific styles.
**Exceptions**: BEM is not used and should not be introduced. Existing non-prefixed classes (Fomantic's `ui button`, `ui menu`, etc.) remain valid — use them as-is rather than re-styling. Avoid `!important` in component CSS unless overriding Fomantic; rely on Tailwind's `important: true` instead.

---

## 6. Fomantic UI — bundled, component-list declared in `semantic.json`, overridden not forked

**Rule**: Use the bundled Fomantic UI 2.8.7 from `web_src/fomantic/`. Enable/disable components by editing the `components` array in `web_src/fomantic/semantic.json`. Override Fomantic styling in `web_src/css/modules/<component>.css`, never by editing Fomantic source under `web_src/fomantic/build/`.
**Why**: Fomantic is built once into `web_src/fomantic/build/semantic.{js,css}` and then loaded as part of the `index` webpack entry (`webpack.config.js:80-83`). The `semantic.json` `components` array (currently `api`, `dimmer`, `dropdown`, `form`, `modal`, `search`, `tab`, ...) is the source of truth for what ships; editing build output gets clobbered on the next `make svg`/frontend build.
**Frequency**: universal. `web_src/fomantic/package.json` pins `fomantic-ui: 2.8.7`. Override pattern: `web_src/css/modules/modal.css`, `web_src/css/modules/dropdown.css`, etc. Fomantic is loaded at `webpack.config.js:80` (`web_src/fomantic/build/semantic.js`).
**Exceptions**: Do not fork Fomantic. If a component is genuinely missing, add it to `semantic.json` and re-run the Fomantic build (`cd web_src/fomantic && npm run build` per its scripts).

---

## 7. Async/HTTP — use `GET/POST/PATCH/PUT/DELETE` from `web_src/js/modules/fetch.js`

**Rule**: All browser-side HTTP must go through the wrappers in `web_src/js/modules/fetch.js` (`request`, `GET`, `POST`, `PATCH`, `PUT`, `DELETE`). Never call raw `fetch()` without the wrapper, and never use `XMLHttpRequest`.
**Why**: The wrapper (`web_src/js/modules/fetch.js:11-41`) auto-sets the CSRF token header `x-csrf-token` for non-safe methods (POST/PATCH/PUT/DELETE) by reading `window.config.csrfToken`, auto-serializes object/array `data` to JSON with the correct `content-type`, and forwards `FormData`/`URLSearchParams` untouched. Bypassing it means hand-rolling CSRF on every call site and silently breaking under session expiry.
**Frequency**: dominant. `grep -rln "XMLHttpRequest" web_src/ | grep -v CLAUDE.md` returns nothing (verified). Import sites: `web_src/js/features/repo-migrate.js:2`, `web_src/js/features/repo-projects.js:4` (`POST, DELETE, PUT`), `web_src/js/features/notification.js:2`, `web_src/js/components/RepoRecentCommits.vue` (`GET`). Test coverage exists at `web_src/js/modules/fetch.test.js`.
**Exceptions**: The only sanctioned raw `fetch` is inside `fetch.js` itself. Web Workers and Service Workers in `eventsource.sharedworker.js` may use lower-level APIs if they are not doing CSRF-protected writes.

---

## 8. jQuery is global, legacy, and should not grow

**Rule**: jQuery is exposed globally (`window.$ = window.jQuery` at `web_src/js/jquery.js:3`) and loaded first in the `index` entry (`webpack.config.js:79`). New code should not introduce new jQuery usage; prefer vanilla DOM APIs (`document.querySelector`, `addEventListener`) and the `init*` pattern.
**Why**: jQuery exists for Fomantic UI and for legacy initializers; new feature code in `web_src/js/features/` is largely jQuery-free. Expanding jQuery usage deepens a dependency the project is not otherwise modernizing in the 1.22 line.
**Frequency**: legacy prevalence, declining in new code. The global assignment is at `web_src/js/jquery.js:3` with an explicit `eslint-disable no-jquery/variable-pattern` comment.
**Exceptions**: Fomantic UI component calls (`$('.ui.dropdown').dropdown(...)`) require jQuery and are fine inside Fomantic glue code under `web_src/js/modules/fomantic.js`.

---

## 9. SVG icons — import via `web_src/js/svg.js`, do not inline `<svg>` in JS

**Rule**: Add SVG files under `public/assets/img/svg/` (e.g. `octicon-*.svg`, `gitea-*.svg`) and reference them through the `SvgIcon` component exported from `web_src/js/svg.js`. Do not paste raw `<svg>` markup into JS strings or templates when an icon already exists.
**Why**: `web_src/js/svg.js` imports each icon as a URL and exposes a Vue `<SvgIcon name="...">` component plus helpers; this keeps the icon set greppable, deduplicated, and consistent with the octicon/gitea naming scheme. The `make svg` Make target regenerates the asset bundle from the source SVGs.
**Frequency**: dominant. `web_src/js/svg.js` imports ~100+ icons (verified: file is ~100+ lines of `import octiconXxx from '../../public/assets/img/svg/octicon-xxx.svg'`). `SvgIcon` is used across components (e.g. `web_src/js/components/RepoRecentCommits.vue`).
**Exceptions**: One-off decorative SVGs that are not icons (logos, illustrations) may live in templates directly.
