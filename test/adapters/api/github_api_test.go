package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oluwatobi1/gh-api-data-fetch/internal/adapters/api"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestFetchRepository(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 1, "full_name": "chromium/chromium", "last_commit_sha": "abc123"}`))
	}))
	defer mockServer.Close()
	logger, _ := zap.NewDevelopment()
	githubApi := api.NewGitHubAPI("", logger)
	repoName := "chromium/chromium"
	repo, err := githubApi.FetchRepository(repoName)
	fmt.Println("repo", repo, "err", err)
	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.Equal(t, repoName, repo.FullName)
}

func TestFetchCommits(t *testing.T) {

	logger, _ := zap.NewDevelopment()
	githubApi := api.NewGitHubAPI("", logger)
	config := models.CommitConfig{
		StartDate: "2023-01-01",
		EndDate:   "2023-12-31",
	}

	repoName := "chromium/chromium"
	commits, _, rl, err := githubApi.FetchCommits(repoName, 1, config)
	assert.NoError(t, err)
	assert.NotNil(t, commits)
	assert.Equal(t, 0, rl)
}
