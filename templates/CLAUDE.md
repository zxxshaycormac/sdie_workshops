# `templates/` — Go HTML Template Conventions

`templates/` holds 486 `.tmpl` files rendered by Go's `html/template` engine via `ctx.HTML(status, tplName)`. This file goes deeper than `design.md` Section 6 on template-side patterns; the cross-cutting i18n rules (key naming, placeholder substitution, locale file workflow) and the handler-side render contract (`ctx.HTML` + `ctx.Data`, see `routers/web/CLAUDE.md` Section 4) remain authoritative and are not repeated here.

---

## 1. Page Inheritance — Wrap Every Page in `base/head` + `base/footer`

**Rule**: Every full-page template starts with `{{template "base/head" .}}` on its first non-comment line and ends with `{{template "base/footer" .}}` on its last line. The two calls open and close the `<html><head><body>` shell defined in `templates/base/head.tmpl` and `templates/base/footer.tmpl`; the page body is the literal content between them.
**Why**: `base/head.tmpl` emits the `<!DOCTYPE>`, `<head>` metadata, asset links, the CSRF header, the navbar, and the `page-content` wrapper. `base/footer.tmpl` closes the wrapper, loads `assets/js/index.js`, and emits the closing tags. Skipping either leaves the document half-rendered or breaks asset loading.
**Frequency**: universal for browser-facing pages — 109 templates open with `{{template "base/head" .}}` and 113 close with `{{template "base/footer" .}}` (a few use whitespace variants of the head call).
**Exceptions**: `templates/mail/**/*.tmpl` are standalone HTML for email clients and do NOT use `base/head` or `base/footer` (email renderers reject shared stylesheets and scripts). `templates/api/`, `templates/swagger/`, and `templates/devtest/` fragments are partials or non-browser payloads and are likewise exempt. Use the devtest pattern (`templates/devtest/tmplerr.tmpl`) as the reference layout for a new full page.

---

## 2. User-Visible Text — `{{ctx.Locale.Tr "key.name" args...}}`, Never Hardcoded English

**Rule**: Every string the user reads comes through `{{ctx.Locale.Tr "locale.key" args...}}`. The locale is exposed on the template context as `ctx.Locale` (the same `context.Context` the handler received); there is no `.i18n` field on the data map. Add the new key to `options/locale/locale_en-US.ini` first.
**Why**: Gitea ships 40+ languages; a hardcoded English literal forces every translator to grep for it and breaks non-English UIs. `ctx.Locale.Tr` also accepts positional args for placeholder substitution, so translators can reorder sentences.
**Frequency**: dominant — 353 templates call `ctx.Locale.Tr` or read `ctx.Locale.Lang`; zero templates use a `.i18n.Tr` field. See `design.md` Section 6 for the locale-file workflow and the no-concatenation rule.
**Exceptions**: Programmatic identifiers never shown to users (`id` attributes, CSS class names, ARIA labels that mirror a visible translated string, data attributes). Inline SVG `<title>` and `alt` text that duplicate a translated string already on the page may stay literal to avoid double-lookups, but prefer `{{ctx.Locale.Tr "..."}}` for clarity.

---

## 3. Template Helpers — Register in `modules/templates/helper.go`, Use by Registered Name

**Rule**: Use the helpers registered in `NewFuncMap()` (`modules/templates/helper.go:30`). The common ones: `{{AssetUrlPrefix}}` (static asset root, do not hardcode `/assets`), `{{DateTime "full|short|long" .CreatedUnix}}` (date formatting — not `.DateFmt`), `{{ctx.AvatarUtils.Avatar .User 48 "class"}}` (avatar rendering — a method on a context-scoped util, not a free function `.Avatar`), `{{svg "octicon-x"}}` (inline SVG icon, 250 call sites), and the `Render*` family (`{{RenderCodeBlock $desc}}`, `{{RenderMarkdownToHtml $body}}`, `{{RenderLabel .}}`). For safe HTML output use `{{SafeHTML ...}}` / `{{HTMLFormat ...}}` explicitly — `html/template` escapes everything else by default.
**Why**: Helpers are the only sanctioned bridge between Go data and rendered HTML; they centralize escaping, asset versioning, avatar URL generation, and markdown rendering. The registered names are the contract — `grep` and tooling depend on them. Note the spec's `.AssetUrl`, `.Render`, `.Avatar`, `.DateFmt` are aspirational names; the real registered names are `AssetUrlPrefix`, `Render*`, `ctx.AvatarUtils.Avatar`, `DateTime`.
**Frequency**: universal; every non-trivial template calls at least one helper. New helpers that touch `models/*` go in `modules/templates/util_*.go`; those files are an accepted structural concession (model imports contained there, do not spread the pattern).
**Exceptions**: Pure presentation logic that needs no Go data may stay as Go template primitives (`{{if}}`, `{{range}}`, `{{with}}`). Do not invent a helper to paper over a missing locale key.

---

## 4. Naming — Template Path Mirrors Route, File Extension Is Always `.tmpl`

