package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/models"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/ports"
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
	req.Header.Set("Authorization", "Bearer "+gh.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		gh.logger.Sugar().Warn("FetchRepository Error, " + err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	var repo models.Repository
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		gh.logger.Sugar().Warn("FetchRepository decode Error, " + err.Error())
		return nil, err
	}
	return &repo, nil
}

func (gh *GitHubAPI) FetchCommits(repoName, startDate, endDate string) ([]models.Commit, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/commits?since=%s&until=%s", repoName, startDate, endDate)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+gh.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		gh.logger.Sugar().Warn("FetchCommits Error, " + err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	var commits []models.Commit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		gh.logger.Sugar().Warn("FetchCommits decode Error, " + err.Error())
		return nil, err
	}
	return commits, nil

}
