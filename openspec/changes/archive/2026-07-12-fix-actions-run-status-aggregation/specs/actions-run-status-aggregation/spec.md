## ADDED Requirements

### Requirement: Workflow run status uses every job status
The system SHALL aggregate an Actions workflow run status from every job status,
including success, failure, cancelled, skipped, waiting, running, and blocked.

#### Scenario: Every job is skipped
- **GIVEN** every job in an Actions workflow run has status skipped
- **WHEN** the run status is aggregated
- **THEN** the workflow run status is skipped

#### Scenario: Cancelled job is terminal
- **GIVEN** an Actions workflow run has at least one cancelled job and no waiting, running, or blocked jobs
- **WHEN** the run status is aggregated
- **THEN** the workflow run status is cancelled

#### Scenario: Blocked job is pending
- **GIVEN** an Actions workflow run has at least one blocked job and no running or waiting jobs
- **WHEN** the run status is aggregated
- **THEN** the workflow run status is blocked

#### Scenario: Successful jobs with skipped jobs remain successful
- **GIVEN** an Actions workflow run has only success and skipped jobs and at least one success job
- **WHEN** the run status is aggregated
- **THEN** the workflow run status is success

### Requirement: Actions run list exposes filterable statuses
The Actions run list SHALL expose filter entries for success, failure,
cancelled, skipped, waiting, running, and blocked statuses.

#### Scenario: User filters blocked runs
- **GIVEN** the repository Actions page is rendered
- **WHEN** the status filter menu is populated
- **THEN** blocked is available as a selectable status filter

#### Scenario: User filters skipped or cancelled runs
- **GIVEN** the repository Actions page is rendered
- **WHEN** the status filter menu is populated
- **THEN** skipped and cancelled are available as selectable status filters
