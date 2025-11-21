package out

import (
	"context"

	"example.com/team-monitoring/domain"
)

type GithubPort interface {
	GetUserActivity(ctx context.Context, user string) ([]domain.GitHubActivity, error)
}
