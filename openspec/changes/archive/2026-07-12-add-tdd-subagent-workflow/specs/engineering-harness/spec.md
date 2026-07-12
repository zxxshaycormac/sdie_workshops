## ADDED Requirements

### Requirement: TDD-first implementation workflow

The repository SHALL require behavior and contract changes to follow a
test-first workflow before production code is changed.

#### Scenario: Contributor starts a behavior change

- **GIVEN** a change affects behavior, a public contract, persistence, frontend
  consumption, or a cross-layer workflow
- **WHEN** a contributor begins implementation
- **THEN** the harness directs them to derive focused tests from traced control
  points and OpenSpec scenarios before changing production code

#### Scenario: Subagent testing support is available

- **GIVEN** an independent subagent can be used for the task
- **WHEN** focused tests are needed for a behavior or contract change
- **THEN** the harness requires the subagent to write or review those tests
  before product code is implemented

#### Scenario: Expected failing evidence is practical

- **GIVEN** the focused test can run before implementation in the local
  environment
- **WHEN** the contributor follows the TDD workflow
- **THEN** the completion evidence includes the expected failing command or
  failure reason before implementation and the passing command after
  implementation

#### Scenario: Test-first work is not applicable

- **GIVEN** the change is documentation-only, formatting-only, or mechanical and
  has no behavioral contract
- **WHEN** the contributor reports completion
- **THEN** the report states why test-first work was not applicable and lists
  the proportionate checks that were run
