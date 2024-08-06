package gorm

import (
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/models"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/types"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/ports"
	"gorm.io/gorm"
)

type CommitRepo struct {
	db *gorm.DB
}

func NewCommitRepo(db *gorm.DB) ports.Commit {
	return &CommitRepo{db: db}
}

func (c *CommitRepo) Create(commit *models.Commit) error {
	return c.db.Create(commit).Error
}

func (c *CommitRepo) CreateMany(commits []models.Commit) error {
	return c.db.Create(commits).Error
}
func (c *CommitRepo) FindByHash(hash string) (*models.Commit, error) {
	var cmt models.Commit
	if err := c.db.Where("hash = ?", hash).First(&cmt).Error; err != nil {
		return nil, err
	}
	return &cmt, nil
}

func (c *CommitRepo) FindByRepoId(repoId uint, page int, pageSize int) ([]*models.Commit, error) {
	var cmt []*models.Commit
	if err := c.db.Where("repo_id = ?", repoId).
		Limit(pageSize).
		Offset((page-1)*pageSize - 1).
		Find(&cmt).Error; err != nil {
		return nil, err
	}
	return cmt, nil
}

func (r *CommitRepo) FindAll() ([]*models.Commit, error) {
	var cmt []*models.Commit
	if err := r.db.Find(&cmt).Error; err != nil {
		return nil, err
	}
	return cmt, nil
}

func (c *CommitRepo) UpsertCommits(commits []models.Commit) error {
	for _, commit := range commits {
		if err := c.db.Save(&commit).Error; err != nil {
			return err
		}
	}
	return nil
}

// Count returns the total number of commits in the database: for logging purpose
func (c *CommitRepo) Count() (int64, error) {
	var count int64
	if err := c.db.Model(&models.Commit{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (c *CommitRepo) GetTopCommitAuthors(page int, pageSize int) ([]types.AuthorCommitsCount, error) {
	var results []types.AuthorCommitsCount
	err := c.db.Model(&models.Commit{}).
		Select("author, COUNT(*) as commit_count").
		Group("author").
		Order("commit_count DESC").
		Limit(pageSize).Offset((page - 1) * (pageSize - 1)).Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
