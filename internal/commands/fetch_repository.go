package commands

import (
	"time"

	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/models"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/ports"
)

type FetchRepositoryCmd struct {
	repoRepo ports.Repository
}

func NewFetchRepositoryCmd(repo ports.Repository) *FetchRepositoryCmd {
	return &FetchRepositoryCmd{repoRepo: repo}
}

func (cmd *FetchRepositoryCmd) Execute(repo string) error {
	r := models.NewRepository(repo, "description", "rl", "eng", 2, 2, 3, 4, time.Now(), time.Now())
	return cmd.repoRepo.Create(r)
}
func (cmd *FetchRepositoryCmd) Fetch() ([]*models.Repository, error) {
	return cmd.repoRepo.FindAll()
}
