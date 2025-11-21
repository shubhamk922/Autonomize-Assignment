package out

import (
	"context"

	"example.com/team-monitoring/domain"
)

type GithubPort interface {
	Get(ctx context.Context, url string, target interface{}) error
	ListUserRepos(ctx context.Context, username string) ([]string, error)
	FetchUserEvents(ctx context.Context, username string) ([]domain.GitHubEvent, error)
	FetchCommitsFromRepo(ctx context.Context, username string, repo string, since string, until string) ([]domain.GitHubCommit, error)
	GetRepoPRs(ctx context.Context, repo string, filter string) ([]domain.PullRequest, error)
	GetUserWidePRs(ctx context.Context, username string, filter string) ([]domain.PullRequest, error)
}
