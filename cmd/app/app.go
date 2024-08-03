package app

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oluwatobi1/gh-api-data-fetch/config"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/adapters/api"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/adapters/db/gorm"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/application/handlers"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/models"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	gm "gorm.io/gorm"
)

type APPServer struct {
}

func NewAPPServer() *APPServer {
	return &APPServer{}
}

var router *gin.Engine

func (s *APPServer) Run() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalln(err)
	}

	db, err := gm.Open(sqlite.Open(config.Env.DB_URL), &gm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	db.AutoMigrate(&models.Repository{}, &models.Commit{})

	logger := zap.Must(zap.NewDevelopment())
	if config.Env.ENVIRONMENT == "release" {
		// production
		logger = zap.Must(zap.NewProduction())
		gin.SetMode(gin.ReleaseMode)
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync()
	router = gin.Default()
	initializeApp(db, logger)
	if err := router.Run(":" + config.Env.PORT); err != nil {
		logger.Sugar().Fatal(err)
	}
}

func initializeApp(db *gm.DB, logger *zap.Logger) {
	logger.Info("initializeApp")
	repoRepo := gorm.NewRepository(db)
	commitRepo := gorm.NewCommitRepo(db)
	ghApi := api.NewGitHubAPI(config.Env.GITHUB_TOKEN, logger)
	appHandler := handlers.NewAppHandler(repoRepo, commitRepo, ghApi, logger)
	appHandler.SetupEventBus()
	setupApp(appHandler, logger)
	configureRoutes(appHandler)
}

func setupApp(app *handlers.AppHandler, logger *zap.Logger) {
	logger.Sugar().Info("setupApp")
	if config.Env.DEFAULT_REPO != "" {
		if _, err := app.InitNewRepository(config.Env.DEFAULT_REPO); err != nil {
			logger.Sugar().Warn("Error fetching repositories::: ", err.Error())
		} else {
			logger.Sugar().Info("Repository fetched successfully")
		}
	} else {
		logger.Sugar().Warn("DEFAULT_REPO Not Specified")
	}

}
