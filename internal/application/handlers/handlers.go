package handlers

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oluwatobi1/gh-api-data-fetch/config"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/models"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/ports"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/events"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/utils"
	"go.uber.org/zap"
)

type AppHandler struct {
	RepositoryRepo    ports.Repository
	CommitRepo        ports.Commit
	GithubService     ports.GithubService
	EventBus          *events.EventBus
	logger            *zap.Logger
	monitoringRunning bool
}

func NewAppHandler(repo ports.Repository, cmt ports.Commit, gh ports.GithubService, logger *zap.Logger) *AppHandler {
	return &AppHandler{
		RepositoryRepo: repo,
		CommitRepo:     cmt,
		GithubService:  gh,
		logger:         logger,
	}
}

func (h *AppHandler) SetupEventBus() {
	eventBus := events.NewEventBus()
	eventBus.Register(reflect.TypeOf(events.AddCommitEvent{}), h.HandleAddCommitEvent)
	eventBus.Register(reflect.TypeOf(events.StartMonitorEvent{}), h.HandleStartMonitoringEvent)
	h.EventBus = eventBus
}

func (h *AppHandler) StartMonitoring() {
	h.monitoringRunning = true
}

func (h *AppHandler) StopMonitoring() {
	h.monitoringRunning = false
}

// Method to check if monitoring is already running
func (h *AppHandler) isMonitoringRunning() bool {
	return h.monitoringRunning
}

// Add a new repository to be pull and monitored
func (h *AppHandler) InitNewRepository(repoName string) (bool, error) {
	if repoName == "" {
		return false, fmt.Errorf("missing repo name")
	}

	repoMeta, err := h.GithubService.FetchRepository(repoName)
	if err != nil {
		return false, err
	}

	if err := h.RepositoryRepo.Create(repoMeta); err != nil {
		// todo: add specific check for already exist error
		h.logger.Sugar().Info("err:", err.Error())
	}
	cmtConfig := models.CommitConfig{
		StartDate: config.Env.START_DATE,
		EndDate:   config.Env.END_DATE,
	}
	h.EventBus.Emit(events.AddCommitEvent{Repo: repoMeta, Config: cmtConfig})
	h.logger.Sugar().Info("::::: AddCommitEvent Emitted for repo:: ", repoMeta.FullName)
	return true, nil
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
		cmtConfig := models.CommitConfig{
			StartDate: config.Env.START_DATE,
			EndDate:   config.Env.END_DATE,
			Sha:       repo.LastCommitSHA,
		}
		h.EventBus.Emit(events.AddCommitEvent{Repo: repo, Config: cmtConfig})
		h.logger.Sugar().Info("::::: AddCommitEvent Emitted for repo:: ", repo.FullName)
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
	ticker := time.NewTicker(1 * time.Minute)
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

func (h *AppHandler) HandleAddCommitEvent(event events.AddCommitEvent) {
	repo := event.Repo
	config := event.Config
	h.logger.Sugar().Info("Received AddCommitEvent repo:: ", repo.FullName)
	if err := h.CommitManager(repo, config); err != nil {
		h.logger.Sugar().Error("CommitManager error: ", err)
		return
	}
	if !h.isMonitoringRunning() {
		h.EventBus.Emit(events.StartMonitorEvent{})
		h.logger.Sugar().Info("::::::: StartMonitorEvent Emitted for repo:: ", repo.FullName)
		h.StartMonitoring()
	}
}

func (h *AppHandler) HandleStartMonitoringEvent(event events.StartMonitorEvent) {
	h.logger.Sugar().Info("Received StartMonitorEvent Emitted for repo:: ")
	go h.MonitorCommits()
	h.logger.Sugar().Info("Started Monitoring all repos")
}
