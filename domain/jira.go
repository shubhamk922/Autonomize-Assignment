package domain

type JiraIssue struct {
	Key     string
	Summary string
	Status  string
	Id      string
}

type JiraQuery struct {
	Project  string
	Assignee string // “currentUser()” or actual name
	Status   string
}

type JiraIssueStatus struct {
	Key    string `json:"key"`
	Status string `json:"status"`
}

type JiraIssueUpdate struct {
	IssueKey string              `json:"issueKey"`
	Updates  []JiraChangelogItem `json:"updates"`
}

type JiraChangelogItem struct {
	Author    string `json:"author"`
	Field     string `json:"field"`
	From      string `json:"from"`
	To        string `json:"to"`
	CreatedAt string `json:"createdAt"`
}
