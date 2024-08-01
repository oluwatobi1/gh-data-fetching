package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/models"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/ports"
)

type AppHandler struct {
	RepositoryRepo ports.Repository
	CommitRepo     ports.Commit
}

func NewAppHandler(repo ports.Repository, cmt ports.Commit) *AppHandler {
	return &AppHandler{
		RepositoryRepo: repo,
		CommitRepo:     cmt,
	}
}

func (h *AppHandler) FetchRepository(w http.ResponseWriter, r *http.Request) {
	repoName := r.URL.Query().Get("repo")
	if repoName == "" {
		http.Error(w, "Missing repo param", http.StatusBadRequest)
		return
	}

	rp := models.Repository{
		Name: repoName,
	}
	if err := h.RepositoryRepo.Create(&rp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rp)

	// Fetch commit
}
