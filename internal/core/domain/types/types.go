package types

import "github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/models"

type Pagination struct {
	Page     int
	PageSize int
}

type AuthorCommitsCount struct {
	Author      string `json:"author"`
	CommitCount int    `json:"commit_count"`
}

type AuthorCommitsCountResponse struct {
	Authors    []AuthorCommitsCount `json:"authors"`
	Pagination PaginationResponse   `json:"pagination"`
}

type FetchCommitsByRepoNameRequest struct {
	RepoName string `form:"repo_name"`
	PaginationRequest
}
type FetchCommitsByRepoNameResponse struct {
	Commits    []*models.Commit   `json:"commits"`
	Pagination PaginationResponse `json:"pagination"`
}

type PaginationRequest struct {
	Page     string `form:"page"`
	PageSize string `form:"page_size"`
}
type PaginationResponse struct {
	Page     string `json:"page"`
	PageSize string `json:"page_size"`
	HasNext  bool   `json:"has_next"`
}

type ApiError struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
	Status           string `json:"status"`
}
