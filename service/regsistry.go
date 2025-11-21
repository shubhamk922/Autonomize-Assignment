package service

import (
	"example.com/team-monitoring/domain"
	"example.com/team-monitoring/service/ports/out"
)

type ToolRegistry struct {
	tools map[string]out.AITool
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{tools: make(map[string]out.AITool)}
}

func (r *ToolRegistry) Register(tool out.AITool) {
	r.tools[tool.Name()] = tool
}

func (r *ToolRegistry) Get(name string) (out.AITool, bool) {
	t, ok := r.tools[name]
	return t, ok
}

func (r *ToolRegistry) GetDefinitions() []domain.ToolDefinition {
	var defintions []domain.ToolDefinition
	for _, val := range r.tools {
		defintions = append(defintions, val.Defintion())
	}
	return defintions
}
