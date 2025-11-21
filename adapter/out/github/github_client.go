package github

import (
	"context"
	"encoding/json"
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
