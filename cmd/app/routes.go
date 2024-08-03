package app

import (
	"github.com/oluwatobi1/gh-api-data-fetch/internal/application/handlers"
)

func configureRoutes(appHandler *handlers.AppHandler) {
	v1 := router.Group("/api/v1")
	v1.GET("/fetch-repo", appHandler.FetchRepository)
	v1.GET("/top-commit-authors", appHandler.GetTopCommitAuthors)
	v1.GET("/commits", appHandler.FetchCommitsByRepoName)
	// v1.GET("/list-repo", appHandler.ListRepositories)
	// v1.GET("/list-commit", appHandler.ListCommits)

}
