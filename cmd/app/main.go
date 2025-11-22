package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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
	"go.uber.org/zap/zapcore"
)

func main() {

	cfg, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.NewZapLogger(zapcore.InfoLevel, "app.log", "err.log")

	redis := cache.NewRedisClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)

	userIdentityDB := user.GetInstance()
	userIdentityDB.InitDB()

	jiraClient := jira.New(os.Getenv("JIRA_URL"), os.Getenv("JIRA_TOKEN"))
	jiraClient.Log = logger
	jiraClient.IdentityDB = userIdentityDB
	jiraClient.Cache = redis

	githubClient := github.New(os.Getenv("GITHUB_TOKEN"))
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

	reader := bufio.NewReader(os.Stdin)
	ctx := context.Background()
	for {
		fmt.Print("\nAsk: ")
		q, _ := reader.ReadString('\n')
		q = strings.TrimSpace(q)
		resp, err := bot.Handle(ctx, q)
		if err != nil {
			fmt.Printf("Error %+v", err)
			continue
		}
		fmt.Println("\n➡ ", resp)
	}
}
