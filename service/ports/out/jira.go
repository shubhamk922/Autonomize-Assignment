package out

import (
	"context"

	"example.com/team-monitoring/domain"
)

type JiraPort interface {
	GetIssues(ctx context.Context, query domain.JiraQuery) ([]domain.JiraIssue, error)
	GetStatus(ctx context.Context, key string) (domain.JiraIssueStatus, error)
	GetUpdates(ctx context.Context, key string, limit int) (domain.JiraIssueUpdate, error)
}
