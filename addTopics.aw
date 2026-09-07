{
  "type": "workflow",
  "name": "addTopics",
  "version": "1.0.0",
  "id": "addTopics.aw",
  "objective": "GitHub repos の topics（メタデータ）を作成・適用する",
  "owner": "bonsai",
  "tools": {
    "aw": {
      "repository": "https://github.com/bonsai/aw.tui",
      "commands": {
        "build": {
          "usage": "repo topics build --output repos.topics.json",
          "description": "GitHub repos 一覧 + data/repos.tags.tsv（curated）を統合し、repo単位の topics メタデータ JSON を生成する",
          "input": {
            "tsv": "data/repos.tags.tsv (curated tags: repo/lang/git/remote/tags, local-only の資産も含む)"
          },
          "output": {
            "file": "repos.topics.json",
            "schema": {
              "owner": "string",
              "generated": "RFC3339",
              "count": "int",
              "repos": [
                {
                  "full_name": "string",
                  "topics": ["string"],
                  "source": "curated|derived|empty"
                }
              ]
            }
          }
        },
        "apply": {
          "usage": "repo topics apply -i repos.topics.json [--apply]",
          "description": "topics メタデータを GitHub へ PUT /repos/{owner}/{repo}/topics で適用。--apply なしは dry-run",
          "gate": "dry-run レビュー後、--apply で本反映（force 操作のため要承認）"
        }
      }
    }
  },
  "data": {
    "repos.tags.tsv": {
      "file": "data/repos.tags.tsv",
      "note": "手メンテ curated インベントリ。~/repos/db/repos.tags.tsv 由来の正本（dataぽい資産として保全）"
    }
  },
  "steps": [
    { "order": 1, "action": "aw repo topics build", "target": "metadata 生成 repos.topics.json" },
    { "order": 2, "action": "review", "target": "dry-run `aw repo topics apply` で差分確認（curated/derived）" },
    { "order": 3, "action": "aw repo topics apply --apply", "target": "GitHub topics へ本反映" },
    { "order": 4, "action": "verify", "target": "gh api repos/bonsai/<repo>/topics で反映確認" }
  ],
  "gate": "step2 の diff レビュー承認後にのみ step3 実行。apply は破壊的（上書き）のため force push 同様の承認ルール適用",
  "files": {
    "go": [
      "repos/topics.go",
      "repos/fetchall.go",
      "cmd/repo/topics.go",
      "cmd/repo/root.go"
    ]
  }
}