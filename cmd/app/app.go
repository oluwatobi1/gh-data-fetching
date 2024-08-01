package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/oluwatobi1/gh-api-data-fetch/internal/adapters/db/gorm"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/core/domain/models"
	"github.com/oluwatobi1/gh-api-data-fetch/internal/handlers"
	"gorm.io/driver/sqlite"
	gm "gorm.io/gorm"
)

type APPServer struct {
	addr string
}

func NewAPPServer(addr string) *APPServer {
	return &APPServer{
		addr: addr,
	}
}

func (s *APPServer) Run() error {

	db, err := gm.Open(sqlite.Open("github_monitoring.db"), &gm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	db.AutoMigrate(&models.Repository{}, &models.Commit{})

	repoRepo := gorm.NewRepository(db)
	commitRepo := gorm.NewCommitRepo(db)
	repoHandler := handlers.NewAppHandler(repoRepo, commitRepo)

	router := mux.NewRouter()

	router.Use(LoggingMiddleware)
	subRouter := router.PathPrefix("/api/v1").Subrouter()
	subRouter.HandleFunc("/fetch-repo", repoHandler.FetchRepository).Methods("GET")

	return http.ListenAndServe(s.addr, router)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
