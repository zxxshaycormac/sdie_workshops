## 1. Persistence Query

- [x] 1.1 Recheck tag creation, release creation, push synchronization, mirror synchronization, and relevant legacy handling; record the `CreateNewTag` existing-tag gap and update the design to query authoritative `tag_name` without a migration.
- [x] 1.2 Add a tag-name keyword field to `repo_model.FindReleasesOptions` and apply the cross-database case-insensitive substring helper to `tag_name` without changing exact `TagNames` filtering or default ordering.
- [x] 1.3 Add focused `models/repo` tests proving partial, case-insensitive tag-name matching and proving non-tag fields do not affect the query.

## 2. Tags Page Behavior

- [x] 2.1 Update `repo.TagsList` to trim `q`, pass it to the release query, load rows and their filtered count together, expose `Keyword`, and paginate with the filtered count while leaving `NumTags` unchanged.
- [x] 2.2 Update `templates/repo/tag/list.tmpl` to reuse the shared search combo, retain the search control after zero matches, render the existing no-results message, and preserve all existing tag-row permission checks and actions.
- [x] 2.3 Add the English `search.tag_kind` locale source entry without editing Crowdin-managed non-English locale files or generated template assets.
- [x] 2.4 Extend the Tags-page integration tests to cover the unchanged unfiltered list, partial and case-insensitive matches, keyword echo, zero matches, repository-wide count compatibility, and filtered pagination retaining `q`.

## 3. Verification

- [x] 3.1 Run `./tools/harness/verify.sh go ./models/repo/...` with project Go 1.22.12 and `CGO_ENABLED=1`; the focused package passed.
- [x] 3.2 Run the focused SQLite `TestViewTagsList`: the Make target was blocked by missing Git LFS and `.git`, while the equivalent directly compiled integration binary passed the target test.
- [x] 3.3 Run `./tools/harness/verify.sh spec` and `openspec validate --all --strict --no-interactive`; both strict validations and `openspec doctor` passed.
- [x] 3.4 Confirm no generated source outputs, public API/Swagger contracts, JavaScript assets, migrations, feeds, or unrelated source files changed; full backend, frontend, non-SQLite database, and E2E suites were not run because the implementation stayed within the planned surface.
