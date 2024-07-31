package ports

import "github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/models"

type Commit interface {
	Create(commit *models.Commit) error
	FindByHash(hash string) (*models.Commit, error)
	FindByRepoId(repoId uint) ([]*models.Commit, error)
}

type Repository interface {
	Create(repo *models.Repository) error
	FindByName(name string) (*models.Repository, error)
	FindAll() ([]*models.Repository, error)
}

type GithubService interface {
	FetchRepository(repoName string) (*models.Repository, error)
	FetchCommits(repoName, startDate, endDate string) ([]models.Commit, error)
}
