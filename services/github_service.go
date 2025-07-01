package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
)

type GitHubServiceImpl struct {
	client *api.RESTClient
}

func NewGitHubService() (GitHubService, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return nil, err
	}
	return &GitHubServiceImpl{client: client}, nil
}

func (g *GitHubServiceImpl) GetRepoStats(owner, repo string) (*RepoStats, error) {
	repoPath := fmt.Sprintf("repos/%s/%s", owner, repo)

	var repoData struct {
		Name  string `json:"name"`
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
		StargazersCount int       `json:"stargazers_count"`
		ForksCount      int       `json:"forks_count"`
		OpenIssuesCount int       `json:"open_issues_count"`
		UpdatedAt       time.Time `json:"updated_at"`
	}

	err := g.client.Get(repoPath, &repoData)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repo %s/%s: %w", owner, repo, err)
	}

	pullRequestsPath := fmt.Sprintf("repos/%s/%s/pulls?state=open", owner, repo)
	var prs []struct{}
	err = g.client.Get(pullRequestsPath, &prs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PRs for %s/%s: %w", owner, repo, err)
	}

	return &RepoStats{
		Name:         repoData.Name,
		Owner:        repoData.Owner.Login,
		Stars:        repoData.StargazersCount,
		Issues:       repoData.OpenIssuesCount,
		PullRequests: len(prs),
		Forks:        repoData.ForksCount,
		UpdatedAt:    repoData.UpdatedAt,
	}, nil
}

func (g *GitHubServiceImpl) GetCurrentUser() (string, error) {
	var user struct {
		Login string `json:"login"`
	}
	err := g.client.Get("user", &user)
	if err != nil {
		return "", err
	}
	return user.Login, nil
}

func ParseRepoString(repoStr string) (owner, repo string, err error) {
	parts := strings.Split(repoStr, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid repo format: %s (expected owner/repo)", repoStr)
	}
	return parts[0], parts[1], nil
}

func CalculateEventSummary(repoStr string, current *RepoStats, previous RepoState) EventSummary {
	summary := EventSummary{
		Repo: repoStr,
	}

	if current.Stars > previous.LastStarCount {
		summary.NewStars = current.Stars - previous.LastStarCount
		summary.HasChanges = true
	}

	if current.Issues > previous.LastIssueCount {
		summary.NewIssues = current.Issues - previous.LastIssueCount
		summary.HasChanges = true
	}

	if current.PullRequests > previous.LastPRCount {
		summary.NewPRs = current.PullRequests - previous.LastPRCount
		summary.HasChanges = true
	}

	if current.Forks > previous.LastForkCount {
		summary.NewForks = current.Forks - previous.LastForkCount
		summary.HasChanges = true
	}

	return summary
}
