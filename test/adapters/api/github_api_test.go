package api_test

import (
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
		w.Write([]byte(`{"id": 1, "full_name": "test/repo", "last_commit_sha": "abc123"}`))
	}))
	defer mockServer.Close()
	logger, _ := zap.NewDevelopment()
	githubApi := api.NewGitHubAPI("mocktoken", logger)

	repo, err := githubApi.FetchRepository("test/repo")
	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.Equal(t, "test/repo", repo.FullName)
}

func TestFetchCommits(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"sha":"abc", "commit": {"author": {"name": "testAuthor"}}}]`))
	}))

	defer mockServer.Close()
	logger, _ := zap.NewDevelopment()
	githubApi := api.NewGitHubAPI("mock_token", logger)
	config := models.CommitConfig{
		StartDate: "2023-01-01T00:00:00Z",
		EndDate:   "2023-12-31T23:59:59Z",
	}
	commits, err := githubApi.FetchCommits("test/repo", 1, config)
	assert.NoError(t, err)
	assert.NotNil(t, commits)
	assert.Equal(t, "abc", commits[0].Hash)
}
