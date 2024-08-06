package ports

import (
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/models"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/types"
)

type Commit interface {
	Create(commit *models.Commit) error
	FindByHash(hash string) (*models.Commit, error)
	FindByRepoId(repoId uint, page int, pageSize int) ([]*models.Commit, error)
	FindAll() ([]*models.Commit, error)
	CreateMany(commits []models.Commit) error
	Count() (int64, error)
	GetTopCommitAuthors(page int, pageSize int) ([]types.AuthorCommitsCount, error)
	UpsertCommits(commits []models.Commit) error
}

type Repository interface {
	Create(repo *models.Repository) error
	FindByName(name string) (*models.Repository, error)
	FindAll() ([]*models.Repository, error)
	UpdateLastCommitSHA(id uint, sha string) error
}

type GithubService interface {
	FetchRepository(repoName string) (*models.Repository, error)
	FetchCommits(repoName string, repoID uint, config models.CommitConfig) ([]models.Commit, string, int, error)
}
