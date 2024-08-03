package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/utils"
)

func (h *AppHandler) FetchRepository(gc *gin.Context) {
	repoName := gc.Query("repo")
	if repoName == "" {
		utils.InfoResponse(gc, "Missing repo param", nil, http.StatusBadRequest)
		return
	}

	_, err := h.InitNewRepository(repoName)
	if err != nil {
		utils.InfoResponse(gc, err.
			Error(), nil, http.StatusInternalServerError)
		return
	}

	h.logger.Sugar().Info("::::: AddCommitEvent Emitted for repo:: ", repoName)
	utils.InfoResponse(gc, "success", nil, http.StatusOK)
}

func (h *AppHandler) ListRepositories(gc *gin.Context) {
	repos, err := h.RepositoryRepo.FindAll()
	if err != nil {
		utils.InfoResponse(gc, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	utils.InfoResponse(gc, "success", repos, http.StatusOK)
}

func (h *AppHandler) ListCommits(gc *gin.Context) {

	repos, err := h.CommitRepo.FindAll()
	if err != nil {
		utils.InfoResponse(gc, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	utils.InfoResponse(gc, "commit success", repos, http.StatusOK)
}

func (h *AppHandler) UpdateCommit(gc *gin.Context) {
	err := h.UpdateAllCommits()
	if err != nil {
		utils.InfoResponse(gc, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	utils.InfoResponse(gc, "success", nil, http.StatusOK)
}
