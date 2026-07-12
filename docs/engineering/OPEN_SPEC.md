# OpenSpec In This Repository

OpenSpec is the change-contract layer of the engineering harness. It keeps the
reason, required behavior, technical decisions, implementation tasks, and final
specification connected to the code change.

OpenSpec does not statically scan or understand the code by itself. In this
repository, a "scan" means the agent reads the real code paths and feeds that
evidence into OpenSpec artifacts. The CLI creates, orders, validates, and
archives those artifacts.

## Why It Helps

- Requirements become reviewable before implementation consumes time.
- Scenarios turn vague intent into potential test cases.
- Designs name real control points and expose cross-layer omissions.
- Tasks keep long AI sessions bounded and resumable.
- Strict validation catches malformed or disconnected artifacts.
- Archive updates main specs, leaving durable current-system knowledge.
- A later contributor can distinguish why a decision exists from how it was
  coded.

OpenSpec does not replace tests, code review, Git, observability, or judgment.
It makes their intended use explicit.

## Artifact Model

```text
proposal.md  -> why, scope, non-goals, capabilities
     |
     +-> specs/<capability>/spec.md -> normative behavior and scenarios
     |
     +-> design.md                  -> control points and decisions
             |
             v
          tasks.md                  -> ordered implementation and checks
             |
             v
          implementation -> validate -> archive -> openspec/specs/
```

The default `spec-driven` schema enforces this dependency order. Use
`openspec status --change <name>` whenever the next valid action is unclear.

## Recommended Entry Points

After restarting Codex so local skills are discovered, conversational commands
are the easiest interface:

- `/opsx:explore <idea>`: investigate code and clarify an idea without coding.
- `/opsx:propose <idea>`: create proposal, specs, design, and tasks.
- `/opsx:apply <change>`: implement pending tasks using all artifacts as context.
- `/opsx:update <change>`: continue or revise change artifacts.
- `/opsx:archive <change>`: finish the lifecycle and update main specs.

The generated skills live under `.codex/skills/`. The CLI remains useful for
inspection, validation, automation, and learning what the skills are doing.

## First Change Workflow

### 1. Explore the real code

Use explore mode when the request is still ambiguous or when the control point
is unknown:

```text
/opsx:explore add repository-level release approvals
```

Explore mode is read-only for product code. A good exploration identifies the
route, handler, service, model, contract, UI, tests, assumptions, and risks.

### 2. Create a complete proposal

```text
/opsx:propose add repository-level release approvals
```

Use a kebab-case name such as `add-release-approvals`. A strong proposal states
why the change matters, exact scope, non-goals, new or modified capabilities,
and affected compatibility surfaces.

Equivalent low-level starting point:

```sh
openspec new change add-release-approvals
openspec status --change add-release-approvals
openspec instructions proposal --change add-release-approvals --json
```

Normally let `/opsx:propose` produce all required artifacts. Hand-authoring is
useful for learning but must follow the template returned by `instructions`.

### 3. Review requirements before implementation

A requirement is normative and testable:

```md
### Requirement: Approval blocks release publication
The system SHALL prevent publication until the configured approval count is met.

#### Scenario: Approval count is insufficient
- **GIVEN** a repository requires two approvals and has one
- **WHEN** an authorized user attempts to publish the release
- **THEN** publication is rejected without creating release artifacts
```

Check especially:

- Does each requirement use SHALL or MUST?
- Does every requirement have a `#### Scenario`?
- Do scenarios reflect MECE branch analysis: mutually exclusive categories,
  important edge paths, and explicit exclusions?
- Are authorization, failure, compatibility, and side effects covered?
- Does the design name concrete paths from `PROJECT_MAP.md`?
- Are non-goals clear enough to prevent scope creep?

### 4. Apply tasks

```text
/opsx:apply add-release-approvals
```

The agent reads all context files, implements one bounded task at a time, runs
adjacent checks, and marks `- [ ]` as `- [x]`. If code evidence invalidates an
assumption, update the artifacts before broadening implementation.

Inspect progress at any time:

```sh
openspec status --change add-release-approvals
openspec instructions apply --change add-release-approvals --json
```

### 5. Validate

```sh
./tools/harness/verify.sh spec
openspec validate add-release-approvals --strict --no-interactive
```

OpenSpec validation checks artifact structure and spec semantics. It does not
prove the Go or frontend implementation works, so run the project checks selected
in the design and `VERIFICATION.md`.

### 6. Archive

After every task and required check is complete:

```text
/opsx:archive add-release-approvals
```

Or with the CLI:

```sh
openspec archive add-release-approvals
```

Archive moves the change to a dated history directory and merges delta specs
into `openspec/specs/`. Do not archive merely to hide incomplete work.

## Daily Commands

```sh
openspec list                 # active changes
openspec list --specs         # main capabilities
openspec show <name>          # render a change or spec
openspec status --change <id> # artifact readiness
openspec context              # current agent brief
openspec doctor               # relationship health
openspec validate --all --strict --no-interactive
openspec view                 # interactive dashboard
```

## Scope Heuristic

Use OpenSpec for behavior, contracts, data, security, async work, architecture,
or cross-layer changes. Skip it for a truly non-behavioral typo or formatting
edit, and say why. When uncertain, a small well-written proposal is cheaper than
discovering a missing contract during implementation.

## Common Failure Modes

- Writing a proposal before tracing code: the design names the wrong control
  point. Explore first.
- Treating specs as implementation notes: requirements should describe visible
  behavior, not function names.
- Using `MODIFIED` for a partial requirement: archive can discard omitted
  behavior. Copy the complete existing requirement block, then edit it.
- Omitting failure scenarios: happy-path code passes while authorization,
  retries, or rollback remain undefined.
- Leaving tasks vague: "implement backend" is not independently verifiable.
- Archiving before checks: main specs then claim behavior the code may not have.
- Assuming OpenSpec replaces Git: it records intent, not source diffs or history.

## Bootstrap Example

The first lifecycle is `establish-engineering-harness`. Its proposal, delta
spec, design, tasks, implementation, strict validation, and archive provide a
local example aligned with this repository rather than a generic demo.
