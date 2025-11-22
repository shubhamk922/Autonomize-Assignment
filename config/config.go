package config

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type AppConfig struct {
	RedisAddr   string
	RedisPass   string
	RedisDB     int
	JiraToken   string
	JiraURL     string
	GithubToken string
}

func LoadConfig(ctx context.Context) (*AppConfig, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	client := ssm.NewFromConfig(cfg)

	get := func(name string, secure bool) (string, error) {
		resp, err := client.GetParameter(ctx, &ssm.GetParameterInput{
			Name:           aws.String(name),
			WithDecryption: &secure,
		})
		if err != nil {
			return "", err
		}
		return *resp.Parameter.Value, nil
	}

	redisAddr, _ := get("/app/redis/addr", false)
	redisPass, _ := get("/app/redis/pass", true)
	redisDB, _ := get("/app/redis/db", false)
	redisD, _ := strconv.Atoi(redisDB)
	jiraUrl, _ := get("/app/jira/url", false)
	jiraToken, _ := get("/app/jira/token", true)
	githubToken, _ := get("/app/github/token", true)

	return &AppConfig{
		RedisAddr:   redisAddr,
		RedisPass:   redisPass,
		RedisDB:     redisD,
		JiraToken:   jiraToken,
		JiraURL:     jiraUrl,
		GithubToken: githubToken,
	}, nil
}
