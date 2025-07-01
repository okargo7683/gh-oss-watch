package cmd

import (
	"fmt"
	"strings"

	"github.com/jackchuka/gh-oss-watch/services"
)

func HandleConfigAdd(repo string, eventArgs []string, configService services.ConfigService, output services.Output) error {
	config, err := configService.Load()
	if err != nil {
		return err
	}

	events := []string{"stars", "issues", "pull_requests", "forks"}
	if len(eventArgs) > 0 {
		events = eventArgs
	}

	err = services.ValidateEvents(events)
	if err != nil {
		return err
	}

	config.AddRepo(repo, events)

	err = configService.Save(config)
	if err != nil {
		return err
	}

	output.Printf("Added %s to watch list with events: %s\n", repo, strings.Join(events, ", "))
	return nil
}

func HandleConfigSet(repo string, eventArgs []string, configService services.ConfigService, output services.Output) error {
	if len(eventArgs) == 0 {
		output.Println("Usage: gh oss-watch set <repo> <events...>")
		output.Println("Available events: stars, issues, pull_requests, forks")
		return fmt.Errorf("no events specified")
	}

	config, err := configService.Load()
	if err != nil {
		return err
	}

	repoConfig := config.GetRepo(repo)
	if repoConfig == nil {
		return fmt.Errorf("repository %s not found in config. Use 'gh oss-watch add' first", repo)
	}

	err = services.ValidateEvents(eventArgs)
	if err != nil {
		return err
	}

	config.AddRepo(repo, eventArgs)

	err = configService.Save(config)
	if err != nil {
		return err
	}

	output.Printf("Updated %s events to: %s\n", repo, strings.Join(eventArgs, ", "))
	return nil
}

func HandleConfigRemove(repo string, configService services.ConfigService, output services.Output) error {
	config, err := configService.Load()
	if err != nil {
		return err
	}

	if !config.RemoveRepo(repo) {
		return fmt.Errorf("repository %s not found in config", repo)
	}

	err = configService.Save(config)
	if err != nil {
		return err
	}

	output.Printf("Removed %s from watch list\n", repo)
	return nil
}
