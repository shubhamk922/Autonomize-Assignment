package openai

import (
	"context"
	"encoding/json"
	"fmt"

	"example.com/team-monitoring/domain"
	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/shared"
	"github.com/openai/openai-go/shared/constant"
)

type OpenAIClient struct {
	client openai.Client
	model  string
}

func NewOpenAIClient(model string) *OpenAIClient {
	// OPENAI_API_KEY is automatically read by NewClient()

	client := openai.NewClient()

	if model == "" {
		model = "gpt-4o-mini" // default model
	}

	return &OpenAIClient{
		client: client,
		model:  model,
	}
}

// Generate implements domain.AIPort
func (o *OpenAIClient) Generate(act *domain.MemberActivity) (string, error) {
	ctx := context.Background()

	// Prepare prompt
	prompt := fmt.Sprintf(
		"Summarize the following engineering activity for member %s.\n\nJIRA: %+v\n\nGitHub: %+v\n",
		act.Name, act.Jira, act.GitHub,
	)

	resp, err := o.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: o.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Role: "user", // correct enum
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String(prompt),
					},
				},
			},
		},
	})

	if err != nil {
		return "", fmt.Errorf("openai chat error: %w", err)
	}

	// Extract text
	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}

func (o *OpenAIClient) Chat(ctx context.Context, msgs []domain.Message, tools []domain.ToolDefinition) (domain.AIResponse, error) {
	var sdkMsgs []openai.ChatCompletionMessageParamUnion
	for _, m := range msgs {
		if m.Role == "system" {
			sdkMsgs = append(sdkMsgs, openai.ChatCompletionMessageParamUnion{OfSystem: &openai.ChatCompletionSystemMessageParam{Content: openai.ChatCompletionSystemMessageParamContentUnion{OfString: openai.String(m.Content)}}})
		} else if m.Role == "user" {
			sdkMsgs = append(sdkMsgs, openai.ChatCompletionMessageParamUnion{OfUser: &openai.ChatCompletionUserMessageParam{Role: "user", Content: openai.ChatCompletionUserMessageParamContentUnion{OfString: openai.String(m.Content)}}})
		}
	}

	// Define get_user_commits tool
	/*getUserCommitsTool := openai.ChatCompletionToolParam{
		Type: constant.Function("function"),
		Function: shared.FunctionDefinitionParam{
			Name:        "get_user_commits",
			Description: param.Opt[string]{Value: "Get recent commits pushed by a GitHub user"},
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"username": map[string]interface{}{"type": "string"},
				},
				"required": []string{"username"},
			},
		},
	}

	memberActivityTool := openai.ChatCompletionToolParam{
		Type: constant.Function("function"),
		Function: shared.FunctionDefinitionParam{
			Name:        "get_member_activity",
			Description: param.Opt[string]{Value: "Get activity summary for a user"},
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"username": map[string]interface{}{"type": "string"},
				},
				"required": []string{"username"},
			},
		},
	}*/

	var sdkTools []openai.ChatCompletionToolParam
	for _, t := range tools {
		sdkTools = append(sdkTools, openai.ChatCompletionToolParam{
			Type: constant.Function("function"),
			Function: shared.FunctionDefinitionParam{
				Name:        t.Name,
				Description: param.Opt[string]{Value: t.Description},
				Parameters:  t.Parameters,
			},
		})
	}

	req := openai.ChatCompletionNewParams{
		Model:    (o.model),
		Messages: sdkMsgs,
		Tools:    sdkTools,
		ToolChoice: openai.ChatCompletionToolChoiceOptionUnionParam{
			OfAuto: param.Opt[string]{Value: "auto"},
		},
	}

	resp, err := o.client.Chat.Completions.New(ctx, req)
	if err != nil {
		return domain.AIResponse{}, fmt.Errorf("openai chat error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return domain.AIResponse{}, fmt.Errorf("no choices")
	}

	ch := resp.Choices[0]

	// Check tool calls
	if len(ch.Message.ToolCalls) > 0 {
		tc := ch.Message.ToolCalls[0]
		return domain.AIResponse{ToolCall: &domain.ToolCall{ID: tc.ID, Name: tc.Function.Name, Arguments: tc.Function.Arguments}}, nil
	}

	return domain.AIResponse{Content: ch.Message.Content}, nil
}

func (o *OpenAIClient) CompleteTool(ctx context.Context, original domain.AIResponse, toolResult interface{}) (string, error) {

	// Build messages: assistant message indicating tool call + tool message containing result
	assistantMsg := openai.ChatCompletionAssistantMessageParam{
		Role: "assistant",
		ToolCalls: []openai.ChatCompletionMessageToolCallParam{
			{
				ID: original.ToolCall.ID,
				Function: openai.ChatCompletionMessageToolCallFunctionParam{
					Name:      original.ToolCall.Name,
					Arguments: original.ToolCall.Arguments,
				},
			},
		},
	}

	toolJSON, _ := json.Marshal(toolResult)
	toolMsg := openai.ChatCompletionToolMessageParam{
		Role:       "tool",
		ToolCallID: original.ToolCall.ID,
		Content: openai.ChatCompletionToolMessageParamContentUnion{
			OfString: openai.String(string(toolJSON)),
		},
	}

	msgs := []openai.ChatCompletionMessageParamUnion{
		{OfAssistant: &assistantMsg},
		{OfTool: &toolMsg},
	}

	req := openai.ChatCompletionNewParams{
		Model:    o.model,
		Messages: msgs,
	}

	resp, err := o.client.Chat.Completions.New(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices on second response")
	}

	return resp.Choices[0].Message.Content, nil
}
