package cmd

import (
	"strings"

	"github.com/jackchuka/gh-oss-watch/services"
)

func HandleDashboard(configService services.ConfigService, githubService services.GitHubService, output services.Output) error {
	config, err := configService.Load()
	if err != nil {
		return err
	}

	if len(config.Repos) == 0 {
		output.Println("No repositories configured. Use 'gh oss-watch add <repo>' to add some.")
		return nil
	}

	output.Println("📊 OSS Watch Dashboard")
	output.Println("======================")

	totalStats := struct {
		Stars  int
		Issues int
		PRs    int
		Forks  int
	}{}

	for _, repoConfig := range config.Repos {
		owner, repo, err := services.ParseRepoString(repoConfig.Repo)
		if err != nil {
			output.Printf("Error parsing repo %s: %v\n", repoConfig.Repo, err)
			continue
		}

		stats, err := githubService.GetRepoStats(owner, repo)
		if err != nil {
			output.Printf("Error fetching stats for %s: %v\n", repoConfig.Repo, err)
			continue
		}

		output.Printf("\n📁 %s\n", repoConfig.Repo)
		output.Printf("   ⭐ Stars: %d\n", stats.Stars)
		output.Printf("   🐛 Issues: %d\n", stats.Issues)
		output.Printf("   🔀 Pull Requests: %d\n", stats.PullRequests)
		output.Printf("   🍴 Forks: %d\n", stats.Forks)
		output.Printf("   📅 Last Updated: %s\n", stats.UpdatedAt.Format("2006-01-02 15:04"))
		output.Printf("   📢 Watching: %s\n", strings.Join(repoConfig.Events, ", "))

		totalStats.Stars += stats.Stars
		totalStats.Issues += stats.Issues
		totalStats.PRs += stats.PullRequests
		totalStats.Forks += stats.Forks
	}

	output.Println("\n📈 Total Across All Repos:")
	output.Printf("   ⭐ Total Stars: %d\n", totalStats.Stars)
	output.Printf("   🐛 Total Issues: %d\n", totalStats.Issues)
	output.Printf("   🔀 Total PRs: %d\n", totalStats.PRs)
	output.Printf("   🍴 Total Forks: %d\n", totalStats.Forks)

	return nil
}
