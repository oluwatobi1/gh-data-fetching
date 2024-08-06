package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/models"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/types"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/ports"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/utils"
	"go.uber.org/zap"
)

type GitHubAPI struct {
	token  string
	logger *zap.Logger
}

func NewGitHubAPI(token string, logger *zap.Logger) ports.GithubService {
	return &GitHubAPI{token: token, logger: logger}
}

func (gh *GitHubAPI) FetchRepository(repoName string) (*models.Repository, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s", repoName)
	req, _ := http.NewRequest("GET", url, nil)
	if gh.token != "" {
		req.Header.Set("Authorization", "Bearer "+gh.token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		gh.logger.Sugar().Warn("FetchRepository Error, " + err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiError types.ApiError
		if err := json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
			return nil, fmt.Errorf("failed to decode error response: %w", err)
		}
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, fmt.Errorf("repository not found: %s", apiError.Message)
		default:
			return nil, fmt.Errorf("failed to fetch repository: %s (status: %d)", apiError.Message, resp.StatusCode)
		}
	}

	var repo models.Repository
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		gh.logger.Sugar().Warn("FetchRepository decode Error, " + err.Error())
		return nil, err
	}
	return &repo, nil
}

func (gh *GitHubAPI) FetchCommits(repoName string, repoId uint, config models.CommitConfig) ([]models.Commit, string, int, error) {
	var allCommits []models.CommitResponse
	var errL error
	var rateLimitDuration int
	url := utils.BuildGHCommitURL(repoName, config)

	gh.logger.Sugar().Info("Fetching Commit in Batches...")
	for len(allCommits) < 1000 {
		commits, nextURL, rL, err := utils.FetchBatch(url, gh.token)
		if err != nil {
			errL = err
			break
		}

		if config.Sha != "" && len(commits) > 0 {
			// remove already fetch hash from hash
			commits = commits[1:]
		}
		allCommits = append(allCommits, commits...)
		if rL != 0 {
			rateLimitDuration = rL
			break
		}
		if len(commits) == 0 || nextURL == "" {
			break
		}
		url = nextURL
	}
	var commitsMd []models.Commit
	for _, cmt := range allCommits {
		commitsMd = append(commitsMd, cmt.ToCommit(repoId))
	}
	gh.logger.Sugar().Info("Total Commits in Current Batch: ", len(commitsMd))

	lastCommitSHA := ""
	if len(allCommits) > 0 {
		lastCommitSHA = allCommits[len(allCommits)-1].SHA
	}

	return commitsMd, lastCommitSHA, rateLimitDuration, errL
}
