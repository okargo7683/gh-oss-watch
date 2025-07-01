package cmd

import (
	"fmt"
	"os"

	"github.com/jackchuka/gh-oss-watch/services"
)

type CLI struct {
	configService services.ConfigService
	cacheService  services.CacheService
	githubService services.GitHubService
	output        services.Output
}

func NewCLI(configService services.ConfigService, cacheService services.CacheService, githubService services.GitHubService, output services.Output) *CLI {
	return &CLI{
		configService: configService,
		cacheService:  cacheService,
		githubService: githubService,
		output:        output,
	}
}

func (c *CLI) Run(args []string) {
	if len(args) < 2 {
		c.printUsage()
		return
	}

	command := args[1]
	cmdArgs := args[2:]

	var err error

	switch command {
	case "init":
		err = HandleInit(c.configService, c.output)
	case "add":
		err = c.handleAddCommand(cmdArgs)
	case "set":
		err = c.handleSetCommand(cmdArgs)
	case "remove":
		err = c.handleRemoveCommand(cmdArgs)
	case "status":
		err = HandleStatus(c.configService, c.cacheService, c.githubService, c.output)
	case "dashboard":
		err = HandleDashboard(c.configService, c.githubService, c.output)
	default:
		c.output.Printf("Unknown command: %s\n", command)
		c.printUsage()
		os.Exit(1)
	}

	if err != nil {
		c.output.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func (c *CLI) handleAddCommand(args []string) error {
	if len(args) < 1 {
		c.output.Println("Usage: gh oss-watch add <repo> [events...]")
		return fmt.Errorf("repository required")
	}
	return HandleConfigAdd(args[0], args[1:], c.configService, c.output)
}

func (c *CLI) handleSetCommand(args []string) error {
	if len(args) < 2 {
		c.output.Println("Usage: gh oss-watch set <repo> <events...>")
		return fmt.Errorf("repository and events required")
	}
	return HandleConfigSet(args[0], args[1:], c.configService, c.output)
}

func (c *CLI) handleRemoveCommand(args []string) error {
	if len(args) < 1 {
		c.output.Println("Usage: gh oss-watch remove <repo>")
		return fmt.Errorf("repository required")
	}
	return HandleConfigRemove(args[0], c.configService, c.output)
}

func (c *CLI) printUsage() {
	c.output.Println("gh-oss-watch - GitHub CLI plugin for OSS maintainers")
	c.output.Println("")
	c.output.Println("Usage:")
	c.output.Println("  gh oss-watch init                    Initialize config file")
	c.output.Println("  gh oss-watch add <repo> [events...]  Add repo to watch list")
	c.output.Println("  gh oss-watch set <repo> <events...>  Configure events for repo")
	c.output.Println("  gh oss-watch remove <repo>           Remove repo from watch list")
	c.output.Println("  gh oss-watch status                  Show new activity")
	c.output.Println("  gh oss-watch dashboard               Show summary across all repos")
}
