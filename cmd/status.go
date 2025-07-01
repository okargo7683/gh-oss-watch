package cmd

import (
	"time"

	"github.com/jackchuka/gh-oss-watch/services"
)

func HandleStatus(configService services.ConfigService, cacheService services.CacheService, githubService services.GitHubService, output services.Output) error {
	config, err := configService.Load()
	if err != nil {
		return err
	}

	if len(config.Repos) == 0 {
		output.Println("No repositories configured. Use 'gh oss-watch add <repo>' to add some.")
		return nil
	}

	cache, err := cacheService.Load()
	if err != nil {
		return err
	}

	hasChanges := false

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

		previousState, exists := cache.Repos[repoConfig.Repo]
		if !exists {
			previousState = services.RepoState{}
		}

		summary := services.CalculateEventSummary(repoConfig.Repo, stats, previousState)

		if summary.HasChanges {
			hasChanges = true
			output.Printf("\n📈 %s:\n", repoConfig.Repo)

			for _, event := range repoConfig.Events {
				switch event {
				case "stars":
					if summary.NewStars > 0 {
						output.Printf("  ⭐ +%d stars (%d total)\n", summary.NewStars, stats.Stars)
					}
				case "issues":
					if summary.NewIssues > 0 {
						output.Printf("  🐛 +%d issues (%d open)\n", summary.NewIssues, stats.Issues)
					}
				case "pull_requests":
					if summary.NewPRs > 0 {
						output.Printf("  🔀 +%d pull requests (%d open)\n", summary.NewPRs, stats.PullRequests)
					}
				case "forks":
					if summary.NewForks > 0 {
						output.Printf("  🍴 +%d forks (%d total)\n", summary.NewForks, stats.Forks)
					}
				}
			}
		}

		cache.Repos[repoConfig.Repo] = services.RepoState{
			LastStarCount:  stats.Stars,
			LastIssueCount: stats.Issues,
			LastPRCount:    stats.PullRequests,
			LastForkCount:  stats.Forks,
			LastUpdated:    stats.UpdatedAt,
		}
	}

	if !hasChanges {
		output.Println("No new activity since last check.")
	}

	cache.LastCheck = time.Now()
	err = cacheService.Save(cache)
	if err != nil {
		output.Printf("Warning: Error saving cache: %v\n", err)
	}

	return nil
}
