# 🌱 AW — Agentic Workflow アーキテクチャ

**AW は「親・仕様・入口」だけを持つ。実行の実体は 3 つのサブシステムに分割された。**

```text
bonsai/aw            = AW architecture / contracts / ecosystem entrypoint
bonsai/aw-router     = routing（どこへ送るか）
bonsai/aw-generator  = generation（何を生成するか）
bonsai/aw-runner     = gh aw execution（どう起動するか）
```

## 責務の境界

| リポジトリ | 責務 | 問い |
|---|---|---|
| **bonsai/aw** | 親・仕様・入口・契約 | What is AW? |
| **bonsai/aw-router** | Route / Select | Where does it go? |
| **bonsai/aw-generator** | Generate / Propose | What is generated? |
| **bonsai/aw-runner** | Dispatch / Execute | How is it launched? |

## 全体像

```text
Issue
  ↓
Agent
  ↓
Solve
  ↓
aw-router
  ↓
aw-runner ──→ gh aw
       │
       ↓
   workflow ──→ gh wf（bonsai/workflow）
       │
       ↓
    Action
       ↓
    State
       ↓
    Goal
```

`gh aw` は aw-runner 側、`gh wf` は `bonsai/workflow` 側という綺麗な境界になっている。

## このリポジトリが持つもの

```text
aw/
├── README.md        # ← エントリポイント
├── spec.md          # AW 再設計仕様
├── design.md        # 設計メモ
├── ux.md            # ユーザー体験物語
├── schema.json      # オントロジー契約
├── docs/            # issue-solving-ecosystem 等
├── examples/        # research.aw.yaml / research.wf.yaml 例
└── issues/          # 設計イシュー
```

## サブシステム

- **bonsai/aw-router** — `router/router.yaml` で観測/要求を capabilities へルーティング。
- **bonsai/aw-generator** — `generators/*.yaml` + スキルで proposal → artifact を生成。
- **bonsai/aw-runner** — 実行境界。`gh aw` integration・TUI・API 設計・evolve ループ。

## 設計原則

1. **Evidence before generation** — すべての proposal は観測にトレース可能
2. **Discovery before declaration** — 推論は黙って ontology にならない
3. **Router before workflow** — セマンティックコンテキストで送り先を決める
4. **Generation is first-class** — Repo / Domain / Agent / Workflow / Synapse は生成される
5. **Proposal before creation** — 新規 repo / ontology 宣言は承認が必要
6. **Execution is downstream** — GitHub Actions / `gh aw` が実行する
7. **The loop is open** — 生成物は次の観測になる

## One-line

> **AW は観測を routing / generation / execution の 3 層に分けて、Bonsai の repo・agent・workflow グラフを進化させる自己組織化ファクトリー。**