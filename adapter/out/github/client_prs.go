package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"example.com/team-monitoring/domain"
)

func (c *GithubClient) GetUserPRs(ctx context.Context, username, repo, filter string) ([]domain.PullRequest, error) {

	if filter == "" {
		filter = "open"
	}

	// If repo provided
	if repo != "" {
		parts := strings.Split(repo, "/")
		if len(parts) != 2 {
			return nil, fmt.Errorf("repo must be owner/repo")
		}

		owner, repoName := parts[0], parts[1]

		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?state=%s", owner, repoName, mapFilter(filter))

		req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
		req.Header.Set("Authorization", "token "+c.Token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var prs []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&prs)

		return mapPRList(prs, repo), nil
	}

	// Search across ALL repos
	search := fmt.Sprintf("q=author:%s+type:pr+state:%s", username, filter)
	url := "https://api.github.com/search/issues?" + search

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Authorization", "token "+c.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		Items []map[string]interface{} `json:"items"`
	}
	json.NewDecoder(resp.Body).Decode(&data)

	return mapSearchPRList(data.Items), nil
}

func mapFilter(f string) string {
	if f == "merged" {
		// GitHub doesn't have merged state in /pulls endpoint
		return "closed"
	}
	if f == "all" {
		return "all"
	}
	return f
}

func mapPRList(in []map[string]interface{}, repo string) []domain.PullRequest {
	out := []domain.PullRequest{}

	for _, pr := range in {
		out = append(out, domain.PullRequest{
			Title:     pr["title"].(string),
			Repo:      repo,
			Number:    int(pr["number"].(float64)),
			State:     pr["state"].(string),
			Url:       pr["html_url"].(string),
			CreatedAt: pr["created_at"].(string),
			Merged:    pr["merged_at"] != nil,
		})
	}
	return out
}

func mapSearchPRList(items []map[string]interface{}) []domain.PullRequest {
	out := []domain.PullRequest{}
	for _, i := range items {
		repo := strings.Split(i["repository_url"].(string), "/repos/")[1]

		out = append(out, domain.PullRequest{
			Title:     i["title"].(string),
			Repo:      repo,
			Number:    int(i["number"].(float64)),
			State:     strings.ToLower(i["state"].(string)),
			Url:       i["html_url"].(string),
			CreatedAt: i["created_at"].(string),
			Merged:    strings.Contains(i["state"].(string), "merged"),
		})
	}
	return out
}
