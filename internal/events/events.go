package events

import "github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/models"

type AddCommitEvent struct {
	Repo   *models.Repository
	Config models.CommitConfig
}

type StartMonitorEvent struct {
}
