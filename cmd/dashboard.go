package cmd

import (
	"strings"

	"github.com/jackchuka/gh-oss-watch/services"
)

type dashboardProcessor struct {
	output     services.Output
	totalStats *struct {
		Stars  int
		Issues int
		PRs    int
		Forks  int
	}
}

func (d *dashboardProcessor) ProcessRepo(repoConfig services.RepoConfig, stats *services.RepoStats, index int) error {
	d.output.Printf("\n📁 %s\n", repoConfig.Repo)
	d.output.Printf("   ⭐ Stars: %d\n", stats.Stars)
	d.output.Printf("   🐛 Issues: %d\n", stats.Issues)
	d.output.Printf("   🔀 Pull Requests: %d\n", stats.PullRequests)
	d.output.Printf("   🍴 Forks: %d\n", stats.Forks)
	d.output.Printf("   📅 Last Updated: %s\n", stats.UpdatedAt.Format("2006-01-02 15:04"))
	d.output.Printf("   📢 Watching: %s\n", strings.Join(repoConfig.Events, ", "))

	d.totalStats.Stars += stats.Stars
	d.totalStats.Issues += stats.Issues
	d.totalStats.PRs += stats.PullRequests
	d.totalStats.Forks += stats.Forks

	return nil
}

func (c *CLI) handleDashboard() error {
	config, err := c.validateConfig()
	if err != nil {
		return err
	}

	if len(config.Repos) == 0 {
		return nil
	}

	c.output.Println("📊 OSS Watch Dashboard")
	c.output.Println("======================")

	totalStats := struct {
		Stars  int
		Issues int
		PRs    int
		Forks  int
	}{}

	processor := &dashboardProcessor{
		output:     c.output,
		totalStats: &totalStats,
	}

	err = c.processReposWithBatch(config, processor)
	if err != nil {
		return err
	}

	c.output.Println("\n📈 Total Across All Repos:")
	c.output.Printf("   ⭐ Total Stars: %d\n", totalStats.Stars)
	c.output.Printf("   🐛 Total Issues: %d\n", totalStats.Issues)
	c.output.Printf("   🔀 Total PRs: %d\n", totalStats.PRs)
	c.output.Printf("   🍴 Total Forks: %d\n", totalStats.Forks)

	return nil
}
