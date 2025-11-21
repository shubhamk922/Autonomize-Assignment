package domain

type GitHubRepoContribution struct {
	Repo            string `json:"repo"`
	LastCommittedAt string `json:"last_committed_at,omitempty"`
	ActivityType    string `json:"activity_type"` // "push" or "pull_request"
}

type GitHubEvent struct {
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	Repo      struct {
		Name string `json:"name"`
	} `json:"repo"`
}

type GitHubPR struct {
	Title     string  `json:"title"`
	Number    int     `json:"number"`
	State     string  `json:"state"`
	HTMLURL   string  `json:"html_url"`
	CreatedAt string  `json:"created_at"`
	MergedAt  *string `json:"merged_at"`
}

type GitHubSearchIssues struct {
	Items []GitHubSearchItem `json:"items"`
}

type GitHubSearchItem struct {
	Title         string `json:"title"`
	Number        int    `json:"number"`
	State         string `json:"state"`
	HTMLURL       string `json:"html_url"`
	CreatedAt     string `json:"created_at"`
	RepositoryURL string `json:"repository_url"`
}

type GitHubActivity struct {
	Repo      string
	CommitMsg string
	PRTitle   string
	Time      string
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
