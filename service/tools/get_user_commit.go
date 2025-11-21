package tools

import (
	"context"
	"encoding/json"

	"example.com/team-monitoring/domain"
	"example.com/team-monitoring/service/ports/out"
)

type GetUserCommitsTool struct {
	Github out.GithubPort
}

func (t *GetUserCommitsTool) Name() string { return "get_user_commits" }

type getUserCommitsArgs struct {
	Username string `json:"username"`
	Repo     string `json:"repo"`
	Since    string `json:"since"`
	Until    string `json:"until"`
}

func (t *GetUserCommitsTool) Execute(ctx context.Context, raw json.RawMessage) (interface{}, error) {
	var args getUserCommitsArgs
	if err := json.Unmarshal(raw, &args); err != nil {
		return nil, err
	}
	return t.Github.GetUserCommits(ctx, args.Username, args.Repo, args.Since, args.Until)
}

func (t *GetUserCommitsTool) Definition() domain.ToolDefinition {
	return domain.ToolDefinition{
		Name:        "get_user_commits",
		Description: "Retrieve commits from GitHub for a user. If repo is not provided, fetch commits across all repos. If username is not provided, assume current authenticated user.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"username": map[string]interface{}{
					"type":        "string",
					"description": "GitHub username. Optional. If not provided, use default authenticated user.",
				},
				"repo": map[string]interface{}{
					"type":        "string",
					"description": "Repository name. Optional. If not provided, fetch across all repos.",
				},
				"since": map[string]interface{}{
					"type":        "string",
					"format":      "date-time",
					"description": "ISO-8601 timestamp. If user uses words like 'this week', 'recently', 'last month', the model must convert them into correct date-time relative to TODAY. Optional.",
				},
				"until": map[string]interface{}{
					"type":        "string",
					"format":      "date-time",
					"description": "ISO-8601 timestamp. If user uses words like 'this week', 'recently', 'last month', the model must convert them into correct date-time relative to TODAY. Optional.",
				},
			},
			"required": []string{}, // NO REQUIRED FIELDS
		},
	}
}
