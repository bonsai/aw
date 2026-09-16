# BQML CSV

CSV is the simple tabular seed format for BQML inside AW.

## Purpose

The CSV layer provides observable records that BQML can use to discover:

- clusters
- similarities
- relationships
- migration patterns
- candidate semantic structures

BQML produces evidence and inferred relationships. AW consumes those results and routes them into proposals and generators.

```text
CSV
 ↓
BQML
 ↓
clusters / relations
 ↓
AW Router
 ↓
Action
 ↓
Event
 ↓
Status
 ↓
Generate / Validate / Execute
```

## Minimal schema

| Field | Meaning |
|---|---|
| `id` | Record identifier |
| `source` | Source repository or entity |
| `target` | Target repository or entity |
| `entity` | Semantic entity represented by the record |
| `relation` | Relationship between source and target |
| `action` | Action associated with the relationship |
| `event` | Observable event produced by the action |
| `status` | Current lifecycle state |

## Example

```csv
id,source,target,entity,relation,action,event,status
MIG-001,bonsai/wiki,bonsai/issues,Issue,migrates_to,inventory,migration_planned,planned
MIG-002,bonsai/wiki,bonsai/issues,Status,migrates_to,inventory,migration_planned,planned
MIG-003,bonsai/wiki,bonsai/issues,Action,migrates_to,inventory,migration_planned,planned
MIG-004,bonsai/wiki,bonsai/issues,Event,migrates_to,inventory,migration_planned,planned
MIG-005,bonsai/wiki,bonsai/solve,Workflow,references,deduplicate,reference_created,planned
MIG-006,bonsai/wiki,bonsai/wiki,Skill,retains,retain,retained,planned
MIG-007,bonsai/wiki,bonsai/wiki,Agent,retains,retain,retained,planned
MIG-008,bonsai/wiki,bonsai/wiki,Tool,retains,retain,retained,planned
```

## BQML boundary

CSV is data, not ontology.

BQML may infer clusters and relationships from CSV records, but those inferences are not automatically declared as Bonsai ontology. AW uses the resulting evidence to propose and generate artifacts; validation remains the promotion boundary.

```text
CSV = observed data
BQML = learning / inference
AW = routing / proposal / generation
Validation = promotion boundary
Ontology = declared meaning
```

## Design rule

> Keep CSV boring. Put semantics in columns, let BQML discover structure, and let AW decide what to do with the evidence.
