package app

import (
	"os"

	"github.com/oluwatobi1/gh-api-data-fetch/config"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/adapters/api"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/adapters/db/gorm"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/application/handlers"
	"go.uber.org/zap"
	gm "gorm.io/gorm"
)

func configureRoutes(db *gm.DB, logger *zap.Logger) {
	logger.Info("Configure routes")

	repoRepo := gorm.NewRepository(db)
	commitRepo := gorm.NewCommitRepo(db)
	ghApi := api.NewGitHubAPI(config.Env.GITHUB_TOKEN, logger)
	appHandler := handlers.NewAppHandler(repoRepo, commitRepo, ghApi, logger)
	appHandler.SetupEventBus()

	configureApp(appHandler, logger)

	v1 := router.Group("/api/v1")
	v1.GET("/fetch-repo", appHandler.FetchRepository)
	v1.GET("/list-repo", appHandler.ListRepositories)
	v1.GET("/list-commit", appHandler.ListCommits)
	v1.GET("/fetch-commit", appHandler.UpdateCommit)
	v1.GET("/monitor-commit", appHandler.TriggerMonitorCommits)

}

func configureApp(app *handlers.AppHandler, logger *zap.Logger) {
	logger.Sugar().Info("Configure App")
	logger.Sugar().Info("os.Args", os.Args, config.Env.DEFAULT_REPO)

	if len(os.Args) > 1 && os.Args[1] == "default-repo" {
		if config.Env.DEFAULT_REPO != "" {
			if _, err := app.InitNewRepository(config.Env.DEFAULT_REPO); err != nil {
				logger.Sugar().Warn("Error fetching repositories::: ", err.Error())
			} else {
				logger.Sugar().Info("Repository fetched successfully")
			}
		} else {
			logger.Sugar().Info("Default Repo Not Specified")
		}
	}
}
