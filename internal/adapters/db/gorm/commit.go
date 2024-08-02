package gorm

import (
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/models"
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

func (c *CommitRepo) FindByHash(hash string) (*models.Commit, error) {
	var cmt models.Commit
	if err := c.db.Where("hash = ?", hash).First(&cmt).Error; err != nil {
		return nil, err
	}
	return &cmt, nil
}

func (c *CommitRepo) FindByRepoId(repoId uint) ([]*models.Commit, error) {
	var cmt []*models.Commit
	if err := c.db.Where("repo_id = ?", repoId).Find(&cmt).Error; err != nil {
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
