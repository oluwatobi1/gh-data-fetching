package app

import (
	"github.com/oluwatobi1/gh-api-data-fetch/config"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/adapters/api"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/adapters/db/gorm"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/handlers"
	"go.uber.org/zap"
	gm "gorm.io/gorm"
)

func configureRoutes(db *gm.DB, logger *zap.Logger) {
	logger.Info("Configure routes")

	repoRepo := gorm.NewRepository(db)
	commitRepo := gorm.NewCommitRepo(db)
	ghApi := api.NewGitHubAPI(config.Env.GITHUB_TOKEN, logger)
	repoHandler := handlers.NewAppHandler(repoRepo, commitRepo, ghApi, logger)

	v1 := router.Group("/api/v1")
	v1.GET("/fetch-repo", repoHandler.FetchRepository)
	v1.GET("/list-repo", repoHandler.ListRepositories)
	v1.GET("/list-commit", repoHandler.ListCommits)
	v1.GET("/fetch-commit", repoHandler.UpdateCommit)

}
