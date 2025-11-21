package domain

type MemberActivity struct {
	Name   string
	Jira   interface{}
	GitHub interface{}
}

type Message struct {
	Role    string
	Content string
}

// AIResponse: carries content and optional tool-call metadata
type AIResponse struct {
	Content  string
	ToolCall *ToolCall
}

type ToolCall struct {
	ID        string
	Name      string
	Arguments string
}

type ToolDefinition struct {
	Name        string
	Description string
	Parameters  map[string]interface{}
}
