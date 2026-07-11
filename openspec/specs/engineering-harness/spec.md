# engineering-harness Specification

## Purpose

Define the repository-level operating model that makes human and AI changes code-grounded, bounded by explicit contracts, and supported by proportionate verification evidence.

## Requirements

### Requirement: Repository operating instructions

The repository SHALL provide a root-level instruction surface that tells human and AI contributors how to orient, plan, implement, verify, and report a change.

#### Scenario: Contributor starts a non-trivial change

- **GIVEN** a change affects product behavior, a public contract, persistence, or more than one architectural layer
- **WHEN** a contributor begins work
- **THEN** the instructions require an OpenSpec change and code-grounded control-point tracing before implementation

#### Scenario: Contributor handles a trivial change

- **GIVEN** a change only corrects wording, comments, or formatting without changing behavior
- **WHEN** a contributor begins work
- **THEN** the instructions permit a direct edit with proportionate verification and an explicit statement that no behavioral contract changed

### Requirement: Code-grounded project navigation

The repository SHALL document its startup path, request flow, source boundaries, contract locations, generated outputs, runtime-only paths, and test locations using concrete repository paths.

#### Scenario: Contributor investigates an API change

- **GIVEN** an API behavior change is proposed
- **WHEN** the contributor consults the project map
- **THEN** the map directs them through route registration, handler, service, model or module, public structs, Swagger, and relevant tests

#### Scenario: Contributor encounters generated or local files

- **GIVEN** a candidate edit is under a generated-output or runtime-only path
- **WHEN** the contributor consults the project map
- **THEN** the map identifies the source-of-truth path or requires explicit runtime-diagnosis scope before editing

### Requirement: Proportionate executable verification

The repository SHALL provide an executable verification entry point that supports strict OpenSpec validation and focused backend or frontend checks without falsely implying full-suite coverage.

#### Scenario: OpenSpec artifacts are changed

- **GIVEN** OpenSpec configuration, specs, or change artifacts were edited
- **WHEN** the specification verification mode runs
- **THEN** all OpenSpec items are validated strictly and relationship health is checked

#### Scenario: Focused code is changed

- **GIVEN** one or more Go packages or frontend tests are affected
- **WHEN** the contributor selects the corresponding focused verification mode
- **THEN** the entry point runs only the explicitly supplied scope and reports the invoked command

#### Scenario: Verification scope is omitted

- **GIVEN** a focused backend or frontend mode requires targets
- **WHEN** no targets are supplied
- **THEN** the entry point exits with usage guidance rather than running an ambiguous broad suite

### Requirement: Safety and evidence boundaries

The repository SHALL distinguish source from secrets and mutable local state, and SHALL require final reports to state checks run, checks not run, and material workspace limitations.

#### Scenario: Runtime state is outside task scope

- **GIVEN** a task does not explicitly concern local runtime diagnosis
- **WHEN** an agent works in the repository
- **THEN** it does not inspect or edit local configuration, data, logs, dependencies, or built binaries

#### Scenario: Git metadata is unavailable

- **GIVEN** the workspace has no Git metadata
- **WHEN** a contributor reports completion
- **THEN** the report states that history, blame, and Git-based diff or rollback verification were unavailable

### Requirement: OpenSpec lifecycle guidance

The repository SHALL document how to explore, propose, apply, validate, and archive changes, including when each action is appropriate.

#### Scenario: First-time OpenSpec user follows a change

- **GIVEN** a user has an idea but no active change
- **WHEN** they follow the documented workflow
- **THEN** they can create complete artifacts, implement tasks, strictly validate the result, and archive it into the main specification
