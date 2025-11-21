package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"example.com/team-monitoring/adapter/out/user"
	"example.com/team-monitoring/domain"
	"example.com/team-monitoring/infra/cache"
	"example.com/team-monitoring/infra/logger"
)

type GithubClient struct {
	Token      string
	Log        logger.Logger
	IdentityDB *user.UserIdentityDB
	Cache      cache.Cache
}

func New(token string) *GithubClient {
	return &GithubClient{Token: token}
}

func (c *GithubClient) Get(ctx context.Context, url string, target interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+c.Token)
	c.Log.Debug("Request", "req", req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("github API error: %s", resp.Status)
	}

	json.NewDecoder(resp.Body).Decode(target)
	c.Log.Debug("Response", "res", target)
	return nil
}

func (c *GithubClient) FetchUserEvents(ctx context.Context, username string) ([]domain.GitHubEvent, error) {
	cacheKey := fmt.Sprintf("github_events_%s", c.IdentityDB.GetGithubId(username))

	// Try cache first
	var events []domain.GitHubEvent
	if c.Cache != nil {
		if err := c.Cache.Get(cacheKey, &events); err == nil {
			return events, nil
		}
	}

	url := fmt.Sprintf("https://api.github.com/users/%s/events", c.IdentityDB.GetGithubId(username))

	if err := c.Get(ctx, url, &events); err != nil {
		return nil, err
	}
	if c.Cache != nil {
		_ = c.Cache.Set(cacheKey, events, 5*time.Minute)
	}
	return events, nil
}

func (c *GithubClient) ListUserRepos(ctx context.Context, username string) ([]string, error) {
	cacheKey := fmt.Sprintf("github_repos_%s", c.IdentityDB.GetGithubId(username))

	// Try cache first
	var cached []string
	if c.Cache != nil {
		if err := c.Cache.Get(cacheKey, &cached); err == nil {
			return cached, nil
		}
	}
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", c.IdentityDB.GetGithubId(username))

	var repos []map[string]interface{}

	if err := c.Get(ctx, url, &repos); err != nil {
		return nil, err
	}

	list := []string{}
	for _, r := range repos {
		list = append(list, r["name"].(string))
	}
	if c.Cache != nil {
		_ = c.Cache.Set(cacheKey, list, 10*time.Minute) // TTL configurable
	}
	return list, nil
}

func (c *GithubClient) FetchCommitsFromRepo(
	ctx context.Context,
	username string,
	repo string,
	since string,
	until string,
) ([]domain.GitHubCommit, error) {

	cacheKey := fmt.Sprintf("github_commits_%s_%s_%s_%s",
		c.IdentityDB.GetGithubId(username), repo, since, until)

	// Try fetching from cache first
	if c.Cache != nil {
		var cached []domain.GitHubCommit
		if err := c.Cache.Get(cacheKey, &cached); err == nil {
			return cached, nil
		}
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?author=%s", c.IdentityDB.GetGithubId(username), repo, c.IdentityDB.GetGithubId(username))

	if since != "" {
		url += "&since=" + since
	}
	if until != "" {
		url += "&until=" + until
	}

	var jsonCommits []map[string]interface{}
	if err := c.Get(ctx, url, &jsonCommits); err != nil {
		return nil, err
	}
	var commits []domain.GitHubCommit

	for _, cmt := range jsonCommits {
		commitObj := cmt["commit"].(map[string]interface{})
		authorObj := commitObj["author"].(map[string]interface{})

		commits = append(commits, domain.GitHubCommit{
			Repo:    repo,
			Sha:     cmt["sha"].(string),
			Message: commitObj["message"].(string),
			Author:  authorObj["name"].(string),
			Date:    authorObj["date"].(string),
			Url:     cmt["html_url"].(string),
		})
	}

	// Save to cache with TTL
	if c.Cache != nil {
		_ = c.Cache.Set(cacheKey, commits, 10*time.Minute)
	}

	return commits, nil
}

func (c *GithubClient) GetRepoPRs(
	ctx context.Context,
	repo string,
	filter string,
) ([]domain.PullRequest, error) {

	owner, repoName, err := splitRepo(repo)
	if err != nil {
		return nil, err
	}
	cacheKey := fmt.Sprintf("github_prs_%s_%s_%s", owner, repoName, filter)

	// Try cache first
	if c.Cache != nil {
		var cached []domain.PullRequest
		if err := c.Cache.Get(cacheKey, &cached); err == nil {
			return cached, nil
		}
	}
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/pulls?state=%s",
		owner, repoName, mapFilter(filter),
	)

	var prs []domain.GitHubPR
	if err := c.Get(ctx, url, &prs); err != nil {
		return nil, err
	}
	result := mapPRList(prs, repo)
	if c.Cache != nil {
		_ = c.Cache.Set(cacheKey, result, 10*time.Minute)
	}

	return result, nil
}

func mapPRList(in []domain.GitHubPR, repo string) []domain.PullRequest {
	out := []domain.PullRequest{}

	for _, pr := range in {
		out = append(out, domain.PullRequest{
			Title:     pr.Title,
			Repo:      repo,
			Number:    int(pr.Number),
			State:     pr.State,
			Url:       pr.HTMLURL,
			CreatedAt: pr.CreatedAt,
			Merged:    pr.MergedAt != nil,
		})
	}
	return out
}

func (c *GithubClient) GetUserWidePRs(
	ctx context.Context,
	username string,
	filter string,
) ([]domain.PullRequest, error) {

	ghUser := c.IdentityDB.GetGithubId(username)
	if ghUser == "" {
		return nil, fmt.Errorf("github identity not found for user: %s", username)
	}
	cacheKey := fmt.Sprintf("github_user_wide_prs_%s_%s", ghUser, filter)
	if c.Cache != nil {
		var cached []domain.PullRequest
		if err := c.Cache.Get(cacheKey, &cached); err == nil {
			return cached, nil
		}
	}
	// build query using url.Values
	params := url.Values{}
	params.Set("q", fmt.Sprintf("author:%s type:pr state:%s", ghUser, filter))
	params.Set("sort", "created")
	params.Set("order", "desc")
	url := "https://api.github.com/search/issues?" + params.Encode()

	var data domain.GitHubSearchIssues
	if err := c.Get(ctx, url, &data); err != nil {
		return nil, err
	}
	result := mapSearchPRList(data.Items)
	if c.Cache != nil {
		_ = c.Cache.Set(cacheKey, result, 10*time.Minute)
	}
	return result, nil
}

func mapSearchPRList(items []domain.GitHubSearchItem) []domain.PullRequest {
	out := []domain.PullRequest{}
	for _, i := range items {
		repo := strings.Split(i.RepositoryURL, "/repos/")[1]

		out = append(out, domain.PullRequest{
			Title:     i.Title,
			Repo:      repo,
			Number:    int(i.Number),
			State:     strings.ToLower(i.State),
			Url:       i.HTMLURL,
			CreatedAt: i.CreatedAt,
			Merged:    strings.Contains(i.State, "merged"),
		})
	}
	return out
}

func splitRepo(repo string) (string, string, error) {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("repo must be owner/repo")
	}
	return parts[0], parts[1], nil
}

func mapFilter(state string) string {
	if state == "merged" {
		return "closed" // GitHub API quirk
	}
	if state == "all" {
		return "all"
	}
	return state
}
