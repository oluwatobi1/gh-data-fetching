package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/ports"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/utils"
	"go.uber.org/zap"
)

type AppHandler struct {
	RepositoryRepo ports.Repository
	CommitRepo     ports.Commit
	GithubService  ports.GithubService
}

func NewAppHandler(repo ports.Repository, cmt ports.Commit, gh ports.GithubService, logger *zap.Logger) *AppHandler {
	return &AppHandler{
		RepositoryRepo: repo,
		CommitRepo:     cmt,
		GithubService:  gh,
	}
}

func (h *AppHandler) FetchRepository(gc *gin.Context) {
	repoName := gc.Query("repo")
	if repoName == "" {
		utils.InfoResponse(gc, "Missing repo param", nil, http.StatusBadRequest)
		return
	}

	repoMeta, err := h.GithubService.FetchRepository(repoName)
	if err != nil {
		utils.InfoResponse(gc, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	// rp := models.Repository{
	// 	Name: repoName,
	// }
	// if err := h.RepositoryRepo.Create(&rp); err != nil {
	// 	utils.InfoResponse(gc, err.Error(), nil, http.StatusInternalServerError)
	// 	return
	// }

	utils.InfoResponse(gc, "success", repoMeta, http.StatusOK)

}

func (h *AppHandler) ListRepositories(gc *gin.Context) {
	repos, err := h.RepositoryRepo.FindAll()
	if err != nil {
		utils.InfoResponse(gc, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	utils.InfoResponse(gc, "success", repos, http.StatusOK)

}
