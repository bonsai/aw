# Issue Solving Ecosystem

`bonsai/aw`, `bonsai/solve`, `bonsai/issues`, and `bonsai/goal` form one ecosystem for taking work from a Goal to a verified Solve.

## Repositories

| Repository | Role | Question |
|---|---|---|
| `bonsai/goal` | Goal | Where do we want to go? |
| `bonsai/issues` | Issue registry | What needs to be solved? |
| `bonsai/solve` | Solving ontology/protocol | How is an Issue solved? |
| `bonsai/aw` | Agent Workflow runtime | How do Agents execute the solving process? |

`issues/issues.md` is the central human-readable issue inventory. `solve.yaml` is the machine-readable solving workflow contract.

## End-to-end flow

```text
Goal
  ↓
Breakdown
  ↓
Issue
  ↓
Observe current Status
  ↓
Plan Actions
  ↓
Design
  ├── Document
  └── UML / Graph
  ↓
Task
  ↓
TODO
  ↓
Ticket
  ↓
AW Router
  ↓
Agent + Skill
  ↓
Action
  ↓
Event / Evidence
  ↓
Progress / Status
  ↓
Verify
  ↓
Solve
```

## Responsibility boundary

### `bonsai/goal`

Defines desired outcomes and goals. A Goal is the starting point for decomposition, not an executable ticket.

### `bonsai/issues`

Owns the list of problems/work that must be solved. Issues are the durable units of work and may be decomposed into Tasks.

```text
Goal
 ├── Issue A
 ├── Issue B
 └── Issue C
```

### `bonsai/solve`

Defines the semantics of solving. Its ontology contains `Issue`, `Action`, `Event`, `Status`, and `Skill`, with the status machine:

```text
open
 → understood
 → planned
 → executing
 → implemented
 → verified
 → solved
```

Work units are distinguished as:

```text
Issue  = problem to solve
Task   = coherent work required
TODO   = concrete next activity
Ticket = executable work package
```

### `bonsai/aw`

Executes the workflow around those definitions. AW observes the current repository/issue state, routes work, discovers required capabilities, proposes and generates artifacts, validates them, executes Tickets through Agents, and evolves from resulting evidence.

## Work decomposition

```text
Goal
  │
  └── Issue
        │
        ├── Task
        │     ├── TODO
        │     └── TODO
        │
        └── Task
              └── TODO
                    ↓ formalize
                  Ticket
```

The distinction is important:

- **Issue** is the problem.
- **Task** is the work.
- **TODO** is the next concrete activity.
- **Ticket** is the executable package.

A TODO becomes a Ticket when enough context, target, action, input, expected output, and acceptance criteria exist for an Agent to execute it.

## `solve.yaml`

`solve.yaml` is the machine-readable contract connecting the issue-solving model to AW.

Conceptually:

```yaml
solve:
  goal: goal reference
  issues: issue references
  observe: status and evidence
  plan: actions
  design:
    document: true
    graph: true
    uml: true
  tasks: task decomposition
  todos: next activities
  tickets: agent-executable work
  verify: acceptance criteria
  terminal_status: solved
```

The exact schema may evolve, but the boundary remains stable: `solve.yaml` describes the solving process; AW performs it.

## AW execution loop

AW applies its existing factory loop to the solving lifecycle:

```text
Observe
  → Route
  → Discover
  → Propose
  → Generate
  → Validate
  → Execute
  → Evolve
```

For Issue solving this becomes:

```text
Observe Issue
  → Route to solving strategy
  → Discover Tasks / Skills / dependencies
  → Propose plan and design
  → Generate TODOs / Tickets / artifacts
  → Validate plan and outputs
  → Execute Tickets with Agents
  → Evolve from Events and verification evidence
```

## Graph model

The ecosystem can be represented as a graph:

```text
Goal
 ↓
Issue ── has_status ──→ Status
 ↓
Task
 ↓
TODO
 ↓
Ticket ── assigned_to ──→ Agent
 ↓                         ↓
Action ── uses_skill ──→ Skill
 ↓
Event
 ↓
Status
```

This graph makes progress observable. A Ticket attempt is not equivalent to a solved Issue; the system must observe Events and verify the resulting state.

## BQML boundary

BQML can operate over the graph to discover clusters, similarities, dependencies, recurring work patterns, and bottlenecks.

```text
Observed data
     ↓
    BQML
     ↓
clusters / relations / patterns
     ↓
    AW
     ↓
proposals / routing / generation
     ↓
GitHub execution + evidence
```

BQML discovers structure; it does not become the source of truth for Issue status. `solve.yaml`, GitHub Issues, workflow events, and verification evidence remain explicit execution artifacts.

## Core principle

> Goal defines the destination. Issue defines the problem. Task structures the work. TODO identifies the next activity. Ticket packages executable work. AW routes and executes it. Solve verifies the result.

The ecosystem is therefore a **Goal → Issue → Work → Agent → Evidence → Solve** system rather than a simple issue tracker.
