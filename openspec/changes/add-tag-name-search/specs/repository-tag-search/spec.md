## ADDED Requirements

### Requirement: Tags page provides tag-name search
The system SHALL provide a keyword search control on the repository Tags page and SHALL submit the keyword through the optional `q` query parameter on the existing Tags page route.

#### Scenario: User opens a repository with tags
- **GIVEN** the user can read a repository that contains tags
- **WHEN** the user opens the repository Tags page
- **THEN** the page displays a search input for tag names

#### Scenario: User submits an empty keyword
- **GIVEN** the user is on the repository Tags page
- **WHEN** the user submits an empty or whitespace-only search keyword
- **THEN** the system displays the normal unfiltered tag list

### Requirement: Search filters by tag name
The system SHALL match tags whose names contain the normalized keyword, without regard to letter case, and SHALL NOT match other release or commit fields.

#### Scenario: Keyword partially matches multiple tag names
- **GIVEN** a repository contains tags named `v1.0`, `v1.1`, and `v2.0`
- **WHEN** the user searches for `v1`
- **THEN** the page displays `v1.0` and `v1.1`
- **AND** the page does not display `v2.0`

#### Scenario: Keyword differs only by letter case
- **GIVEN** a repository contains a tag whose name includes uppercase or lowercase letters
- **WHEN** the user searches with different letter casing
- **THEN** the tag is included in the results

#### Scenario: Keyword occurs only outside the tag name
- **GIVEN** a release title, note, commit message, or commit SHA contains the keyword but its tag name does not
- **WHEN** the user searches for that keyword on the Tags page
- **THEN** that tag is not included solely because another field matches

### Requirement: Filtered results remain navigable
The system SHALL base Tags-page pagination on the number of matching tags and SHALL preserve the normalized keyword in the search input and pagination links.

#### Scenario: Matching tags span multiple pages
- **GIVEN** the matching tag count exceeds the configured Tags-page size
- **WHEN** the user follows a pagination link
- **THEN** the next page contains results for the same keyword
- **AND** the search input retains the keyword

#### Scenario: Search has no matches
- **GIVEN** a repository contains tags
- **WHEN** the user searches for a keyword that matches no tag name
- **THEN** the page keeps the search control visible with the keyword populated
- **AND** the page displays a no-matching-results state
- **AND** the page does not display result pagination

### Requirement: Search preserves existing tag-page behavior
The system MUST preserve the repository-wide tag count, tag ordering, permissions, and tag actions when tag-name search is unused or applied.

#### Scenario: Search result count differs from repository tag count
- **GIVEN** a repository contains more tags than the current keyword matches
- **WHEN** the filtered Tags page is rendered
- **THEN** repository navigation continues to display the total number of repository tags
- **AND** result pagination uses only the matching tag count

#### Scenario: Reader searches tags
- **GIVEN** the user can read the repository but cannot modify its tags
- **WHEN** the user searches on the Tags page
- **THEN** the user can view matching tags
- **AND** the page does not expose tag actions that the user could not access before the change

