# 🌱 AW — Router + Generator

**Observe → Route → Discover → Propose → Generate → Validate → Execute → Evolve**

AW is the self-organizing factory between observed GitHub repositories and declared organizational structure. It routes evidence to capabilities and generates new domains, repositories, agents, workflows, and synapses.

## Core definition

> **AW routes observations into generators, creating new repositories, domains, agents, workflows, and synapses.**

```text
GitHub / repositories
        ↓
     OBSERVE
        ↓
       AW
   ┌────┴────┐
 ROUTER   GENERATORS
   │          │
   ↓          ↓
repo/cluster  Domain / Repo / Agent / Workflow / Synapse
   │          │
   └────┬─────┘
        ↓
     VALIDATE
        ↓
     EXECUTE
        ↓
     EVIDENCE
        ↓
      EVOLVE
        ↺
```

## Router

`router/router.yaml` defines semantic routing rules. AW decides **where an observation or request should go**, rather than hard-coding one workflow for every repository.

```text
repository observation → repo2cluster
semantic cluster       → cluster2domain
domain candidate       → domain2repo
domain candidate       → domain2synapse
repository capability  → repo2agent
approved workflow       → workflow-generator
```

## Generators

Generators turn proposals into concrete artifacts:

```text
generators/
├── domain.yaml
├── repo.yaml
└── synapse.yaml
```

The target is to make Repo / Domain / Agent / Workflow / Synapse / Ontology all first-class generated artifacts.

A discovered domain can become a new repository proposal:

```text
cluster → proposed domain → validated domain → repo generator → new repository
```

Generated repository structure:

```text
README.md
AGENTS.md
domain.yaml
ontology.yaml
synapse.yaml
agents/
.github/workflows/
```

Repository creation remains an approval boundary.

## Self-organization workflow

`workflows/evolve.md` defines the common loop:

```text
OBSERVE → DISCOVER → PROPOSE → GENERATE → VALIDATE → EVOLVE → OBSERVE
```

The key difference from a simple `repo2agent` compiler is that AW can generate **new semantic structure and new repositories**, not merely workflows inside an existing repository.

## Domain and Synapse discovery

`repo2cluster` discovers semantic communities from repository evidence. `cluster2domain` turns an inferred cluster into a domain proposal. `domain2synapse` proposes weighted relationships, while `domain2repo` prepares a new repository structure.

```text
repos
 ↓
clusters
 ↓
proposed domain ─────→ proposed synapses
 ↓
validated domain
 ↓
new repo
 ↓
new agents / workflows
 ↓
evidence
 ↺
```

**AW discovers and generates. `bonsai/ecosystem` validates and declares.**

## Lifecycle boundary

```text
observed   → GitHub metadata, code, README, issues, workflows
inferred   → clusters, similarity, capabilities
proposed   → domains, synapses, agents, repositories
generated  → YAML, repository manifests, workflow source
validated  → reviewed structure
declared   → ecosystem ontology
```

Generation must never silently promote a proposal to declared ontology.

## Relation to Bonsai

```text
repos       = observation
AW          = route / discover / generate / orchestrate
ecosystem   = declared semantic world
ontology    = meaning
synapse     = relationship
matrix      = computation
BQML        = learning
.company    = organization / policy
GitHub      = execution and evidence
```

BQML can discover clusters and relationships; AW consumes those results and turns them into proposals and generated artifacts. BQML is not the source of truth.

## Repository layout

```text
aw/
├── README.md
├── schema.json
├── router/
├── generators/
├── skills/
└── workflows/
```

## Design principles

1. **Evidence before generation** — every proposal is traceable to evidence.
2. **Discovery before declaration** — inference does not silently become ontology.
3. **Router before workflow** — route according to semantic context.
4. **Generation is first-class** — Repo, Domain, Agent, Workflow, and Synapse can be generated.
5. **Proposal before creation** — new repositories and ontology declarations require validation.
6. **Execution is downstream** — GitHub Actions / `gh-aw` executes approved workflows.
7. **The loop is open** — generated structure becomes new evidence for the next cycle.

## One-line architecture

> **AW is a self-organizing factory that routes observations to generators and evolves the Bonsai repository/agent/workflow graph through evidence-backed proposals.**
