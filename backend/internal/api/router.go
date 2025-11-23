package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/codeforgood-org/community-issue-mapper/internal/api/handlers"
	"github.com/codeforgood-org/community-issue-mapper/internal/api/middleware"
	"github.com/codeforgood-org/community-issue-mapper/internal/auth"
	"github.com/codeforgood-org/community-issue-mapper/internal/config"
	"github.com/codeforgood-org/community-issue-mapper/internal/models"
	"github.com/codeforgood-org/community-issue-mapper/internal/repository"
	"github.com/codeforgood-org/community-issue-mapper/internal/service"
	ws "github.com/codeforgood-org/community-issue-mapper/internal/websocket"
)

func NewRouter(db *sql.DB, cfg *config.Config) http.Handler {
	r := mux.NewRouter()

	// Initialize JWT service
	jwtService := auth.NewJWTService(cfg.Server.Env, 24*time.Hour) // Use proper secret and expiry from config

	// Initialize repositories
	issueRepo := repository.NewIssueRepository(db)
	userRepo := repository.NewUserRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	activityRepo := repository.NewActivityRepository(db)

	// Initialize services
	issueService := service.NewIssueService(issueRepo)
	authService := service.NewAuthService(userRepo, jwtService)
	commentService := service.NewCommentService(commentRepo, activityRepo)
	analyticsService := service.NewAnalyticsService(db, userRepo, activityRepo)

	// Initialize WebSocket hub
	hub := ws.NewHub()
	go hub.Run()

	// Initialize handlers
	issueHandler := handlers.NewIssueHandler(issueService, cfg)
	healthHandler := handlers.NewHealthHandler(db)
	authHandler := handlers.NewAuthHandler(authService)
	commentHandler := handlers.NewCommentHandler(commentService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	wsHandler := handlers.NewWebSocketHandler(hub)

	// Health check routes
	r.HandleFunc("/health", healthHandler.Health).Methods("GET")
	r.HandleFunc("/ready", healthHandler.Ready).Methods("GET")

	// WebSocket route
	r.HandleFunc("/ws", wsHandler.ServeWS)

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Public routes (no auth required)
	public := api.PathPrefix("").Subrouter()

	// Auth routes
	public.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	public.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// Issue routes (public read, auth for write)
	public.HandleFunc("/issues", issueHandler.ListIssues).Methods("GET")
	public.HandleFunc("/issues/{id}", issueHandler.GetIssue).Methods("GET")
	public.HandleFunc("/issues/{id}/comments", commentHandler.GetComments).Methods("GET")

	// Stats route (public)
	public.HandleFunc("/stats", issueHandler.GetStats).Methods("GET")
	public.HandleFunc("/analytics", analyticsHandler.GetAnalytics).Methods("GET")

	// Protected routes (auth required)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.Authentication(jwtService))

	// User routes
	protected.HandleFunc("/auth/profile", authHandler.GetProfile).Methods("GET")
	protected.HandleFunc("/auth/profile", authHandler.UpdateProfile).Methods("PATCH")

	// Issue routes (authenticated)
	protected.HandleFunc("/issues", issueHandler.CreateIssue).Methods("POST")
	protected.HandleFunc("/issues/{id}", issueHandler.UpdateIssue).Methods("PUT", "PATCH")
	protected.HandleFunc("/issues/{id}", issueHandler.DeleteIssue).Methods("DELETE")
	protected.HandleFunc("/issues/{id}/upload", issueHandler.UploadImage).Methods("POST")
	protected.HandleFunc("/issues/{id}/vote", issueHandler.VoteIssue).Methods("POST")

	// Comment routes (authenticated)
	protected.HandleFunc("/issues/{id}/comments", commentHandler.CreateComment).Methods("POST")

	// Admin routes
	admin := protected.PathPrefix("").Subrouter()
	admin.Use(middleware.RequireRole(models.RoleAdmin, models.RoleModerator))
	admin.HandleFunc("/admin/issues/{id}/status", issueHandler.UpdateIssue).Methods("PATCH")

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
