package jira

import (
	"context"
	"encoding/json"

	"example.com/team-monitoring/domain"
	"example.com/team-monitoring/infra/logger"
	"example.com/team-monitoring/service/ports/out"
)

type GetIssueStatusTool struct {
	Jira out.JiraPort
	Log  logger.Logger
}

// Arguments accepted from the LLM/tool call
type getIssueStatusArgs struct {
	IssueKey string `json:"issueKey"`
}

func (t *GetIssueStatusTool) Name() string { return "get_issue_status" }

func (t *GetIssueStatusTool) Definition() domain.ToolDefinition {
	return domain.ToolDefinition{
		Name:        "get_issue_status",
		Description: "Fetch the current status of a Jira issue by issue key (e.g. PROJ-123). Returns status name and basic issue info.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"issueKey": map[string]interface{}{
					"type":        "string",
					"description": "Jira issue key, e.g. PROJ-101",
				},
			},
			"required": []string{"issueKey"},
		},
	}
}

func (t *GetIssueStatusTool) Execute(ctx context.Context, raw json.RawMessage) (interface{}, error) {
	var args getIssueStatusArgs
	if err := json.Unmarshal(raw, &args); err != nil {
		return nil, err
	}
	t.Log.Infof("Running %s", t.Name())
	// Make a call to the JiraPort
	return t.Jira.GetStatus(ctx, args.IssueKey)
}
