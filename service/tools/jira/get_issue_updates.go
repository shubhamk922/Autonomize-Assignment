package jira

import (
	"context"
	"encoding/json"

	"example.com/team-monitoring/domain"
	"example.com/team-monitoring/infra/logger"
	"example.com/team-monitoring/service/ports/out"
)

type GetIssueUpdatesTool struct {
	Jira out.JiraPort
	Log  logger.Logger
}

func (t *GetIssueUpdatesTool) Name() string { return "get_issue_updates" }

type getIssueUpdatesArgs struct {
	IssueKey string `json:"issueKey"`
}

func (t *GetIssueUpdatesTool) Definition() domain.ToolDefinition {
	return domain.ToolDefinition{
		Name:        "get_issue_updates",
		Description: "Fetch recent updates for a Jira issue including changelog, status transitions, and update timestamps.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"issueKey": map[string]interface{}{
					"type":        "string",
					"description": "Jira issue key, e.g. PROJ-123",
				},
			},
			"required": []string{"issueKey"},
		},
	}
}

func (t *GetIssueUpdatesTool) Execute(ctx context.Context, raw json.RawMessage) (interface{}, error) {
	var args getIssueUpdatesArgs
	if err := json.Unmarshal(raw, &args); err != nil {
		return nil, err
	}
	t.Log.Infof("Running %s", t.Name())
	return t.Jira.GetUpdates(ctx, args.IssueKey, 10)
}
