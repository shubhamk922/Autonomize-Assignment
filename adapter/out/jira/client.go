package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"example.com/team-monitoring/adapter/out/user"
	"example.com/team-monitoring/domain"
	"example.com/team-monitoring/infra/cache"
	"example.com/team-monitoring/infra/logger"
)

type JiraClient struct {
	BaseURL    string
	Token      string
	Log        logger.Logger
	IdentityDB *user.UserIdentityDB
	Cache      cache.Cache
}

func New(baseURL, token string) *JiraClient {
	return &JiraClient{BaseURL: baseURL, Token: token}
}

func (c *JiraClient) GetIssues(ctx context.Context, q domain.JiraQuery) ([]domain.JiraIssue, error) {

	q.Assignee = c.IdentityDB.GetJiraId(q.Assignee)
	jql := BuildJQL(q)
	encoded := url.QueryEscape(jql)
	cacheKey := "jira_issues:" + jql

	// Try fetching from cache first
	var cached []domain.JiraIssue
	if c.Cache != nil {

		if err := c.Cache.Get(cacheKey, &cached); err == nil {
			return cached, nil
		}
	}
	req, _ := http.NewRequest(
		"GET",
		c.BaseURL+"/rest/api/3/search/jql?jql="+encoded+"&maxResults=10"+"&fields=*all",
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
	if c.Cache != nil {
		_ = c.Cache.Set(cacheKey, issues, 5*time.Minute)
	}

	return issues, nil
}

func (t *JiraClient) GetStatus(ctx context.Context, key string) (domain.JiraIssueStatus, error) {
	cacheKey := "jira_status:" + key

	// 1️⃣ Try fetching from cache
	if t.Cache != nil {
		var cached domain.JiraIssueStatus
		if err := t.Cache.Get(cacheKey, &cached); err == nil {
			return cached, nil
		}
	}
	req, _ := http.NewRequest(
		"GET",
		t.BaseURL+"/rest/api/3/issue/"+key+"?fields=status",
		nil,
	)
	req.Header.Set("Authorization", "Basic "+t.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return domain.JiraIssueStatus{}, err
	}
	defer resp.Body.Close()

	var data struct {
		Key    string `json:"key"`
		Fields struct {
			Status struct {
				Name string `json:"name"`
			} `json:"status"`
		} `json:"fields"`
	}

	json.NewDecoder(resp.Body).Decode(&data)

	result := domain.JiraIssueStatus{
		Key:    data.Key,
		Status: data.Fields.Status.Name,
	}
	if t.Cache != nil {
		_ = t.Cache.Set(cacheKey, result, 5*time.Minute)
	}
	return result, nil
}

func (t *JiraClient) GetUpdates(ctx context.Context, key string, limit int) (domain.JiraIssueUpdate, error) {
	cacheKey := fmt.Sprintf("jira_updates:%s", key)

	// 1️⃣ Check cache first
	if t.Cache != nil {
		var cached domain.JiraIssueUpdate
		if err := t.Cache.Get(cacheKey, &cached); err == nil {
			return cached, nil
		}
	}
	req, _ := http.NewRequest(
		"GET",
		t.BaseURL+"/rest/api/3/issue/"+key+"?expand=changelog",
		nil,
	)
	req.Header.Set("Authorization", "Basic "+t.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return domain.JiraIssueUpdate{}, err
	}
	defer resp.Body.Close()

	var data struct {
		Key       string `json:"key"`
		Changelog struct {
			Histories []struct {
				Created string `json:"created"`
				Author  struct {
					DisplayName string `json:"displayName"`
				} `json:"author"`
				Items []struct {
					Field string `json:"field"`
					From  string `json:"fromString"`
					To    string `json:"toString"`
				} `json:"items"`
			} `json:"histories"`
		} `json:"changelog"`
	}

	json.NewDecoder(resp.Body).Decode(&data)

	updates := []domain.JiraChangelogItem{}

	for _, h := range data.Changelog.Histories {
		for _, item := range h.Items {
			updates = append(updates, domain.JiraChangelogItem{
				Author:    h.Author.DisplayName,
				Field:     item.Field,
				From:      item.From,
				To:        item.To,
				CreatedAt: h.Created,
			})
		}
	}

	result := domain.JiraIssueUpdate{
		IssueKey: data.Key,
		Updates:  updates,
	}

	if t.Cache != nil {
		_ = t.Cache.Set(cacheKey, result, 5*time.Minute)
	}

	return result, nil
}
