package main

import (
	"fmt"
	"os"

	"github.com/jackchuka/gh-oss-watch/cmd"
	"github.com/jackchuka/gh-oss-watch/services"
)

func main() {
	configService := services.NewConfigService()
	cacheService := services.NewCacheService()
	output := services.NewConsoleOutput()

	githubService, err := services.NewConcurrentGitHubService()
	if err != nil {
		fmt.Printf("Error creating GitHub service: %v\n", err)
		os.Exit(1)
	}

	cli := cmd.NewCLI(configService, cacheService, githubService, output)
	cli.Run(os.Args)
}
