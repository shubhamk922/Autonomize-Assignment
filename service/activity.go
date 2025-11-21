package service

import (
	"context"

	"example.com/team-monitoring/domain"
	"example.com/team-monitoring/service/ports/out"
)

type ActivityService struct {
	Jira   out.JiraPort
	Github out.GithubPort
	AI     out.AIPort
}

func NewActivityService(j out.JiraPort, g out.GithubPort, ai out.AIPort) *ActivityService {
	return &ActivityService{Jira: j, Github: g, AI: ai}
}

func (s *ActivityService) GetMemberActivity(ctx context.Context, name string) (*domain.MemberActivity, error) {
	jiraData, _ := s.Jira.GetUserIssues(domain.JiraQuery{})
	githubData, _ := s.Github.GetUserActivity(ctx, name)

	activity := &domain.MemberActivity{
		Name:   name,
		Jira:   jiraData,
		GitHub: githubData,
	}

	return activity, nil
}

func (s *ActivityService) GenerateResponse(activity *domain.MemberActivity) (string, error) {
	return s.AI.Generate(activity)
}
