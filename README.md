# 🌱 repo2agent

**Repository → Agent Workflow**

`repo2agent` is the **midwife of agents**: it reads a repository's existing codebase and discovers what agents and agentic workflows the repository should have.

The goal is not to generate YAML blindly. The goal is to turn the repository's **README, source code, data, tests, configuration, issues, and existing workflows** into a repository ontology, then propose and generate the agents/workflows that naturally belong to that repository.

## Concept

```text
GitHub Repository
       │
       ├── README
       ├── source code
       ├── data
       ├── tests
       ├── config
       ├── Issues
       └── existing workflows
       │
       ▼
  Repository Understanding
       │
       ▼
  Repository Ontology
       │
       ├── capabilities
       ├── responsibilities
       ├── data flows
       ├── constraints
       └── evidence
       │
       ▼
   Agent Candidates
       │
       ▼
 Workflow Candidates
       │
       ▼
 Human Review / Policy Gate
       │
       ▼
.github/workflows/*.md
       │
       │  YAML frontmatter + natural-language instructions
       ▼
   gh aw compile
       │
       ▼
*.lock.yml
       │
       ▼
 GitHub Actions / gh-aw
```

## Why `repo2agent`?

A repository already contains the strongest evidence of what its agents should do.

Instead of starting with:

```text
"What agent should I build?"
```

start with:

```text
"What does this repository already know how to do?"
```

`repo2agent` extracts that knowledge and turns it into an **agent workforce proposal**.

## The Midwife Model

`repo2agent` is neither the CEO nor the worker.

It does not decide the business direction, and it does not replace the agents it creates.

It acts as a **midwife**:

- discovers capabilities already present in the repository
- identifies missing responsibilities
- proposes suitable agents
- proposes workflows connecting those agents
- records the evidence behind each proposal
- generates GitHub Agentic Workflow source after approval
- hands execution to `gh-aw` / GitHub Actions

```text
.company       = organization / policy
repo2agent     = agent midwife / compiler
agent          = specialist
workflow       = orchestration
repos/         = repository data
GitHub Actions = execution
```

## Repository → Agent

The fundamental transformation is:

```text
Repository
   ↓
Understand
   ↓
Model
   ↓
Discover
   ↓
Propose
   ↓
Approve
   ↓
Generate
   ↓
Compile
   ↓
Execute
```

### Example

For a quiz repository, `repo2agent` might discover:

```text
Repository capabilities
├── question JSON
├── quiz UI
├── validation
├── historical source data
└── tests

Agent candidates
├── quiz-generator
├── quiz-validator
├── fact-checker
└── data-maintainer

Workflow candidates
├── generate-quiz
├── validate-quiz
├── check-facts
└── update-data
```

The important point is that these agents are **derived from repository evidence**, rather than being hard-coded as a universal organization chart.

## Agent Proposal Schema

`schema.json` is intended to describe the intermediate model between a repository and generated workflows.

```json
{
  "repository": {},
  "ontology": {},
  "agent_candidates": [],
  "workflow_candidates": [],
  "evidence": [],
  "risk": {},
  "approval": {},
  "generated_workflows": []
}
```

This intermediate layer is important. `repo2agent` should **not** perform a direct `README → YAML` conversion.

```text
README ─┐
Code   ─┤
Issues ─┤→ Ontology → Agent Proposal → Workflow Proposal
Tests  ─┤                                      ↓
Config ─┘                                  Approval
                                             ↓
                                      Agentic Workflow
```

## GitHub Agentic Workflows

GitHub Agentic Workflows use Markdown workflow source with YAML frontmatter. The source is compiled by `gh aw compile` into a machine-ready `.lock.yml` workflow.

Therefore the output of `repo2agent` is conceptually:

```text
repo2agent
    ↓
.github/workflows/<agent-or-workflow>.md
    ↓
gh aw compile
    ↓
.github/workflows/<agent-or-workflow>.lock.yml
```

`repo2agent` is consequently a **compiler front-end / workforce designer**, while `gh-aw` is the execution-oriented workflow compiler/runtime layer.

## Design Principles

### 1. Evidence before agents

Every proposed agent should be traceable to repository evidence.

### 2. Ontology before workflow

Understand the repository before deciding how to automate it.

### 3. Proposal before execution

Agent creation should pass through a review/policy gate rather than silently modifying the repository.

### 4. No hard-coded departments

The repository determines its candidate workforce. Similar repositories may produce different organizational structures.

### 5. Agents are replaceable

The durable asset is the repository ontology and workflow contract, not a particular model or agent implementation.

### 6. Execution is downstream

`repo2agent` discovers and generates. GitHub Actions / `gh-aw` executes.

## Architecture

```text
                 bonsai.company
                       │
                 organization
                       │
                       ▼
                  repo2agent
                       │
          ┌────────────┼────────────┐
          ▼            ▼            ▼
       repos/       ontology     evidence
          │            │            │
          └────────────┼────────────┘
                       ▼
                agent candidates
                       │
                       ▼
              workflow candidates
                       │
                       ▼
                 approval gate
                       │
                       ▼
                .aw.md / .md
                       │
                       ▼
                 gh aw compile
                       │
                       ▼
                  .lock.yml
                       │
                       ▼
                GitHub execution
```

## Relation to `aw`

`aw` is the foundation for this transformation.

The original `aw.tui` concept focused on observing repositories and discovering communities. That observation layer remains useful, but the strategic direction is broader:

> **Observe repositories → understand them → discover agents → generate workflows.**

The repository itself becomes the starting point of agent organization.

## Roadmap

- [ ] Repository inventory and evidence extraction
- [ ] README / code / issue understanding
- [ ] Repository ontology schema
- [ ] Agent candidate generation
- [ ] Workflow candidate generation
- [ ] Evidence and confidence scoring
- [ ] Human approval gate
- [ ] `.github/workflows/*.md` generation
- [ ] `gh aw compile` integration
- [ ] Workflow validation and safety checks
- [ ] CLI/API interface
- [ ] MCP interface
- [ ] BQML-based repository clustering and workforce analysis

## One-line definition

> **repo2agent turns a repository's existing knowledge into a proposed agent workforce and executable agentic workflows.**
