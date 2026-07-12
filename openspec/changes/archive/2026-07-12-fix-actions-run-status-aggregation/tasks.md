## 1. Test-First Coverage

- [x] 1.1 Add focused model tests for Actions run status aggregation branches.
- [x] 1.2 Add focused model test coverage for the Actions run status filter entries.
- [x] 1.3 Run the focused model package test before production changes and record the expected failure.

## 2. Core Implementation

- [x] 2.1 Update `aggregateJobStatus` to handle skipped, cancelled, and blocked statuses explicitly.
- [x] 2.2 Update `GetStatusInfoList` so the Actions list filter exposes all supported run statuses.

## 3. Verification

- [x] 3.1 Run the focused model package test through the verification harness.
- [x] 3.2 Run strict OpenSpec validation for all changes and specs.
- [x] 3.3 Review diffs for generated-output or unrelated changes and document residual risk.
