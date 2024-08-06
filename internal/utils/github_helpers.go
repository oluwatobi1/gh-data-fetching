package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/models"
)

const (
	baseURL           = "https://api.github.com/repos/%s/commits?per_page=100"
	authHeader        = "Authorization"
	rateLimitErrorMsg = "rate limit exceeded"
	tooManyRequests   = http.StatusTooManyRequests
	successStatus     = http.StatusOK
)

func BuildGHCommitURL(repoName string, config models.CommitConfig) string {
	url := fmt.Sprintf(baseURL, repoName)
	if config.StartDate != "" {
		url += fmt.Sprintf("&since=%s", config.StartDate)
	}
	if config.EndDate != "" {
		url += fmt.Sprintf("&until=%s", config.EndDate)
	}
	if config.Sha != "" {
		url += fmt.Sprintf("&sha=%s", config.Sha)
	}
	return url
}

func FetchBatch(url, token string) ([]models.CommitResponse, string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", err
	}
	if token != "" {
		req.Header.Set(authHeader, "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == tooManyRequests {
		if err := HandleRateLimit(resp); err != nil {
			return nil, "", err
		}
		return nil, "", fmt.Errorf(rateLimitErrorMsg)
	}

	if resp.StatusCode != successStatus {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("failed to fetch commits: %s", string(bodyBytes))
	}

	var commits []models.CommitResponse
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return nil, "", err
	}

	linkHeader := resp.Header.Get("Link")
	links := ParseLinkHeader(linkHeader)
	nextURL := links["next"]

	return commits, nextURL, nil
}
