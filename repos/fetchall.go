package repos

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
)

// ghRepoOut mirrors the JSON shape returned by `gh repo list --json`.
type ghRepoOut struct {
	Name            string      `json:"name"`
	NameWithOwner   string      `json:"nameWithOwner"`
	Description     string      `json:"description"`
	PrimaryLanguage *ghLang     `json:"primaryLanguage"`
	URL             string      `json:"url"`
	DefaultBranch   *ghBranch   `json:"defaultBranchRef"`
	IsFork          bool        `json:"isFork"`
	IsArchived      bool        `json:"isArchived"`
	StargazerCount  int         `json:"stargazerCount"`
	ForkCount       int         `json:"forkCount"`
	UpdatedAt       string      `json:"updatedAt"`
	Topics          *ghTopics   `json:"repositoryTopics"`
}

type ghTopics struct {
	Nodes []struct {
		Topic struct {
			Name string `json:"name"`
		} `json:"topic"`
	} `json:"nodes"`
}

type ghLang struct {
	Name string `json:"name"`
}

type ghBranch struct {
	Name string `json:"name"`
}

// FetchAll returns repositories of the owner, ordered deterministically by
// full name. It prefers the `gh` CLI (GraphQL, like github-observatory) which
// returns the FULL inventory; the REST users/{owner}/repos endpoint silently
// drops hundreds of bonsai repos, so it is only a fallback when gh is absent.
func FetchAll(owner string, limit int) ([]Repo, error) {
	if gh, err := exec.LookPath("gh"); err == nil {
		rs, err := fetchAllGH(gh, owner, limit)
		if err == nil {
			return rs, nil
		}
		fmt.Fprintln(os.Stderr, "fetchAllGH failed, falling back to REST:", err)
	}
	return fetchAllREST(owner, limit)
}

func fetchAllGH(gh, owner string, limit int) ([]Repo, error) {
	args := []string{"repo", "list", owner, "--limit", "5000",
		"--json", "name,nameWithOwner,description,primaryLanguage,url,defaultBranchRef,isFork,isArchived,stargazerCount,forkCount,updatedAt,repositoryTopics"}
	cmd := exec.Command(gh, args...)
	cmd.Env = os.Environ()
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("gh repo list: %w", err)
	}
	var raw []ghRepoOut
	if err := json.Unmarshal(out, &raw); err != nil {
		return nil, fmt.Errorf("gh json decode: %w", err)
	}
	rs := make([]Repo, 0, len(raw))
	for _, item := range raw {
		lang := ""
		if item.PrimaryLanguage != nil {
			lang = item.PrimaryLanguage.Name
		}
		branch := ""
		if item.DefaultBranch != nil {
			branch = item.DefaultBranch.Name
		}
		var topics []string
		if item.Topics != nil {
			for _, n := range item.Topics.Nodes {
				if n.Topic.Name != "" {
					topics = append(topics, n.Topic.Name)
				}
			}
		}
		rs = append(rs, Repo{
			Name:          item.Name,
			FullName:      item.NameWithOwner,
			Description:   item.Description,
			HTMLURL:       item.URL,
			Language:      lang,
			DefaultBranch: branch,
			Topics:        topics,
			Stars:         item.StargazerCount,
			Forks:         item.ForkCount,
			UpdatedAt:     item.UpdatedAt,
			Archived:      item.IsArchived,
			Fork:          item.IsFork,
			Owner:         owner,
		})
	}
	if limit > 0 && len(rs) > limit {
		rs = rs[:limit]
	}
	sort.Slice(rs, func(i, j int) bool { return rs[i].FullName < rs[j].FullName })
	return rs, nil
}

func fetchAllREST(owner string, limit int) ([]Repo, error) {
	var rs []Repo
	for page := 1; ; page++ {
		endpoint := "https://api.github.com/users/" + urlPathEscape(owner) +
			"/repos?per_page=100&page=" + strconv.Itoa(page) + "&type=owner"
		batch, next, err := fetchPage(endpoint)
		if err != nil {
			return nil, err
		}
		rs = append(rs, batch...)
		if limit > 0 && len(rs) >= limit {
			rs = rs[:limit]
			break
		}
		if !next {
			break
		}
	}
	sort.Slice(rs, func(i, j int) bool { return rs[i].FullName < rs[j].FullName })
	return rs, nil
}

// fetchPage performs one REST API call and reports whether a next page exists.
func fetchPage(endpoint string) ([]Repo, bool, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2026-03-10")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("github API: %s", resp.Status)
	}
	var raw []githubRepo
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, false, err
	}
	rs := make([]Repo, 0, len(raw))
	for _, item := range raw {
		r := item.Repo
		r.Owner = item.Owner.Login
		rs = append(rs, r)
	}
	next, err := nextPage(resp)
	return rs, next, err
}

// nextPage inspects the Link header for a rel="next" link.
func nextPage(resp *http.Response) (bool, error) {
	for _, link := range resp.Header.Values("Link") {
		if contains(link, `rel="next"`) {
			return true, nil
		}
	}
	return false, nil
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

// urlPathEscape escapes a path segment without touching slashes.
func urlPathEscape(v string) string {
	var esc string
	for _, c := range v {
		ok := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '-' || c == '.'
		if ok {
			esc += string(c)
		} else {
			esc += fmt.Sprintf("%%%02X", c)
		}
	}
	return esc
}