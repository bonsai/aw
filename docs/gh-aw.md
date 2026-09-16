# gh-aw (GitHub Agentic Workflows) メモ

> ひとこと理解: **マークダウンの指示書(`.md`)を、CLI(`gh aw`)で起動できるクラウド(Actions上のAIエージェント)ワークフローにする仕組み。**

## 仕組み

1. `.github/workflows/〇〇.md` — YAML frontmatter(設定)+ 本文(日本語でAIにやらせたいタスク)を書く
2. `gh aw compile` — 通常の GitHub Actions の `.lock.yml`(実行形式)に変換・検証
3. `gh aw run` — GitHub-hosted runner 上で AI エージェントがタスクを実行
4. `gh aw status / logs / audit` — 実行監視・ログ・コスト分析

- 実行の実体は通常の Actions なので、Pages 配信などの既存CIと同じ土俵で動く。
- 書き込みは安全に: エージェントは**読み取りのみ**。変更は `safe-outputs`(PR/Issue/コメント作成)経由で提案 → 人間がマージ。
  直接 commit させないのが基本設計(権限を `contents: read` 等に絞る)。
- 対象変更ファイルは `safe-outputs.create-pull-request.allowed-files` で縛れる(例: `services/rakugo/data/rakugo/**`)。

## フロントマターの代表キー

```yaml
---
description: ...
intent: 最終成果の一文(実装でなく目的)
on:
  schedule: [{ cron: "..." }]   # or workflow_dispatch / issues 等
permissions: { contents: read, issues: read, pull-requests: read }
tools: { github: { mode: gh-proxy }, ... }
network: { allowed: [ドメイン…] }
safe-outputs:
  create-pull-request: { allowed-files: [...] }
---
```

## 費用（無料範囲）

| 項目 | 公開リポジトリ | private | 備考 |
|---|---|---|---|
| Actions 時間(compute) | **無料(無制限)** | Free 2,000分/月 | 標準 GitHub-hosted runner のみ |
| AI 推論(エンジン) | **有料の可能性** | 同左 | Copilot の月間 **AI Credits(AIC)** から消費。1 AIC = $0.01 |
| ガードレール(既定) | — | — | 1,000 AIC/run・5,000 AIC/日(≈$50)・agent 20分タイムアウト |

- Copilot Pro($10/月) = **base 1,000 + flex 500 = 1,500 AIC/月**。週1回程度の小型ワークフローなら誤差範囲の消費。
- Copilot Free でも AI Credits の割り当てがあり、限定的に動く(モデルは制限)。
- AI 推論コストは**エンジン別**で GitHub とは別枠のチャージになる可能性(PAT ベースだとトークン所有者の Copilot 権利に紐づく)。
- 上限: フロントマターで `max-ai-credits` / `max-daily-ai-credits` / `timeout-minutes` / `max-turns` を明示できる。

## 認証が無いと動かない（最初の壁）

- エンジンは既定 **copilot**。実行には
  - repository secret `COPILOT_GITHUB_TOKEN`(fine-grained PAT。Account permissions → Copilot Requests: Read) **または**
  - `permissions: copilot-requests: write`(organizationの集中請求設定前提)
  が必要。未設定だと run は失敗する(実例あり)。
- Pi を使う場合: `model:` のプロバイダ接頭辞で認証を選択(Copilot 既定 / `ANTHROPIC_API_KEY` / `CODEX_API_KEY`)。

## Bonsai での使いどころ

- yose-db: `.github/workflows/rakugo-zenza-update.md` — 落語家の前座数(zenza_count)の更新を週次で実装済み。
  「協会名簿を取得 → lists/*.json 更新 → go run ./cmd/zenza-count → go test → PR 提案(または noop)」。
- 原則: 読み取り専用 + safe-outputs + 人間レビュー。自動マージは現状やらない。