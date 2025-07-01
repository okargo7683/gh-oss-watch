package services

//go:generate go tool mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/mock_$GOFILE

import "time"

type ConfigService interface {
	Load() (*Config, error)
	Save(config *Config) error
	GetConfigPath() (string, error)
}

type CacheService interface {
	Load() (*CacheData, error)
	Save(cache *CacheData) error
}

type GitHubService interface {
	GetRepoStats(owner, repo string) (*RepoStats, error)
	GetCurrentUser() (string, error)
}

type Output interface {
	Printf(format string, args ...any)
	Println(args ...any)
}

type Config struct {
	Repos []RepoConfig `yaml:"repos"`
}

type RepoConfig struct {
	Repo   string   `yaml:"repo"`
	Events []string `yaml:"events"`
}

type CacheData struct {
	LastCheck time.Time            `yaml:"last_check"`
	Repos     map[string]RepoState `yaml:"repos"`
}

type RepoState struct {
	LastStarCount  int       `yaml:"last_star_count"`
	LastIssueCount int       `yaml:"last_issue_count"`
	LastPRCount    int       `yaml:"last_pr_count"`
	LastForkCount  int       `yaml:"last_fork_count"`
	LastUpdated    time.Time `yaml:"last_updated"`
}

type RepoStats struct {
	Name         string
	Owner        string
	Stars        int
	Issues       int
	PullRequests int
	Forks        int
	UpdatedAt    time.Time
}

type EventSummary struct {
	Repo       string
	NewStars   int
	NewIssues  int
	NewPRs     int
	NewForks   int
	HasChanges bool
}

func (c *Config) AddRepo(repo string, events []string) {
	for i, r := range c.Repos {
		if r.Repo == repo {
			c.Repos[i].Events = events
			return
		}
	}
	c.Repos = append(c.Repos, RepoConfig{
		Repo:   repo,
		Events: events,
	})
}

func (c *Config) GetRepo(repo string) *RepoConfig {
	for _, r := range c.Repos {
		if r.Repo == repo {
			return &r
		}
	}
	return nil
}

func (c *Config) RemoveRepo(repo string) bool {
	for i, r := range c.Repos {
		if r.Repo == repo {
			c.Repos = append(c.Repos[:i], c.Repos[i+1:]...)
			return true
		}
	}
	return false
}
