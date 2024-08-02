package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oluwatobi1/gh-api-data-fetch/config"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/models"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/ports"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/utils"
	"go.uber.org/zap"
)

type AppHandler struct {
	RepositoryRepo ports.Repository
	CommitRepo     ports.Commit
	GithubService  ports.GithubService
	logger         *zap.Logger
}

func NewAppHandler(repo ports.Repository, cmt ports.Commit, gh ports.GithubService, logger *zap.Logger) *AppHandler {
	return &AppHandler{
		RepositoryRepo: repo,
		CommitRepo:     cmt,
		GithubService:  gh,
		logger:         logger,
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

	if err := h.RepositoryRepo.Create(repoMeta); err != nil {
		utils.InfoResponse(gc, err.Error(), nil, http.StatusInternalServerError)
		return
	}
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

func (h *AppHandler) UpdateAllCommits() error {
	repos, err := h.RepositoryRepo.FindAll()
	if err != nil {
		return err
	}
	if len(repos) < 1 {
		return fmt.Errorf("no repository added yet. add repo to fetch commits")
	}
	for _, repo := range repos {
		h.logger.Sugar().Info("Fetching Repo Commit:: ", repo.FullName)
		cmtConfig := models.CommitConfig{
			StartDate: config.Env.START_DATE,
			EndDate:   config.Env.END_DATE,
			Sha:       repo.LastCommitSHA,
		}
		err := h.CommitManager(repo, cmtConfig)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *AppHandler) CommitManager(repo *models.Repository, config models.CommitConfig) error {
	commits, err := h.GithubService.FetchCommits(repo.FullName, repo.ID, config)
	if err != nil {
		return err
	}
	for _, commit := range commits {
		if _, err := h.CommitRepo.FindByHash(commit.Hash); err != nil {
			commit.RepoID = repo.ID
			if err := h.CommitRepo.Create(&commit); err != nil {
				return err
			}
			h.logger.Sugar().Info("saved commit ", commit.Hash)
		} else {
			h.logger.Sugar().Info("commit already exist ", commit.Hash)
		}
	}
	if len(commits) > 0 {
		if err := h.RepositoryRepo.UpdateLastCommitSHA(repo.ID, commits[len(commits)-1].Hash); err != nil {
			h.logger.Sugar().Error("error updating commit Sha  ", err)
		}
	}
	return nil
}

func (h *AppHandler) TriggerMonitorCommits(gc *gin.Context) {
	go h.MonitorCommits()
	utils.InfoResponse(gc, "Commit monitoring started", nil, http.StatusOK)

}
func (h *AppHandler) MonitorCommits() {
	h.logger.Sugar().Info("MonitorCommits")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			h.logger.Sugar().Info("Starting hourly commit update check")
			err := h.UpdateAllCommits()
			if err != nil {
				h.logger.Sugar().Error("Error updating commits: ", err)

			} else {
				h.logger.Sugar().Info("Successfully updated commits")

			}

		}
	}
}
