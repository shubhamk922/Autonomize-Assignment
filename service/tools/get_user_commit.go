package tools

import (
	"context"
	"encoding/json"

	"example.com/team-monitoring/adapter/out/github"
	"example.com/team-monitoring/domain"
)

type GetUserCommitsTool struct {
	Github *github.GithubClient
}

func (t *GetUserCommitsTool) Name() string { return "get_user_commits" }

type getUserCommitsArgs struct {
	Username string `json:"username"`
}

func (t *GetUserCommitsTool) Execute(ctx context.Context, raw json.RawMessage) (interface{}, error) {
	var args getUserCommitsArgs
	if err := json.Unmarshal(raw, &args); err != nil {
		return nil, err
	}
	return t.Github.GetUserActivity(ctx, args.Username)
}

func (t *GetUserCommitsTool) Defintion() domain.ToolDefinition {
	return domain.ToolDefinition{
		Name:        "get_user_commits",
		Description: "Get recent commits pushed by a GitHub user",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"username": map[string]interface{}{"type": "string"},
			},
			"required": []string{"username"},
		},
	}
}
