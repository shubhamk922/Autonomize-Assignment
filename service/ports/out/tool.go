package out

import (
	"context"
	"encoding/json"

	"example.com/team-monitoring/domain"
)

type AITool interface {
	Name() string
	Definition() domain.ToolDefinition
	Execute(ctx context.Context, args json.RawMessage) (interface{}, error)
}
