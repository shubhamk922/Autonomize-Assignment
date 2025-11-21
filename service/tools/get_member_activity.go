package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"example.com/team-monitoring/domain"
	"example.com/team-monitoring/infra/logger"
	"example.com/team-monitoring/service/ports/out"
)

type GetMemberActivityTool struct {
	Jira   out.JiraPort
	Github out.GithubPort
	Log    logger.Logger
	Steps  map[string]out.AITool
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
	t.Log.Infof("Running %s", t.Name())
	if args.Username == "" {
		args.Username = "shubham"
	}
	return t.GetMemberActivity(ctx, args.Username)
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

func (s *GetMemberActivityTool) GetMemberActivity(ctx context.Context, name string) (*domain.MemberActivity, error) {
	type result struct {
		name string
		data interface{}
		err  error
	}

	ch := make(chan result, len(s.Steps))

	payloads := map[string][]byte{}

	jiraPayload, _ := json.Marshal(map[string]interface{}{
		"assignee": name,
	})
	payloads["get_user_issues"] = jiraPayload

	githubPayload, _ := json.Marshal(map[string]interface{}{
		"username": name,
	})
	payloads["get_user_commits"] = githubPayload

	for toolName, tool := range s.Steps {
		if payload, ok := payloads[toolName]; ok {
			go func(toolName string, tool out.AITool, payload []byte) {
				out, err := tool.Execute(ctx, payload)
				ch <- result{toolName, out, err}
			}(toolName, tool, payload)
		}
	}
	activity := &domain.MemberActivity{Name: name}
	var collectedErrors []error

	for i := 0; i < len(payloads); i++ {
		r := <-ch

		if r.err != nil {
			collectedErrors = append(collectedErrors, fmt.Errorf("%s failed: %w", r.name, r.err))
			continue
		}

		switch r.name {

		case "get_user_issues":
			if issues, ok := r.data.([]domain.JiraIssue); ok {
				activity.Jira = issues
			}

		case "get_user_commits":
			if commits, ok := r.data.([]domain.GitHubCommit); ok {
				activity.GitHub = commits
			}
		}
	}

	if len(collectedErrors) > 0 {
		return activity, collectedErrors[0]
	}

	return activity, nil
}
