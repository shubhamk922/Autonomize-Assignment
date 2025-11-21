package tools

import (
	"context"
	"encoding/json"
	"time"

	"example.com/team-monitoring/domain"
	"example.com/team-monitoring/infra/logger"
	"example.com/team-monitoring/service/ports/out"
)

type GetUserContributedReposTool struct {
	Github out.GithubPort
	Log    logger.Logger
}

type GetUserContributedReposArgs struct {
	Username string `json:"username"`
	Since    string `json:"since"` // ISO-8601 timestamp, optional
}

func NewGetUserContributedReposTool() *GetUserContributedReposTool {
	return &GetUserContributedReposTool{}
}

func (t *GetUserContributedReposTool) Name() string {
	return "get_user_contributed_repos"
}

func (t *GetUserContributedReposTool) Execute(ctx context.Context, raw json.RawMessage) (interface{}, error) {
	var args GetUserContributedReposArgs
	if err := json.Unmarshal(raw, &args); err != nil {
		return "", err
	}
	t.Log.Infof("Running %s", t.Name())
	return t.GetUserContributedRepos(ctx, args.Username, args.Since)
}

func (t *GetUserContributedReposTool) Definition() domain.ToolDefinition {
	return domain.ToolDefinition{
		Name:        "get_user_contributed_repos",
		Description: "List repositories the user recently contributed to. If the user mentions time windows like 'this week', 'yesterday', 'last 7 days', convert them into actual date-time values in ISO format based on the current date.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"username": map[string]interface{}{
					"type":        "string",
					"description": "GitHub username. Optional if you assume current user.",
				},
				"since": map[string]interface{}{
					"type":        "string",
					"format":      "date-time",
					"description": "ISO-8601 timestamp. If user uses words like 'this week', 'recently', 'last month', the model must convert them into correct date-time relative to TODAY.",
				},
			},
			"required": []string{},
		},
	}
}

func (c *GetUserContributedReposTool) GetUserContributedRepos(
	ctx context.Context,
	username string,
	since string,
) ([]domain.GitHubRepoContribution, error) {

	events, err := c.Github.FetchUserEvents(ctx, username)
	if err != nil {
		return nil, err
	}

	sinceTime := parseTime(since)
	seen := make(map[string]domain.GitHubRepoContribution)

	for _, e := range events {
		created := parseTime(e.CreatedAt)
		if !sinceTime.IsZero() && created.Before(sinceTime) {
			continue
		}

		if e.Repo.Name == "" {
			continue
		}

		activityType := classifyEventType(e.Type)
		if activityType == "" {
			continue
		}

		seen[e.Repo.Name] = domain.GitHubRepoContribution{
			Repo:            e.Repo.Name,
			ActivityType:    activityType,
			LastCommittedAt: e.CreatedAt,
		}
	}

	return mapToSlice(seen), nil
}

func classifyEventType(t string) string {
	switch t {
	case "PushEvent":
		return "push"
	case "PullRequestEvent":
		return "pull_request"
	default:
		return ""
	}
}

func mapToSlice(m map[string]domain.GitHubRepoContribution) []domain.GitHubRepoContribution {
	out := make([]domain.GitHubRepoContribution, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	return out
}

func parseTime(t string) time.Time {
	if t == "" {
		return time.Time{}
	}
	parsed, _ := time.Parse(time.RFC3339, t)
	return parsed
}
