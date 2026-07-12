## ADDED Requirements

### Requirement: Frontend API contract guidance

The repository SHALL document how frontend code participates in shared backend
and frontend API contracts.

#### Scenario: Frontend consumes changed API behavior

- **GIVEN** a frontend change depends on an API field, status, error shape,
  permission, ordering, or side effect
- **WHEN** the contributor follows the frontend API contract guidance
- **THEN** they can identify the backend source of truth, the frontend consumer,
  the generated output boundary, and the required backend and frontend
  verification evidence
