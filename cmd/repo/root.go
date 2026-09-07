package repo

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd is the "aw repo" CLI root. Sub-commands (select, topics, …) live in
// this package and are registered through init().
var rootCmd = &cobra.Command{
	Use:   "repo",
	Short: "Repository operations for the aw.tui agentic workflow",
}

// Execute runs the repo command tree with the given arguments.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("owner", "o", "", "GitHub owner / organization (default from GITHUB_OWNER env)")
	rootCmd.PersistentFlags().String("token", "", "GitHub token (default from GITHUB_TOKEN env)")
	viper.BindPFlag("repo.owner", rootCmd.PersistentFlags().Lookup("owner"))
	viper.BindPFlag("repo.token", rootCmd.PersistentFlags().Lookup("token"))
}

// owner returns the effective owner: flag → GITHUB_OWNER → "bonsai".
func owner() string {
	if o := viper.GetString("repo.owner"); o != "" {
		return o
	}
	if o := os.Getenv("GITHUB_OWNER"); o != "" {
		return o
	}
	return "bonsai"
}

// token returns the effective GitHub token.
func token() string {
	if t := viper.GetString("repo.token"); t != "" {
		return t
	}
	return os.Getenv("GITHUB_TOKEN")
}