package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/team-monitoring/domain"
)

type GithubClient struct {
	Token string
}

func New(token string) *GithubClient {
	return &GithubClient{Token: token}
}

func (c *GithubClient) GetUserActivity(ctx context.Context, user string) ([]domain.GitHubActivity, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/users/"+user+"/events", nil)
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
	json.NewDecoder(resp.Body).Decode(&events)

	var activity []domain.GitHubActivity
	for _, e := range events {
		activity = append(activity, domain.GitHubActivity{
			Repo: e["repo"].(map[string]interface{})["name"].(string),
			Time: e["created_at"].(string),
		})
	}

	return activity, nil
}

func (c *GithubClient) GetUserCommits(
	ctx context.Context,
	username string,
	repo string,
	since string,
	until string,
) ([]domain.GitHubCommit, error) {

	if repo != "" {
		// CASE: commits from specific repo
		return c.fetchCommitsFromRepo(ctx, username, repo, since, until)
	}

	// CASE: all repos → list repos first
	repos, err := c.listUserRepos(ctx, username)
	if err != nil {
		return nil, err
	}

	allCommits := []domain.GitHubCommit{}

	for _, r := range repos {
		commits, err := c.fetchCommitsFromRepo(ctx, username, r, since, until)
		fmt.Println("Commits %+v", commits)
		if err != nil {
			// skip repo failures → not fatal
			continue
		}
		allCommits = append(allCommits, commits...)
	}

	return allCommits, nil
}

func (c *GithubClient) fetchCommitsFromRepo(
	ctx context.Context,
	username string,
	repo string,
	since string,
	until string,
) ([]domain.GitHubCommit, error) {

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?author=%s", username, repo, username)

	if since != "" {
		url += "&since=" + since
	}
	if until != "" {
		url += "&until=" + until
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+c.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jsonCommits []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&jsonCommits)

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

	return commits, nil
}

func (c *GithubClient) listUserRepos(ctx context.Context, username string) ([]string, error) {

	url := fmt.Sprintf("https://api.github.com/users/%s/repos", username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+c.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repos []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&repos)

	list := []string{}
	for _, r := range repos {
		list = append(list, r["name"].(string))
	}
	return list, nil
}
