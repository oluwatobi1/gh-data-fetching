package models

import "time"

type Commit struct {
	ID        uint   `gorm:"primaryKey"`
	RepoID    uint   `gorm:"index;not null"`
	Hash      string `gorm:"unique;not null"`
	Message   string `gorm:"type:text"`
	Author    string
	Date      time.Time
	URL       string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCommit(repoID uint, hash, message, author, url string, date time.Time) *Commit {
	return &Commit{
		RepoID:    repoID,
		Hash:      hash,
		Message:   message,
		Author:    author,
		Date:      date,
		URL:       url,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
