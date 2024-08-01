package app

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oluwatobi1/gh-api-data-fetch/config"
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
	configureRoutes(db, logger)

	if err := router.Run(":" + config.Env.PORT); err != nil {
		logger.Sugar().Fatal(err)
	}
}
