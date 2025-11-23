package api

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/codeforgood-org/community-issue-mapper/internal/api/handlers"
	"github.com/codeforgood-org/community-issue-mapper/internal/api/middleware"
	"github.com/codeforgood-org/community-issue-mapper/internal/config"
	"github.com/codeforgood-org/community-issue-mapper/internal/repository"
	"github.com/codeforgood-org/community-issue-mapper/internal/service"
)

func NewRouter(db *sql.DB, cfg *config.Config) http.Handler {
	r := mux.NewRouter()

	// Initialize repositories
	issueRepo := repository.NewIssueRepository(db)

	// Initialize services
	issueService := service.NewIssueService(issueRepo)

	// Initialize handlers
	issueHandler := handlers.NewIssueHandler(issueService, cfg)
	healthHandler := handlers.NewHealthHandler(db)

	// Health check routes
	r.HandleFunc("/health", healthHandler.Health).Methods("GET")
	r.HandleFunc("/ready", healthHandler.Ready).Methods("GET")

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Issue routes
	api.HandleFunc("/issues", issueHandler.ListIssues).Methods("GET")
	api.HandleFunc("/issues", issueHandler.CreateIssue).Methods("POST")
	api.HandleFunc("/issues/{id}", issueHandler.GetIssue).Methods("GET")
	api.HandleFunc("/issues/{id}", issueHandler.UpdateIssue).Methods("PUT", "PATCH")
	api.HandleFunc("/issues/{id}", issueHandler.DeleteIssue).Methods("DELETE")
	api.HandleFunc("/issues/{id}/upload", issueHandler.UploadImage).Methods("POST")
	api.HandleFunc("/issues/{id}/vote", issueHandler.VoteIssue).Methods("POST")

	// Stats route
	api.HandleFunc("/stats", issueHandler.GetStats).Methods("GET")

	// Serve uploaded files
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir(cfg.Upload.UploadDir))))

	// Serve static frontend files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend")))

	// Apply middleware
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimit.Requests, cfg.RateLimit.Duration)
	handler := middleware.Logging(r)
	handler = rateLimiter.Limit(handler)
	handler = middleware.CORS(cfg)(handler)

	return handler
}
