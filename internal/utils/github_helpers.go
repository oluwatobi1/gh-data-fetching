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

func FetchBatch(url, token string) ([]models.CommitResponse, string, int, error) {
	var rL int
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", rL, err
	}
	if token != "" {
		req.Header.Set(authHeader, "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", rL, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == tooManyRequests {
		rateLimit, err := HandleRateLimit(resp)
		if err != nil {
			return nil, "", rateLimit, err
		}
		return nil, "", rL, fmt.Errorf(rateLimitErrorMsg)
	}

	if resp.StatusCode != successStatus {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, "", rL, fmt.Errorf("failed to fetch commits: %s", string(bodyBytes))
	}

	var commits []models.CommitResponse
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return nil, "", rL, err
	}

	linkHeader := resp.Header.Get("Link")
	links := ParseLinkHeader(linkHeader)
	nextURL := links["next"]

	return commits, nextURL, rL, nil
}
