## ADDED Requirements

### Requirement: Directory roadmap guidance

The repository SHALL provide a concise roadmap document inside every Git-tracked
first-level source directory.

#### Scenario: Contributor enters a first-level directory

- **GIVEN** a contributor opens a Git-tracked first-level directory
- **WHEN** they look for local orientation
- **THEN** the directory contains a `ROADMAP.md` describing its purpose,
  control paths, expected edits, verification signals, and boundaries

#### Scenario: Contributor sees a runtime-only directory

- **GIVEN** a first-level directory is untracked local runtime state,
  dependencies, Git metadata, or generated evidence
- **WHEN** the contributor follows the harness
- **THEN** the root instructions and project map identify it as outside normal
  source edits instead of requiring a local roadmap

### Requirement: Root instructions remain concise

The repository SHALL keep the root operating instructions short enough to scan
and SHALL link to detailed supporting documents for deeper rules.

#### Scenario: Root instructions list directory roadmaps

- **GIVEN** directory-local roadmaps exist
- **WHEN** a contributor reads `AGENTS.md`
- **THEN** it links to the roadmaps and stays below 200 lines

#### Scenario: Detailed API test guidance is needed

- **GIVEN** a contributor changes HTTP API behavior
- **WHEN** they read `AGENTS.md`
- **THEN** it points them to the detailed API testing contract instead of
  embedding the full contract inline
