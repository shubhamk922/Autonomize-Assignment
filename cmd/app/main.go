package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"example.com/team-monitoring/adapter/out/github"
	"example.com/team-monitoring/adapter/out/jira"
	"example.com/team-monitoring/adapter/out/openai"
	"example.com/team-monitoring/service"
	"example.com/team-monitoring/service/tools"
)

func main() {
	jiraClient := jira.New(os.Getenv("JIRA_URL"), os.Getenv("JIRA_TOKEN"))
	githubClient := github.New(os.Getenv("GITHUB_TOKEN"))
	aiClient := openai.NewOpenAIClient("")

	activityservice := service.NewActivityService(jiraClient, githubClient, aiClient)

	registry := service.NewToolRegistry()

	registry.Register(&tools.GetUserCommitsTool{Github: githubClient})
	registry.Register(&tools.GetUserPRsTool{Github: githubClient})
	registry.Register(&tools.GetUserContributedReposTool{Github: githubClient})
	registry.Register(&tools.GetMemberActivityTool{Svc: activityservice})

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
