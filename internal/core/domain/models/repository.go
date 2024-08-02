package models

import "time"

type Repository struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	FullName        string    `gorm:"unique;not null" json:"full_name"`
	Name            string    `json:"name"`
	Description     string    `gorm:"type:text"  json:"description"`
	URL             string    `gorm:"type:text"`
	Language        string    `json:"language"`
	ForksCount      int       `json:"forks_count"`
	StarsCount      int       `json:"stargazers_count"`
	OpenIssuesCount int       `json:"open_issues"`
	WatchersCount   int       `json:"watchers"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	FetchedAt       time.Time `json:"fetched_at"`
	LastCommitSHA   string    `json:"last_commit_sha"`
}

func NewRepository(full_name, name, description, url, language string, forksCount, starsCount, openIssuesCount, watchersCount int, createdAt, updatedAt time.Time) *Repository {
	return &Repository{
		FullName:        full_name,
		Name:            name,
		Description:     description,
		URL:             url,
		Language:        language,
		ForksCount:      forksCount,
		StarsCount:      starsCount,
		OpenIssuesCount: openIssuesCount,
		WatchersCount:   watchersCount,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		FetchedAt:       time.Now(),
	}
}
