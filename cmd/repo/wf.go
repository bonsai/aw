package repo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
	wfCmd.AddCommand(wfJSON2YamlCmd, wfYaml2JSONCmd)
	rootCmd.AddCommand(wfCmd)
}

var wfCmd = &cobra.Command{
	Use:   "wf",
	Short: "Workflow file (.aw) format conversion: JSON ↔ YAML",
}

var wfJSON2YamlCmd = &cobra.Command{
	Use:     "json2yaml",
	Aliases: []string{"json2yamal"},
	Short:   "Convert JSON workflow (.aw) to YAML (human-readable wf)",
	RunE:    runJSON2Yaml,
	Example: "  repo wf json2yaml -i addTopics.aw -o addTopics.yaml",
}

var wfYaml2JSONCmd = &cobra.Command{
	Use:     "yaml2json",
	Aliases: []string{"yamal2json"},
	Short:   "Convert YAML workflow back to canonical JSON",
	RunE:    runYaml2JSON,
	Example: "  repo wf yaml2json -i addTopics.yaml -o addTopics.aw",
}

func init() {
	wfJSON2YamlCmd.Flags().String("input", "addTopics.aw", "JSON workflow input")
	wfJSON2YamlCmd.Flags().String("output", "", "YAML output path (default: input with .yaml ext)")
	wfYaml2JSONCmd.Flags().String("input", "addTopics.yaml", "YAML workflow input")
	wfYaml2JSONCmd.Flags().String("output", "", "JSON output path (default: input with .aw ext)")
}

// normalizeJSON decodes the input into a map so ordering is stable (keys are
// sorted alphabetically by encoding/json -> yaml.v3 preserves node order).
func normalizedMap(buf []byte) (map[string]any, error) {
	var v any
	if err := json.Unmarshal(buf, &v); err != nil {
		return nil, err
	}
	m, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected a JSON object at the top level")
	}
	return m, nil
}

func runJSON2Yaml(cmd *cobra.Command, _ []string) error {
	in, _ := cmd.Flags().GetString("input")
	out, _ := cmd.Flags().GetString("output")
	if out == "" {
		out = in + ".yaml"
	}
	buf, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	m, err := normalizedMap(buf)
	if err != nil {
		return err
	}
	y, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}
	return os.WriteFile(out, y, 0o644)
}

func runYaml2JSON(cmd *cobra.Command, _ []string) error {
	in, _ := cmd.Flags().GetString("input")
	out, _ := cmd.Flags().GetString("output")
	if out == "" {
		if len(in) > len(".yaml") && filepath.Ext(in) == ".yaml" {
			out = in[:len(in)-len(".yaml")] + ".aw"
		} else {
			out = in + ".aw"
		}
	}
	buf, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	var v any
	if err := yaml.Unmarshal(buf, &v); err != nil {
		return err
	}
	j, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}
	return os.WriteFile(out, append(j, '\n'), 0o644)
}