## ADDED Requirements

### Requirement: MEMC branch analysis

The repository SHALL require behavior and contract changes to use MEMC-style
branch analysis before implementation and test selection.

#### Scenario: Contributor designs behavior coverage

- **GIVEN** a change affects behavior, public contracts, persistence, frontend
  consumption, asynchronous work, or cross-layer orchestration
- **WHEN** the contributor prepares OpenSpec scenarios or TDD tests
- **THEN** they identify mutually exclusive and maximally complete relevant
  branches for states, inputs, permissions, side effects, failures, and outputs

#### Scenario: Branch is outside the contract

- **GIVEN** a possible branch is impossible, unreachable, or outside the stated
  OpenSpec non-goals
- **WHEN** the contributor applies MEMC branch analysis
- **THEN** they record the exclusion instead of silently ignoring it or writing
  speculative tests

#### Scenario: Completion evidence is reported

- **GIVEN** MEMC branch analysis shaped test selection
- **WHEN** the contributor reports completion
- **THEN** the report identifies covered important branches and any intentionally
  skipped branches with residual risk
