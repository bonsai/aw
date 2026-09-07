# 🌱 repo2agent

**Repository → Structure → Agent Workflow**

`repo2agent` is the **midwife of agents**: it reads a repository's existing codebase, discovers its structure, proposes semantic clusters and agents, and generates agentic workflows after approval.

The goal is not to generate YAML blindly. The goal is to turn repository evidence into an intermediate ontology and workforce proposal.

## Core transformation

```text
Repository
   ↓
Observe
   ↓
Understand
   ↓
Structure
   ↓
Cluster
   ↓
Propose domain / capability
   ↓
Propose agents
   ↓
Propose workflows
   ↓
Approve
   ↓
Generate
   ↓
Compile
   ↓
Execute
```

### repo2cluster

Repository clustering belongs in AW because it is **discovery**, not ontology declaration.

```text
GitHub
   ↓
 bonsai/repos
   ↓
name / description / topics / language / README
   ↓
semantic features
   ↓
vectorization
   ↓
KMeans / clustering
   ↓
cluster profile
   ↓
proposed domain
   ↓
ontology validation
   ↓
ecosystem/domains/*.yaml
```

The reusable experiment is defined in [`skills/repo2cluster.yaml`](skills/repo2cluster.yaml).

## Boundary

```text
observed   → GitHub metadata, topics, repository text
inferred   → similarity, clusters, dominant topics
proposed   → domains, capabilities, synapses, agents
validated  → reviewed structure
 declared  → ecosystem ontology
```

**AW discovers. Ecosystem declares.**

Clustering must never silently rewrite the ontology. Its output is evidence-backed proposal data that can be reviewed and then promoted into `bonsai/ecosystem`.

## The Midwife Model

`repo2agent` is neither the CEO nor the worker.

It acts as a **midwife**:

- discovers capabilities already present in repositories
- discovers semantic communities across repositories
- identifies missing responsibilities
- proposes domains and capabilities
- proposes suitable agents
- proposes workflows connecting those agents
- records evidence and confidence
- generates workflow source after approval
- hands execution to `gh-aw` / GitHub Actions

```text
.company       = organization / policy
AW             = observer / clusterer / agent midwife / compiler
agent          = specialist
workflow       = orchestration
repos/         = repository observations
ecosystem/     = declared semantic world
GitHub Actions = execution
```

## Agent Proposal Schema

`schema.json` describes the intermediate model between repository evidence and generated workflows:

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

This intermediate layer is essential. AW should not perform a direct `README → YAML` conversion.

```text
README ─┐
Code   ─┤
Issues ─┤→ Understanding → Structure → Proposal
Tests  ─┤                         │
Config ─┘                         ├→ Domain / Cluster
                                  ├→ Agent
                                  └→ Workflow
                                         ↓
                                      Approval
                                         ↓
                                   Agentic Workflow
```

## GitHub Agentic Workflows

GitHub Agentic Workflows use Markdown workflow source with YAML frontmatter. The source is compiled by `gh aw compile` into a machine-ready `.lock.yml` workflow.

```text
AW
 ↓
.github/workflows/<agent-or-workflow>.md
 ↓
gh aw compile
 ↓
.github/workflows/<agent-or-workflow>.lock.yml
 ↓
GitHub execution
```

## Design Principles

### 1. Evidence before agents
Every proposed agent or domain should be traceable to repository evidence.

### 2. Discovery before declaration
Clustering is an inference mechanism. It does not become ontology automatically.

### 3. Ontology before workflow
Understand the repository and its semantic neighborhood before deciding how to automate it.

### 4. Proposal before execution
Agent and workflow creation passes through a review/policy gate.

### 5. No hard-coded departments
The repository evidence determines candidate structure. Similar repositories may produce different workforces.

### 6. Agents are replaceable
The durable assets are ontology, evidence, and workflow contracts.

### 7. Execution is downstream
AW discovers and generates. `gh-aw` / GitHub Actions executes.

## Architecture

```text
                  bonsai.company
                        │
                  organization
                        │
                        ▼
                       AW
                repo2agent / repo2cluster
                        │
          ┌─────────────┼─────────────┐
          ▼             ▼             ▼
       repos/        ontology      evidence
          │             │             │
          └─────────────┼─────────────┘
                        ▼
                semantic structure
                        │
          ┌─────────────┼─────────────┐
          ▼             ▼             ▼
       clusters     capabilities    agents
          │             │             │
          └─────────────┼─────────────┘
                        ▼
                workflow candidates
                        │
                        ▼
                  approval gate
                        │
                        ▼
                 gh aw compile
                        │
                        ▼
                  GitHub execution
                        │
                        ▼
                    evidence
                        │
                        └────→ semantic review
```

## Relation to the Bonsai ecosystem

```text
repos
  ↓ observation
AW
  ↓ discovery / proposal
ecosystem
  ↓ semantic declaration
ontology
  ↓ meaning
synapse
  ↓ relationships
matrix
  ↓ computation
BQML
  ↓ learning
AW
  ↓ reorganize / regenerate
```

This makes AW the **self-organization engine at the boundary between observed repositories and declared organizational structure**.

## Roadmap

- [ ] Repository inventory and evidence extraction
- [ ] README / code / issue understanding
- [ ] `repo2cluster` implementation
- [ ] cluster profiling and confidence scoring
- [ ] proposed domain generation
- [ ] repository ontology extraction
- [ ] agent candidate generation
- [ ] workflow candidate generation
- [ ] human approval gate
- [ ] `.github/workflows/*.md` generation
- [ ] `gh aw compile` integration
- [ ] workflow validation and safety checks
- [ ] CLI/API interface
- [ ] MCP interface
- [ ] BQML-based repository clustering and workforce analysis

## One-line definition

> **AW observes repositories, discovers their semantic structure, proposes the workforce, and turns approved structure into executable agentic workflows.**
