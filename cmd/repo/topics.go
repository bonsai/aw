package repo

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/bonsai/aw.tui/repos"
)

func init() {
	topicsCmd.AddCommand(topicsBuildCmd, topicsApplyCmd)
	rootCmd.AddCommand(topicsCmd)
}

var topicsCmd = &cobra.Command{
	Use:   "topics",
	Short: "repos topics メタデータ作成・適用 workflow",
}

// topicsBuildCmd fetches the owner's repositories, merges the curated
// repos.tags.tsv inventory, and emits per-repo topics metadata JSON.
var topicsBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build repos topics metadata (curated + derived) as JSON",
	RunE: runTopicsBuild,
	Example: "  repo topics build --output repos.topics.json\n" +
		"  repo topics build --tsv data/repos.tags.tsv --output -",
}

var topicsApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply topics metadata to GitHub (PUT /repos/{o}/{r}/topics)",
	RunE: runTopicsApply,
	Example: "  repo topics apply -i repos.topics.json --apply\n" +
		"  repo topics apply -i repos.topics.json            # dry-run",
}

func init() {
	topicsBuildCmd.Flags().String("tsv", "data/repos.tags.tsv", "path to curated repos.tags.tsv")
	topicsBuildCmd.Flags().String("output", "repos.topics.json", "output file ('-' for stdout)")
	topicsBuildCmd.Flags().Int("limit", 0, "max repos to fetch (0 = all)")

	topicsApplyCmd.Flags().StringP("input", "i", "repos.topics.json", "topics metadata JSON input")
	topicsApplyCmd.Flags().Bool("apply", false, "actually PUT topics (default: dry-run)")
	topicsApplyCmd.Flags().Bool("verbose", true, "log each applied repo")
}

func runTopicsBuild(cmd *cobra.Command, _ []string) error {
	tsv, _ := cmd.Flags().GetString("tsv")
	outFlag, _ := cmd.Flags().GetString("output")
	limit, _ := cmd.Flags().GetInt("limit")
	o := owner()

	curated, err := repos.LoadCurated(tsv)
	if err != nil {
		// Curated inventory is optional: continue derived-only when absent.
		curated = map[string]repos.CuratedTag{}
	}
	fmt.Fprintln(os.Stderr, fmt.Sprintf("curated tags: %d (tsv=%s)", len(curated), tsv))

	rs, err := repos.FetchAll(o, limit)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, fmt.Sprintf("fetched repos: %d (owner=%s)", len(rs), o))

	set := repos.RepoTopicSet{
		Owner:     o,
		Generated: time.Now().UTC().Format(time.RFC3339),
		Count:     0,
		Repos:     repos.BuildTopics(rs, curated),
	}
	set.Count = len(set.Repos)

	buf, err := json.MarshalIndent(set, "", "  ")
	if err != nil {
		return err
	}
	if outFlag == "-" {
		_, err = os.Stdout.Write(append(buf, '\n'))
		return err
	}
	if err := os.WriteFile(outFlag, append(buf, '\n'), 0o644); err != nil {
		return err
	}
	cur, der, empty := 0, 0, 0
	for _, r := range set.Repos {
		switch r.Source {
		case "curated":
			cur++
		case "derived":
			der++
		default:
			empty++
		}
	}
	fmt.Fprintln(os.Stderr, fmt.Sprintf("wrote %s (curated=%d derived=%d empty=%d)", outFlag, cur, der, empty))
	return nil
}

func runTopicsApply(cmd *cobra.Command, _ []string) error {
	inFlag, _ := cmd.Flags().GetString("input")
	doApply, _ := cmd.Flags().GetBool("apply")
	verbose, _ := cmd.Flags().GetBool("verbose")
	o := owner()
	tok := token()

	buf, err := os.ReadFile(inFlag)
	if err != nil {
		return err
	}
	var set repos.RepoTopicSet
	if err := json.Unmarshal(buf, &set); err != nil {
		return err
	}
	changed := 0
	for _, r := range set.Repos {
		if len(r.Topics) == 0 {
			continue
		}
		if !doApply {
			fmt.Printf("[dry-run] %s/%s -> %v\n", o, r.Name, r.Topics)
			changed++
			continue
		}
		got, err := repos.ApplyTopics(o, r.Name, r.Topics, tok)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("[error] %s: %v", r.Name, err))
			continue
		}
		changed++
		if verbose {
			fmt.Printf("[ok] %s/%s -> %v\n", o, r.Name, got)
		}
	}
	fmt.Fprintln(os.Stderr, fmt.Sprintf("%d repos (mode=%s)", changed, map[bool]string{true: "apply", false: "dry-run"}[doApply]))
	return nil
}