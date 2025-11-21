package tools

import (
	"context"
	"encoding/json"

	"example.com/team-monitoring/domain"
	"example.com/team-monitoring/service/ports/out"
)

type GetUserContributedReposTool struct {
	Github out.GithubPort
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
	return t.Github.GetUserContributedRepos(ctx, args.Username, args.Since)
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
