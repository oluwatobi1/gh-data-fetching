package gorm

import (
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/models"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/ports"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) ports.Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(repo *models.Repository) error {
	return r.db.Create(repo).Error
}

func (r *Repository) FindAll() ([]*models.Repository, error) {
	var repos []*models.Repository
	if err := r.db.Find(&repos).Error; err != nil {
		return nil, err
	}
	return repos, nil
}

func (r *Repository) FindByName(name string) (*models.Repository, error) {
	var repo *models.Repository
	if err := r.db.Where("full_name = ?", name).First(&repo).Error; err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *Repository) UpdateLastCommitSHA(id uint, sha string) error {
	if err := r.db.Model(&models.Repository{}).
		Where("id = ?", id).
		Update("last_commit_sha", sha).Error; err != nil {
		return err

	}
	return nil
}
