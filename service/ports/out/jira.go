package out

import "example.com/team-monitoring/domain"

type JiraPort interface {
	GetUserIssues(domain.JiraQuery) ([]domain.JiraIssue, error)
}
