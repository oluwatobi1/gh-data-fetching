package gorm_test

import (
	"log"
	"os"
	"testing"

	"github.com/oluwatobi1/gh-api-data-fetch/internal/adapters/db/gorm"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	gm "gorm.io/gorm"
)

var db *gm.DB
var dbFilePath = "test.db"

func setupTestDB() *gm.DB {
	db, _ = gm.Open(sqlite.Open(dbFilePath), &gm.Config{})
	db.AutoMigrate(&models.Commit{})
	return db
}
func teardownTestDB() {
	// Close the database connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get underlying sql.DB: %v", err)
	}

	if err := sqlDB.Close(); err != nil {
		log.Fatalf("failed to close database connection: %v", err)
	}

	// Remove the database file
	if err := os.Remove(dbFilePath); err != nil {
		log.Fatalf("failed to remove database file: %v", err)
	}
}

func TestCreateCommit(t *testing.T) {
	db := setupTestDB()
	repo := gorm.NewCommitRepo(db)
	commit := &models.Commit{Hash: "abc123", RepoID: 1}
	err := repo.Create(commit)
	assert.NoError(t, err)
	teardownTestDB()
}

func TestFindByHash(t *testing.T) {
	db := setupTestDB()
	repo := gorm.NewCommitRepo(db)
	cmt := &models.Commit{Hash: "testHash123", RepoID: 12}
	repo.Create(cmt)
	found, err := repo.FindByHash("testHash123")
	assert.NoError(t, err)
	assert.Equal(t, "testHash123", found.Hash)
	teardownTestDB()
}

func TestFindAll(t *testing.T) {
	db := setupTestDB()
	repo := gorm.NewCommitRepo(db)
	cmt := &models.Commit{Hash: "testHash123"}
	repo.Create(cmt)
	found, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(found))
	teardownTestDB()
}
