package repos

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
)

// RepoTopic is the topics metadata entry produced by "repo topics build".
type RepoTopic struct {
	FullName    string   `json:"full_name"`
	Name        string   `json:"name"`
	Language    string   `json:"language"`
	Description string   `json:"description"`
	Topics      []string `json:"topics"`
	Source      string   `json:"source"` // "curated" | "derived" | "empty"
}

// CuratedTag is one row of the repos.tags.tsv inventory.
type CuratedTag struct {
	Repo   string   `json:"repo"`
	Lang   string   `json:"lang"`
	Git    string   `json:"git"`
	Remote string   `json:"remote"`
	Tags   []string `json:"tags"`
}

// LoadCurated parses the hand-maintained repos.tags.tsv inventory:
//   repo<TAB>lang<TAB>git<TAB>remote<TAB>tags      (tags comma/space separated)
func LoadCurated(path string) (map[string]CuratedTag, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	out := map[string]CuratedTag{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "\t") {
			continue
		}
		cols := strings.Split(line, "\t")
		for len(cols) < 5 {
			cols = append(cols, "")
		}
		key := strings.TrimSpace(cols[0])
		if key == "" {
			continue
		}
		tags := splitTags(cols[4])
		out[key] = CuratedTag{
			Repo:   key,
			Lang:   strings.TrimSpace(cols[1]),
			Git:    strings.TrimSpace(cols[2]),
			Remote: strings.TrimSpace(cols[3]),
			Tags:   tags,
		}
	}
	return out, sc.Err()
}

// stop is a small list of words too generic to become a derived topic.
var stop = map[string]bool{
	"the": true, "a": true, "and": true, "for": true, "with": true, "using": true,
	"via": true, "web": true, "app": true, "lib": true, "tool": true, "agent": true,
	"api": true, "github": true, "repo": true, "repos": true, "to": true, "in": true,
	"of": true, "on": true, "by": true, "an": true, "is": true, "are": true,
}

// deriveTopics computes fallback topics from name + description tokens.
func deriveTopics(r Repo) []string {
	var out []string
	seen := map[string]bool{}
	for _, tk := range Tokens(r) {
		if len(tk) < 2 || len(tk) > 18 || stop[tk] || seen[tk] {
			continue
		}
		seen[tk] = true
		out = append(out, tk)
		if len(out) >= 6 {
			break
		}
	}
	return out
}

// BuildTopics merges curated tags with derived fallbacks for every repo.
func BuildTopics(rs []Repo, curated map[string]CuratedTag) []RepoTopic {
	byName := map[string]Repo{}
	for _, r := range rs {
		byName[r.Name] = r
		byName[strings.ToLower(r.Name)] = r
	}
	// Include curated entries even if the repo is not in the fetched list
	// (covers local-only / data-like inventory rows).
	keys := map[string]bool{}
	for _, r := range rs {
		keys[r.Name] = true
	}
	for c := range curated {
		keys[c] = true
	}
	var names []string
	for n := range keys {
		names = append(names, n)
	}
	sort.Strings(names)

	out := make([]RepoTopic, 0, len(names))
	for _, name := range names {
		r, hasRepo := byName[name]
		if !hasRepo {
			r = byName[strings.ToLower(name)]
		}
		tags, found := curated[name], false
		if tg, ok := curated[name]; ok {
			tags, found = tg, true
		} else if tg, ok := curated[strings.ToLower(name)]; ok {
			tags, found = tg, true
		}
		entry := RepoTopic{
			FullName:    name,
			Name:        name,
			Language:    r.Language,
			Description: r.Description,
		}
		switch {
		case found && len(tags.Tags) > 0:
			entry.Topics = tags.Tags
			entry.Source = "curated"
		case hasRepo:
			entry.Topics = deriveTopics(r)
			entry.Source = "derived"
		default:
			entry.Source = "empty"
		}
		out = append(out, entry)
	}
	return out
}

// RepoTopicSet is the metadata document produced by "repo topics build".
type RepoTopicSet struct {
	Owner    string      `json:"owner"`
	Generated string     `json:"generated"`
	Count    int         `json:"count"`
	Repos    []RepoTopic `json:"repos"`
}

// ApplyTopics adds topics on GitHub for one repo, keeping any existing topics
// (union, not replace). Requires GITHUB_TOKEN.
func ApplyTopics(owner, name string, topics []string, tok string) ([]string, error) {
	if tok == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN required to apply topics")
	}
	endpoint := "https://api.github.com/repos/" + urlPathEscape(owner) + "/" + urlPathEscape(name) + "/topics"
	// 1) read current topics
	get, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	get.Header.Set("Accept", "application/vnd.github+json")
	get.Header.Set("X-GitHub-Api-Version", "2026-03-10")
	if tok != "" {
		get.Header.Set("Authorization", "Bearer "+tok)
	}
	var cur struct {
		Names []string `json:"names"`
	}
	if resp, err := http.DefaultClient.Do(get); err != nil {
		return nil, err
	} else if resp.StatusCode == http.StatusOK {
		_ = json.NewDecoder(resp.Body).Decode(&cur)
		resp.Body.Close()
	}
	// 2) union with new topics
	seen := map[string]bool{}
	names := make([]string, 0, len(cur.Names)+len(topics))
	for _, t := range cur.Names {
		if !seen[t] {
			seen[t] = true
			names = append(names, t)
		}
	}
	for _, t := range topics {
		if !seen[t] {
			seen[t] = true
			names = append(names, t)
		}
	}
	body, _ := json.Marshal(map[string]any{"names": names})
	req, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2026-03-10")
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("topics PUT %s: %s", name, resp.Status)
	}
	var got struct {
		Names []string `json:"names"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		return nil, err
	}
	sort.Strings(got.Names)
	return got.Names, nil
}

func splitTags(v string) []string {
	parts := strings.FieldsFunc(v, func(r rune) bool {
		return r == ',' || r == ';' || r == ' ' || r == '\t'
	})
	var out []string
	seen := map[string]bool{}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" || seen[p] {
			continue
		}
		seen[p] = true
		out = append(out, p)
	}
	return out
}