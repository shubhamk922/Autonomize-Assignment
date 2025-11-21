package user

import (
	"strings"
	"sync"

	"example.com/team-monitoring/domain"
)

type UserIdentityDB struct {
	store map[string]domain.UserIdentity
}

var (
	instance *UserIdentityDB
	once     sync.Once
)

func GetInstance() *UserIdentityDB {
	once.Do(
		func() {
			instance = &UserIdentityDB{
				store: make(map[string]domain.UserIdentity),
			}
		})
	return instance
}

func (db *UserIdentityDB) InitDB() {
	db.store["shubham"] = domain.UserIdentity{
		DisplayName:   "Shubham",
		JiraAccountID: "ShubhamK",
		GitHubHandle:  "shubhamk922",
	}
}

func (db *UserIdentityDB) GetJiraId(name string) string {
	if val, ok := db.store[strings.ToLower(name)]; ok {
		return val.JiraAccountID
	}
	return ""
}

func (db *UserIdentityDB) GetGithubId(name string) string {
	if val, ok := db.store[strings.ToLower(name)]; ok {
		return val.GitHubHandle
	}
	return ""
}
