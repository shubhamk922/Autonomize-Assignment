package service

import (
	"context"
	"fmt"

	"example.com/team-monitoring/domain"
	"example.com/team-monitoring/service/ports/out"
)

type ChatBot struct {
	AI       out.AIPort
	Registry *ToolRegistry
}

func NewBot(ai out.AIPort, registry *ToolRegistry) *ChatBot {
	return &ChatBot{AI: ai, Registry: registry}
}

func (b *ChatBot) Handle(ctx context.Context, input string) (string, error) {
	// build system + user messages
	msgs := []domain.Message{
		{Role: "system", Content: "You are an assistant that calls backend tools when needed."},
		{Role: "user", Content: input},
	}

	resp, err := b.AI.Chat(ctx, msgs, b.Registry.GetDefinitions())
	if err != nil {
		return "", err
	}

	// If no tool call → return message
	if resp.ToolCall == nil {
		return resp.Content, nil
	}
	tool, ok := b.Registry.Get(resp.ToolCall.Name)
	if !ok {
		return "", fmt.Errorf("unknown tool: %s", resp.ToolCall.Name)
	}

	// Execute tool
	result, err := tool.Execute(ctx, []byte(resp.ToolCall.Arguments))
	if err != nil {
		result = err
	}

	// Let AI format final answer
	final, err := b.AI.CompleteTool(ctx, resp, result)
	if err != nil {
		return "Sorry ! I am unable to process your request ", nil
	}

	return final, nil
}
