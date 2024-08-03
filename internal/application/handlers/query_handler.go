package handlers

import (
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/types"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/utils"
)

func (h *AppHandler) GetTopCommitAuthors(gc *gin.Context) {
	var req types.PaginationRequest
	if err := gc.ShouldBindQuery(&req); err != nil {
		utils.InfoResponse(gc, err.Error(), nil, http.StatusBadRequest)
		return
	}

	pagination, err := utils.ParsePaginationParams(req.Page, req.PageSize)
	if err != nil {
		utils.InfoResponse(gc, err.Error(), nil, http.StatusBadRequest)
		return
	}
	authors, err := h.CommitRepo.GetTopCommitAuthors(pagination.Page, pagination.PageSize+1)
	if err != nil {
		h.logger.Sugar().Warn("Error fetching top commit authors: ", err)
		return
	}
	hasNext := false
	if len(authors) > pagination.PageSize {
		hasNext = true
	}
	pageLen := int(math.Min(float64(pagination.PageSize), float64(len(authors))))
	resp := types.AuthorCommitsCountResponse{
		Authors: authors[:pageLen],
		Pagination: types.PaginationResponse{
			Page:     fmt.Sprint(pagination.Page),
			PageSize: fmt.Sprint(pageLen),
			HasNext:  hasNext,
		},
	}
	utils.InfoResponse(gc, "success", resp, http.StatusOK)
}

func (h *AppHandler) FetchCommitsByRepoName(gc *gin.Context) {
	var req types.FetchCommitsByRepoNameRequest
	if err := gc.ShouldBindQuery(&req); err != nil {
		utils.InfoResponse(gc, err.Error(), nil, http.StatusBadRequest)
		return
	}
	if req.RepoName == "" {
		utils.InfoResponse(gc, "missing repoName", nil, http.StatusBadRequest)
		return
	}
	pagination, err := utils.ParsePaginationParams(req.Page, req.PageSize)
	if err != nil {
		utils.InfoResponse(gc, err.Error(), nil, http.StatusBadRequest)
		return
	}
	repo, err := h.RepositoryRepo.FindByName(req.RepoName)
	if err != nil {
		h.logger.Sugar().Error("Error finding repository: ", err)
		utils.InfoResponse(gc, err.Error(), nil, http.StatusBadRequest)
		return
	}
	commits, err := h.CommitRepo.FindByRepoId(repo.ID, pagination.Page, pagination.PageSize+1)
	if err != nil {
		h.logger.Sugar().Error("Error fetching commits by: ", err)
		utils.InfoResponse(gc, err.Error(), nil, http.StatusBadRequest)
		return
	}
	hasNext := false
	if len(commits) > pagination.PageSize {
		hasNext = true
	}
	pageLen := int(math.Min(float64(pagination.PageSize), float64(len(commits))))
	resp := types.FetchCommitsByRepoNameResponse{
		Commits: commits[:pageLen],
		Pagination: types.PaginationResponse{
			Page:     fmt.Sprint(pagination.Page),
			PageSize: fmt.Sprint(pageLen),
			HasNext:  hasNext,
		},
	}
	utils.InfoResponse(gc, "success", resp, http.StatusOK)
}
