package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/models"
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
	req.Header.Set("Authorization", "Bearer "+gh.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		gh.logger.Sugar().Warn("FetchRepository Error, " + err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiError struct {
			Message          string `json:"message"`
			DocumentationURL string `json:"documentation_url"`
			Status           string `json:"status"`
		}

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

// func (gh *GitHubAPI) FetchCommits(repoName, startDate, endDate string, repoId uint) ([]models.Commit, error) {
// 	url := fmt.Sprintf("https://api.github.com/repos/%s/commits?&per_page=2", repoName)
// 	gh.logger.Sugar().Info("urlll:::", url)
// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Set("Authorization", "Bearer "+gh.token)
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		gh.logger.Sugar().Warn("FetchCommits Error, " + err.Error())
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		var apiError struct {
// 			Message          string `json:"message"`
// 			DocumentationURL string `json:"documentation_url"`
// 			Status           string `json:"status"`
// 		}

// 		if err := json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
// 			return nil, fmt.Errorf("failed to decode error response: %w", err)
// 		}
// 		switch resp.StatusCode {
// 		case http.StatusNotFound:
// 			return nil, fmt.Errorf("repository not found: %s", apiError.Message)
// 		default:
// 			return nil, fmt.Errorf("failed to fetch repository: %s (status: %d)", apiError.Message, resp.StatusCode)
// 		}
// 	}
// 	var commitResp []models.CommitResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&commitResp); err != nil {
// 		gh.logger.Sugar().Warn("FetchCommits decode Error, " + err.Error())
// 		return nil, err
// 	}
// 	var commits []models.Commit
// 	for _, cmt := range commitResp {
// 		commits = append(commits, cmt.ToCommit(repoId))
// 	}
// 	return commits, nil

// }
func (gh *GitHubAPI) FetchCommits(repoName string, repoId uint, config models.CommitConfig) ([]models.Commit, error) {
	var allCommits []models.CommitResponse
	url := fmt.Sprintf("https://api.github.com/repos/%s/commits?per_page=20", repoName)

	if config.StartDate != "" {
		url += fmt.Sprintf("&since=%s", config.StartDate)
	}
	if config.EndDate != "" {
		url += fmt.Sprintf("&until=%s", config.EndDate)
	}
	if config.Sha != "" {
		url += fmt.Sprintf("&until=%s", config.Sha)
	}
	count := 2
	for {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+gh.token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		gh.logger.Sugar().Info("Status Code:", resp.StatusCode)

		if resp.StatusCode == http.StatusTooManyRequests {
			if err := utils.HandleRateLimit(resp); err != nil {
				return nil, err
			}
			continue
		}
		if resp.StatusCode != http.StatusOK {

			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			return nil, fmt.Errorf("failed to fetch commits: %s", string(bodyBytes))
		}
		var commits []models.CommitResponse
		if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
			return nil, err
		}

		allCommits = append(allCommits, commits...)

		linkHeader := resp.Header.Get("Link")
		links := utils.ParseLinkHeader(linkHeader)
		nextURL, hasNext := links["next"]
		if !hasNext {
			break
		}
		url = nextURL
		count -= 1
		if count < 1 {
			break
		}
	}

	var commits []models.Commit
	for _, cmt := range allCommits {
		commits = append(commits, cmt.ToCommit(repoId))
	}

	gh.logger.Sugar().Info("Total Commits: ", commits)
	return commits, nil

}
