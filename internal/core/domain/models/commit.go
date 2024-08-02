package models

import "time"

type Commit struct {
	ID          uint      `gorm:"primaryKey"`
	RepoID      uint      `gorm:"index;not null"`
	Hash        string    `gorm:"unique;not null" json:"sha"`
	Message     string    `gorm:"type:text" json:"message"`
	Author      string    `json:"author"`
	AuthorEmail string    `json:"author_email"`
	Date        time.Time `json:"author_date"`
	URL         string    `gorm:"type:text" json:"url"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
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

type CommitResponse struct {
	SHA    string `json:"sha"`
	NodeID string `json:"node_id"`
	Commit struct {
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Message      string `json:"message"`
		URL          string `json:"url"`
		CommentCount int    `json:"comment_count"`
	} `json:"commit"`
	URL string `json:"url"`
}

func (c *CommitResponse) ToCommit(repoId uint) Commit {
	now := time.Now()

	return Commit{
		RepoID:      repoId,
		Hash:        c.SHA,
		Message:     c.Commit.Message,
		Author:      c.Commit.Author.Name,
		AuthorEmail: c.Commit.Author.Email,
		Date:        c.Commit.Author.Date,
		URL:         c.Commit.URL,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

type CommitConfig struct {
	StartDate string
	EndDate   string
	Sha       string
}
