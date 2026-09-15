# AW Compiler

AW Compiler converts an approved `*.aw.yaml` intent into an executable `*.wf.yaml` plan.

## Principle

It is **not** `AW YAML -> LLM -> trust output`.

```text
AW YAML
  -> schema validation
  -> prompt construction
  -> examples / cases retrieval
  -> LLM planner
  -> structured WF candidate
  -> WF schema validation
  -> repair loop (LLM)
  -> approved WF
```

The LLM is used for the part that requires planning and interpretation. Schemas remain deterministic boundaries. Examples are supplied as few-shot planning cases so the model learns the repository's workflow conventions instead of inventing an arbitrary format.

## Inputs

- `schema/aw.schema.json`
- `schema/wf.schema.json`
- `*.aw.yaml`
- approved examples under `examples/`
- optional domain/function contracts from `skills/` and `generators/`

## Output

- `*.wf.yaml`

## Example

```bash
aw compile examples/research.aw.yaml -o examples/research.wf.yaml
```

The first implementation should keep the LLM provider behind an adapter interface. This allows LM Studio, OpenAI-compatible endpoints, or another planner to be used without changing AW semantics.
