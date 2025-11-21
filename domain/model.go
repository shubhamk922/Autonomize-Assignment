package domain

type JiraIssue struct {
	Key     string
	Summary string
	Status  string
	Id      string
}

type GitHubActivity struct {
	Repo      string
	CommitMsg string
	PRTitle   string
	Time      string
}

type MemberActivity struct {
	Name   string
	Jira   []JiraIssue
	GitHub []GitHubActivity
}

type Message struct {
	Role    string
	Content string
}

type Commit struct {
	Repo      string `json:"repo"`
	Message   string `json:"message"`
	URL       string `json:"url"`
	Timestamp string `json:"timestamp"`
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

type JiraQuery struct {
	Project  string
	Assignee string // “currentUser()” or actual name
	Status   string
}
