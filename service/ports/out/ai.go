package out

import (
	"context"

	"example.com/team-monitoring/domain"
)

type AIPort interface {
	Generate(activity *domain.MemberActivity) (string, error)
	Chat(ctx context.Context, msgs []domain.Message, tools []domain.ToolDefinition) (domain.AIResponse, error)
	CompleteTool(ctx context.Context, original domain.AIResponse, toolResult interface{}) (string, error)
}
