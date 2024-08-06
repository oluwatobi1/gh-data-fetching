package events

import "github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/models"

type Event interface {
	EventType() string
}

type AddCommitEvent struct {
	Repo   *models.Repository
	Config models.CommitConfig
}

type StartMonitorEvent struct {
}

func (e AddCommitEvent) EventType() string {
	return "AddCommitEvent"
}

func (e StartMonitorEvent) EventType() string {
	return "StartMonitorEvent"
}
