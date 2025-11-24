package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"example.com/team-monitoring/adapter/in/controllers"
	"example.com/team-monitoring/adapter/out/github"
	"example.com/team-monitoring/adapter/out/jira"
	"example.com/team-monitoring/adapter/out/openai"
	"example.com/team-monitoring/adapter/out/user"
	"example.com/team-monitoring/config"
	"example.com/team-monitoring/infra/cache"
	"example.com/team-monitoring/infra/logger"
	"example.com/team-monitoring/service"
	"example.com/team-monitoring/service/ports/out"
	"example.com/team-monitoring/service/tools"
	githubtool "example.com/team-monitoring/service/tools/github"
	jiratool "example.com/team-monitoring/service/tools/jira"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap/zapcore"
)

func main() {

	cfg, err := config.LoadConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.NewZapLogger(zapcore.InfoLevel, "app.log", "err.log")

	redis := cache.NewRedisClient(cfg.RedisAddr, cfg.RedisPass, cfg.RedisDB)

	userIdentityDB := user.GetInstance()
	userIdentityDB.InitDB()

	jiraClient := jira.New(cfg.JiraURL, cfg.JiraToken)
	jiraClient.Log = logger
	jiraClient.IdentityDB = userIdentityDB
	jiraClient.Cache = redis

	githubClient := github.New(cfg.GithubToken)
	githubClient.Log = logger
	githubClient.IdentityDB = userIdentityDB
	githubClient.Cache = redis

	aiClient := openai.NewOpenAIClient("")
	aiClient.Log = logger

	registry := service.NewToolRegistry()
	registry.Register(&githubtool.GetUserCommitsTool{Github: githubClient, Log: logger})
	registry.Register(&githubtool.GetUserPRsTool{Github: githubClient, Log: logger})
	registry.Register(&githubtool.GetUserContributedReposTool{Github: githubClient, Log: logger})
	registry.Register(&jiratool.GetUserIssuesTool{Jira: jiraClient, Log: logger})
	registry.Register(&jiratool.GetIssueStatusTool{Jira: jiraClient, Log: logger})
	registry.Register(&jiratool.GetIssueUpdatesTool{Jira: jiraClient, Log: logger})
	registry.Register(&tools.GetMemberActivityTool{Jira: jiraClient, Github: githubClient,
		Steps: map[string]out.AITool{
			"get_user_issues":  &jiratool.GetUserIssuesTool{Jira: jiraClient, Log: logger},
			"get_user_commits": &githubtool.GetUserCommitsTool{Github: githubClient, Log: logger},
		},
		Log: logger,
	})

	bot := service.NewBot(aiClient, registry)

	http.HandleFunc("/chat", (&controllers.ChatController{
		Service: bot,
	}).Handle)
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("🚀 Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
