package tools

import (
	"context"
	"encoding/json"

	"example.com/team-monitoring/domain"
	"example.com/team-monitoring/service/ports/out"
)

type GetUserPRsTool struct {
	Github out.GithubPort
}

type GetUserPRsArgs struct {
	Username string `json:"username"`
	Repo     string `json:"repo"`
	Filter   string `json:"filter"` // open, closed, merged, all
}

func NewGetUserPRsTool() *GetUserPRsTool {
	return &GetUserPRsTool{}
}

func (t *GetUserPRsTool) Name() string {
	return "get_user_prs"
}

func (t *GetUserPRsTool) Execute(ctx context.Context, raw json.RawMessage) (interface{}, error) {
	var args GetUserPRsArgs
	if err := json.Unmarshal(raw, &args); err != nil {
		return "", err
	}

	return t.Github.GetUserPRs(ctx, args.Username, args.Repo, args.Filter)
}

func (t *GetUserPRsTool) Definition() domain.ToolDefinition {
	return domain.ToolDefinition{
		Name:        "get_user_prs",
		Description: "Fetch pull requests for a GitHub user. Supports filtering by state (open/closed/merged) and optionally by repository.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"username": map[string]interface{}{
					"type":        "string",
					"description": "GitHub username. Optional — if omitted, current authenticated user is used.",
				},
				"repo": map[string]interface{}{
					"type":        "string",
					"description": "Repository name (optional)",
				},
				"filter": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"open", "closed", "merged", "all"},
					"description": "Filter PRs by state",
				},
			},
		},
	}
}
