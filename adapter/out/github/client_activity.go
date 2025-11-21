package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/team-monitoring/domain"
)

func (c *GithubClient) GetUserContributedRepos(ctx context.Context, username, since string) ([]domain.GitHubRepoContribution, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+c.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var events []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, err
	}

	sinceTime := time.Time{}
	if since != "" {
		t2, err := time.Parse(time.RFC3339, since)
		if err == nil {
			sinceTime = t2
		}
	}

	seen := map[string]domain.GitHubRepoContribution{}

	for _, e := range events {
		// parse event creation time
		createdAtStr, ok := e["created_at"].(string)
		if !ok {
			continue
		}
		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			continue
		}
		if !sinceTime.IsZero() && createdAt.Before(sinceTime) {
			continue
		}

		// identify repo
		repoObj, ok := e["repo"].(map[string]interface{})
		if !ok {
			continue
		}
		repoName, _ := repoObj["name"].(string)

		typ, _ := e["type"].(string)
		activityType := ""
		switch typ {
		case "PushEvent":
			activityType = "push"
		case "PullRequestEvent":
			activityType = "pull_request"
		default:
			continue
		}

		// record
		contrib := seen[repoName]
		contrib.Repo = repoName
		contrib.ActivityType = activityType
		contrib.LastCommittedAt = createdAtStr
		seen[repoName] = contrib
	}

	var result []domain.GitHubRepoContribution
	for _, v := range seen {
		result = append(result, v)
	}

	return result, nil
}