**Rule**: The template path under `templates/` mirrors the URL route and the handler package. `templates/repo/home.tmpl` is rendered by `routers/web/repo/view.go` (constant `tplRepoHome base.TplName = "repo/home"`); `templates/user/dashboard/dashboard.tmpl` by `routers/web/user/dashboard.go`; `templates/admin/cron.tmpl` by `routers/web/admin/cron.go`. Shared cross-feature fragments go under `templates/shared/<feature>/` (e.g. `templates/shared/user/authorlink.tmpl`, `templates/shared/combomarkdowneditor.tmpl`); feature-private fragments stay in their own subdir (e.g. `templates/repo/diff/section_split.tmpl`, `templates/repo/issue/view_content/add_reaction.tmpl`).
**Why**: Mirroring the route makes it trivial to find the template from a URL and vice versa; the `tplXxx base.TplName` constant in the handler is the single source of truth for the path string. `.tmpl` is the only extension — 486 of 486 files use it.
**Frequency**: universal. The constant lives at the top of the handler file (e.g. `routers/web/repo/view.go:62-69` declares `tplRepoHome`, `tplRepoViewList`, `tplRepoEMPTY`, `tplMigrating`); the path string in the constant drops the `.tmpl` suffix.
**Exceptions**: `templates/base/`, `templates/custom/`, `templates/mail/`, `templates/api/`, `templates/swagger/`, `templates/status/` are not tied to a single route — they are infrastructure shared by many routes.

---

## 5. Partials — `{{template "name" .}}`, Pass Data With `dict` When the Signature Differs

**Rule**: Include a shared fragment with `{{template "dotted.path.no.ext" .}}` — the path is the template file path with no leading `templates/` and no `.tmpl`. When the partial needs a different data shape than the caller has, build it inline with `{{template "repo/diff/section_split" dict "file" . "root" $}}` (the `dict` helper constructs a map). Partials live in `templates/base/` for global fragments (`alert.tmpl`, `paginate.tmpl`, `modal_actions_confirm.tmpl`, `disable_form_autofill.tmpl`), in `templates/shared/` for cross-feature fragments, or in a per-feature subdir for private fragments.
**Why**: `{{template "x" .}}` is Go's native include — it re-renders with the passed dot. `dict` is the established escape hatch when a partial expects named keys (e.g. `dict "ActionURL" ... "Attachments" ...`) rather than the caller's full page context; 143 templates use this pattern. Naming partials by file path makes them greppable and avoids a separate partial registry.
**Frequency**: common; partials are heavily reused (`base/alert`, `base/paginate`, `shared/user/authorlink`, `repo/branch_dropdown`, `repo/commit_statuses`). Sub-features like diff rendering fan out across many partials under `templates/repo/diff/`.
**Exceptions**: A partial that genuinely needs the whole page context still takes `.` and never a `dict` — do not wrap the caller's dot in a one-key dict for symmetry.

---

## 6. Custom Overrides — Drop-in Files Under `custom/templates/` Shadow `templates/`

**Rule**: Templates are loaded through a layered asset filesystem (`modules/templates/base.go:AssetFS`) that layers `CustomAssets()` (on disk at `custom/templates/`) on top of `BuiltinAssets()` (compiled-in bindata). A file placed at `custom/templates/repo/home.tmpl` overrides the built-in `templates/repo/home.tmpl` for the entire instance without recompiling. The built-in tree ships placeholder files under `templates/custom/` (`header.tmpl`, `footer.tmpl`, `body_inner_pre.tmpl`, `body_outer_pre.tmpl`, `extra_links.tmpl`, `extra_tabs.tmpl`) specifically as override points that `base/head.tmpl` and `base/footer.tmpl` already include via `{{template "custom/header" .}}` etc.
**Why**: Site operators customize branding, analytics, nav links, and footer content without forking. The layered filesystem means the built-in file is the fallback and the custom file is the override — no code change required. The `templates/custom/*.tmpl` placeholders are empty by design so they render nothing unless overridden.
**Frequency**: universal as a mechanism; per-instance in practice. New branding hooks should be added as new empty placeholders under `templates/custom/` and included from `base/head.tmpl` or `base/footer.tmpl` rather than scattering `{{if}}` blocks across feature templates.
**Exceptions**: Do not override `base/head.tmpl` or `base/footer.tmpl` themselves in `custom/` — it breaks asset loading and CSRF wiring on upgrades. Use the `custom/header`, `custom/footer`, `custom/body_*` hooks instead.

---

## 7. Frontend Boundary — No Inline `<script>`, Tailwind `tw-` Classes for Styling

**Rule**: Do not add inline `<script>` tags inside templates; put JS in a module under `web_src/js/`, export an `init*` function, and call it from `web_src/js/index.js`. Use Tailwind utility classes with the `tw-` prefix (e.g. `tw-flex tw-mt-2`) for layout; component-level CSS lives in `web_src/css/modules/`. SVG icons go through `{{svg "octicon-x"}}`, not raw `<svg>` markup.
**Why**: Inline scripts break webpack bundling, escape the CSRF/fetch wrapper, and break CSP-friendly deployment. The `tw-` prefix lets Tailwind utilities override Fomantic UI without specificity wars. These rules apply equally to JS modules under `web_src/js/` and are restated here only as a pointer.
**Frequency**: dominant. Only 4 legacy templates contain inline `<script>` (`base/head_script.tmpl`, `repo/clone_script.tmpl`, `repo/diff/box.tmpl`, `shared/combomarkdowneditor.tmpl`); each has a structural reason and a TODO. `grep -rn "tw-" templates/` returns ~1030 occurrences.
**Exceptions**: `templates/base/head_script.tmpl` must stay inline because it bootstraps `window.config` before the JS bundle loads. New exceptions require a TODO and a structural justification (e.g. avoiding FOUC).
