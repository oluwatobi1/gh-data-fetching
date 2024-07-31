package models

import "time"

type Repository struct {
	ID              uint   `gorm:"primaryKey"`
	Name            string `gorm:"unique;not null"`
	Description     string `gorm:"type:text"`
	URL             string `gorm:"type:text"`
	Language        string
	ForksCount      int
	StarsCount      int
	OpenIssuesCount int
	WatchersCount   int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	FetchedAt       time.Time
}

func NewRepository(name, description, url, language string, forksCount, starsCount, openIssuesCount, watchersCount int, createdAt, updatedAt time.Time) *Repository {
	return &Repository{
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
