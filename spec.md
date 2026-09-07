# BONSAI Agentic Workflow 再設計

## 1. 目的

本設計は **Agentic Workflow に特化** する。

- `soubi` がすべてをオーケストラする
- `aw` は 7 つの道具のうちの 1 つ
- ユーザーは `soubi` だけを開けばよい

| 問い | 答え |
|------|------|
| What do we have? | `soubi` |
| What are we doing? | `soubi`（aw run タブ） |
| What is it? | `ontology/` |
| What is its state? | `BQML` |
| What should we do? | `AW / Agent` |
| Do it. | `pi` |

## 2. 全体像

```text
                    BONSAI
                      │
          ┌───────────┴───────────┐
          │                       │
        MODEL                   VIEW
          │                       │
   ┌──────┼──────┐                │
   │      │      │                │
 repos   BQML  ontology         soubi
   │      │                  （オーケストラ）
   │    state                     │
   │      │                       │ event
   └──────┴───────────────────────┘
                                  │
                              CONTROLLER
                                  │
                              AW / Agent
                                  │
                                  ▼
                                  pi
```

## 3. 7 つの道具

Agent は以下の 7 つの道具を装備できる。

```text
bonsai.equipment
│
├── aw           ← Agentic Workflow 実行エンジン
├── mcp          ← 外部接続（GitHub, BigQuery など）
├── model        ← 推論モデル
├── skill        ← 能力・関数
├── dataset      ← データ・知識
├── repository   ← コード資産
└── agent        ← 他の Agent（マルチエージェント）
```

`soubi` はこの 7 道具を閲覧・装備・実行する唯一のインターフェース。

## 4. オントロジー（命名空間）

```text
bonsai
│
├── ontology/           ← 定義と契約
│   ├── organization.yaml
│   ├── repository.yaml
│   ├── asset.yaml
│   ├── agent.yaml
│   ├── equipment.yaml
│   ├── state.yaml
│   └── action.yaml
│
├── repos/              ← 生のリポジトリデータ
├── github-observatory/ ← GitHub → BigQuery → BQML
├── yaml-as-agent/      ← Agent 装備定義
└── soubi/              ← 唯一の TUI・オーケストレータ
```

### Namespace

```text
bonsai.organization
bonsai.repository
bonsai.asset
bonsai.agent
bonsai.equipment
bonsai.state
bonsai.action
```

## 5. soubi の責務

`soubi` は **指揮者（conductor）** である。

### 5.1 閲覧

7 道具を見る。

```text
soubi
├── aw
├── mcp
├── model
├── skill
├── dataset
├── repository
└── agent
```

### 5.2 装備

Agent に道具を装備する。

```text
agent: repo-coordinator

equipment:
  repositories:
    - github-observatory
    - repos
  skills:
    - repo-discovery
    - bqml-analysis
  models:
    - gemma
  mcps:
    - github
  aw:
    - repo-health-check
```

### 5.3 実行監視

aw run を起動・監視・承認する。

```text
soubi › aw
├── plans
├── workflows
├── runs
├── approvals
└── results
```

### 5.4 イベント発行

`soubi` はイベントを発行するだけ。判断は Agent / pi 側。

```text
VIEW_EQUIPMENT    {genre, id}
EQUIP             {agent, equipment}
UNEQUIP           {agent, equipment}
INSPECT           {genre, id}
CALCULATE         {target}
NEW_RUN           {aw, input}
APPROVE           {run_id, step}
REJECT            {run_id, step}
STOP_RUN          {run_id}
```

## 6. MVC 分離

### Model

```text
bonsai/
├── repos/              ← 生のリポジトリデータ
├── github-observatory/ ← GitHub → BigQuery → BQML
├── yaml-as-agent/      ← Agent 装備定義
└── repo.yaml           ← 組織全体の状態・契約
```

計算フロー:

```text
GitHub
  ↓
Raw
  ↓
BigQuery
  ↓
BQML
  ├── Health
  ├── Activity
  ├── Score
  ├── Cluster
  └── Relationship
          ↓
    Organization State
          ↓
    soubi
```

### View

`soubi` のみ。

```text
soubi
├── equipment view     ← 7 道具を見る
├── state view         ← score/health/cluster
└── aw view            ← run / approval / result
```

### Controller

```text
soubi
  ↓ event
AW / Agent
  ↓
pi
```

## 7. 装備のオントロジー

```text
Organization
     │
     │ owns
     ▼
Asset
     │
     │ capability
     ▼
Agent
     │
     │ equipped-with
     ├──── aw
     ├──── mcp
     ├──── model
     ├──── skill
     ├──── dataset
     ├──── repository
     └──── agent
```

## 8. State は別 namespace

```text
bonsai.state
│
├── activity
├── health
├── score
├── cluster
├── relation
├── execution
└── assignment
```

## 9. Action（ドラクエ風操作）

```text
みる      → view
せってい  → configure
そうび    → equip
はずす    → unequip
かくにん  → inspect
けいさん  → calculate
クエスト  → run （aw 実行）
承認      → approve
```

## 10. 移行

### 現状

```text
aw.tui/
├── cmd/aw-tui/main.go   ← View + Controller が混在
├── graph/graph.go       ← Model
└── repos/repos.go       ← Model（GitHub クライアント）
```

### 移行後

```text
bonsai/soubi/
├── cmd/soubi/main.go    ← 唯一の TUI
├── internal/
│   ├── tui/             ← Bubble Tea View
│   ├── events/          ← イベント発行
│   └── aw/              ← aw run / approval / result 表示
└── go.mod

bonsai/github-observatory/
├── cmd/observatory/
├── internal/
│   ├── github/          ← repos.FetchRecent 相当
│   ├── bigquery/
│   └── bqml/
└── go.mod

bonsai/ontology/
├── organization.yaml
├── repository.yaml
├── asset.yaml
├── agent.yaml
├── equipment.yaml
├── state.yaml
└── action.yaml

bonsai/yaml-as-agent/
├── agents/
│   └── repo-coordinator.yaml
└── schemas/
    └── equipment.json
```

## 11. 非機能

- `soubi` は Bubble Tea + lipgloss
- Model は Go の独立モジュール
- ontology は YAML + JSON Schema で検証
- BQML は BigQuery 上で実行

## 12. 完了基準

- [ ] `ontology/` に 7 つの YAML スキーマが存在
- [ ] `soubi` が 7 道具を閲覧できる
- [ ] `soubi` から Agent への装備・解除ができる
- [ ] `soubi` から aw run を起動・承認・停止できる
- [ ] `github-observatory` が Model として独立
- [ ] `yaml-as-agent` で Agent 装備を定義
- [ ] `soubi` からのイベントが Agent / pi に到達
