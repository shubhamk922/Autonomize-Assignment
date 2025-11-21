package jira

import (
	"encoding/json"
	"net/http"
	"net/url"

	"example.com/team-monitoring/domain"
)

type JiraClient struct {
	BaseURL string
	Token   string
}

func New(baseURL, token string) *JiraClient {
	return &JiraClient{BaseURL: baseURL, Token: token}
}

func (c *JiraClient) GetUserIssues(q domain.JiraQuery) ([]domain.JiraIssue, error) {

	jql := BuildJQL(q)
	encoded := url.QueryEscape(jql)

	req, _ := http.NewRequest(
		"GET",
		c.BaseURL+"/rest/api/3/search?jql="+encoded+"&maxResults=10",
		nil,
	)
	req.Header.Set("Authorization", "Basic "+c.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var data struct {
		Issues []struct {
			Id     string `json:id`
			Key    string `json:"key"`
			Fields struct {
				Summary string `json:"summary"`
				Status  struct {
					Name string `json:"name"`
				} `json:"status"`
			} `json:"fields"`
		} `json:"issues"`
	}
	json.NewDecoder(resp.Body).Decode(&data)

	var issues []domain.JiraIssue
	for _, i := range data.Issues {
		issues = append(issues, domain.JiraIssue{
			Key:     i.Key,
			Id:      i.Id,
			Summary: i.Fields.Summary,
			Status:  i.Fields.Status.Name,
		})
	}

	return issues, nil
}
