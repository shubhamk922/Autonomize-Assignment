package out

import (
	"context"

	"example.com/team-monitoring/domain"
)

type GithubPort interface {
	GetUserActivity(ctx context.Context, user string) ([]domain.GitHubActivity, error)
	GetUserCommits(ctx context.Context, username, repo, since, until string) ([]domain.GitHubCommit, error)
	GetUserPRs(ctx context.Context, username, repo, filter string) ([]domain.PullRequest, error)
	GetUserContributedRepos(ctx context.Context, username, since string) ([]domain.GitHubRepoContribution, error)
}
