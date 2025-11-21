package tools

import (
	"context"
	"encoding/json"

	"example.com/team-monitoring/domain"
	"example.com/team-monitoring/service"
)

type GetMemberActivityTool struct {
	Svc *service.ActivityService
}

func (t *GetMemberActivityTool) Name() string { return "get_member_activity" }

type memberActivityArgs struct {
	Username string `json:"username"`
}

func (t *GetMemberActivityTool) Execute(ctx context.Context, raw json.RawMessage) (interface{}, error) {
	var args memberActivityArgs
	if err := json.Unmarshal(raw, &args); err != nil {
		return nil, err
	}
	return t.Svc.GetMemberActivity(ctx, args.Username)
}

func (t *GetMemberActivityTool) Definition() domain.ToolDefinition {
	return domain.ToolDefinition{
		Name:        "get_member_activity",
		Description: "Get Git and Jira activity summary for a user",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"username": map[string]interface{}{"type": "string"},
			},
			"required": []string{"username"},
		},
	}
}
