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
	GitHub []GitHubCommit
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

type GitHubCommit struct {
	Repo    string `json:"repo"`
	Sha     string `json:"sha"`
	Message string `json:"message"`
	Author  string `json:"author"`
	Date    string `json:"date"`
	Url     string `json:"url"`
}

type PullRequest struct {
	Title     string `json:"title"`
	Repo      string `json:"repo"`
	Number    int    `json:"number"`
	State     string `json:"state"`
	Url       string `json:"url"`
	CreatedAt string `json:"created_at"`
	Merged    bool   `json:"merged"`
}

type GitHubRepoContribution struct {
	Repo            string `json:"repo"`
	LastCommittedAt string `json:"last_committed_at,omitempty"`
	ActivityType    string `json:"activity_type"` // "push" or "pull_request"
}
