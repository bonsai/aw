# Journal — aw を gh-aw フレームワークへ転換（2026-09-08）

> aw の設計を「自前 .aw 形式（JSON⇄YAML）+ 自前 CLI/API/MCP」から、
> **公式 gh-aw（GitHub Agentic Workflows）中心** へ向ける方向性を確認した session の記録。
> 本 journal に該当: `docs/journal.md`（dev-recap 規約に従い docs/recap.md と併用可）。

## 状態（次回の出発点）

- **gh-aw CLI を opencode スキルとして取り込んだ。**
  - スキル本体: `~/repos/gh-aw/SKILL.md`（v0.88.2 のフルコマンド、workflow 構成、落とし穴、MCP 連携）
  - 登録: `~/.config/opencode/skills/gh-aw` → `~/repos/gh-aw/`（symlink）
  - 付随メンテ: `.config/opencode/skills` の18個の壊れ symlink を修復（natalie→port-manager 変更、mala/microsoft-foundry 削除）+ 未登録9スキル追加（dev-* / agent-mala / voice-bbs-dev / port-manager）。計27スキル
  - `~/repos`（=bonsai/my-skills index）: `.gitignore` に `gh-aw/` `agent-mala/` 追加、mala 削除、commit `b94018b` push 済
- **理解の確認**: 「aw で issue.md → workflow*.yaml」＝ gh-aw の実行モデルそのもの。
  - `.github/workflows/<id>.md`（本文＝プロンプト、frontmatter＝on/permissions/network/tools/safe-outputs の契約）
  - `gh aw compile <id>` → `<id>.lock.yml`（GitHub Actions）→ merge → `gh aw run <id>`
  - 自前 `.aw`(JSON⇄YAML) は gh-aw の markdown + YAML frontmatter に吸収可能
- ローカル gh aw は v0.88.2 導入済み。エンジン: claude/codex/copilot/gemini/pi。
- **スキル公開**: `bonsai/gh-aw-skill`（https://github.com/bonsai/gh-aw-skill 、PUBLIC、topics: skill/gh-aw/agentic-workflows/cli/github/mcp）。
  既存 `bonsai/gh-aw`（内製 Archimedes実行基盤）と名前衝突のため `-skill` 接尾辞で新規作成。Gym mode セクション追記済。
- SKILL.md は Gym mode（validate→trial→audit/outcomes→forecast サイクル）を含む 155 行。

## 今セッションの完了（証拠）

| 項目 | 証拠 |
|------|------|
| gh-aw スキル作成 | `~/repos/gh-aw/SKILL.md` + `.config/opencode/skills/gh-aw` |
| skills symlink 全修復 | `.config/opencode/skills/` 27エントリ、壊れ0 |
| my-skills index gitignore | commit `b94018b`（push 済） |
| gh-aw 現行テンプレ確認 | `/tmp/opencode/gh-aw-skill/.github/workflows/addTopics.md`（v0.88.2 の frontmatter all-option） |
| MCP 公式文書反映 | `aw#5` に 2026-07-28 仕様（stateless）+ Go SDK + gh aw mcp-server を追記 |

## 決定事項（ブロッカー・保留）

- aw の実装方針: 「自前フレームワーク一択」→ **「gh-aw 公式 + OAuth token」が主線**。aw/CLI/API/MCP は
  inventory 層（topics メタデータ・タスク管理）の薄いラッパーへ縮退。
- 従量課金回避: 既定 `anthropic_api_key`（$1.3〜/回）→ OAuth token に差し替え必須。
  `claude setup-token`（`sk-ant-oat…`）→ Actions secret `CLAUDE_CODE_OAUTH_TOKEN_FOR_AW`。
- 保留: repo topics 全量（gh GraphQL 1451件）の `topics.json` 再生成は中断中（user abort）→ aw#2。
- 保留: `--apply` は明示承認まで dry-run のみ。

## 次の一手（優先順）

1. `addTopics` を gh-aw 形式で作成: `bonsai/aw` の `.github/workflows/addTopics.md`
   （frontmatter: workflow_dispatch / safe-outputs: update-issue, create-pull-request / engine: claude）
2. `gh aw validate addTopics` → `gh aw compile addTopics --engine claude` → PR → OAuth token sed → merge → `gh aw run addTopics`
3. opencode へ `gh aw mcp-server`（stdio）を MCP 登録するかは次回判断
4. aw#2/#3/#5 に対応コメント追記（gh-aw 転換の決定を issue に反映）

## 数量メモ

- skills: `.config/opencode/skills` 27 エントリ（壊れ 0）
- gh aw: v0.88.2
- my-skills index: commit `b94018b`（push 済、working tree = gitignore 2件追加のみ）
## aw-api 決定（2026-09-08）

- **構成**: DB キュー主線 + GCP Cloud Run REST（人間承認済）。
- Spec 草案: `api/design.md` / `api/openapi.yaml` / `api/ddl.sql` / `api/task.schema.json`
- フロー: Client →(IAP)→ Cloud Run aw-api →(enqueue)→ aw_tasks(Cloud SQL) →(claim by scheduler worker)→ `gh aw run` で GitHub Actions へ dispatch → Result を DB へ readback。📀 repos#7 の data 正本化の対象。
- GitHub 認証: GitHub App `bonsai-aw-app`（PAT 非推奨）。コストガード: max_aic + forecast。
- 未決定: DB=Cloud SQL vs Firestore / IAP 投入時期 / 認証単位 / schedule 拡張。→ Spec レビュー待ち。
