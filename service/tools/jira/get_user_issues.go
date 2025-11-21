package jira

import (
	"context"
	"encoding/json"

	"example.com/team-monitoring/domain"
	"example.com/team-monitoring/infra/logger"
	"example.com/team-monitoring/service/ports/out"
)

type GetUserIssuesTool struct {
	Jira out.JiraPort // <-- clean architecture: depends on a port
	Log  logger.Logger
}

type getUserIssuesArgs struct {
	Assignee string `json:"assignee"`
	Status   string `json:"status"`
	Project  string `json:"project"`
}

func (t *GetUserIssuesTool) Name() string { return "get_user_issues" }

func (t *GetUserIssuesTool) Definition() domain.ToolDefinition {
	return domain.ToolDefinition{
		Name:        "get_user_issues",
		Description: "Search Jira issues using structured filters. Converts natural language (e.g., 'my open tickets') into proper JQL.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"assignee": map[string]interface{}{
					"type":        "string",
					"description": "Filter by assignee (Jira username). Optional.",
				},
				"status": map[string]interface{}{
					"type":        "string",
					"description": "Issue status like 'To Do', 'In Progress', 'Done'. Optional.",
				},
				"project": map[string]interface{}{
					"type":        "string",
					"description": "Jira project key. Optional.",
				},
				"text": map[string]interface{}{
					"type":        "string",
					"description": "Free-text search across summary/description. Optional.",
				},
			},
			"required": []string{},
		},
	}
}

func (t *GetUserIssuesTool) Execute(ctx context.Context, raw json.RawMessage) (interface{}, error) {
	var args getUserIssuesArgs
	if err := json.Unmarshal(raw, &args); err != nil {
		return nil, err
	}
	if args.Assignee == "" {
		args.Assignee = "shubham"
	}
	q := domain.JiraQuery{
		Assignee: args.Assignee,
		Status:   args.Status,
		Project:  args.Project,
	}
	t.Log.Infof("Running %s", t.Name())
	return t.Jira.GetIssues(ctx, q)
}
