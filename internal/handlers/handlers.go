package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oluwatobi1/gh-api-data-fetch/config"
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
	// repoId := gc.Query("repo")
	// if repoId == "" {
	// 	utils.InfoResponse(gc, "Missing repo id param", nil, http.StatusBadRequest)
	// 	return
	// }

	repos, err := h.CommitRepo.FindAll()
	if err != nil {
		utils.InfoResponse(gc, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	utils.InfoResponse(gc, "commit success", repos, http.StatusOK)

}

func (h *AppHandler) UpdateCommit(gc *gin.Context) {
	err := h.CommitHandler()
	if err != nil {
		utils.InfoResponse(gc, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	utils.InfoResponse(gc, "success", nil, http.StatusOK)
}

func (h *AppHandler) CommitHandler() error {
	h.logger.Sugar().Info("CommitHandler ")
	repos, err := h.RepositoryRepo.FindAll()
	if err != nil {
		return err
	}
	if len(repos) < 1 {
		return fmt.Errorf("no repository added yet. add repo to fetch commits")
	}
	for _, repo := range repos {
		h.logger.Sugar().Info("Fetching Repo Commit:: ", repo.FullName)
		commits, err := h.GithubService.FetchCommits(repo.FullName, config.Env.START_DATE, config.Env.END_DATE, repo.ID)
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
		h.logger.Sugar().Info("commit load completed")
	}
	return nil
}
