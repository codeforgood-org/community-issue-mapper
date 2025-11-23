package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/codeforgood-org/community-issue-mapper/internal/config"
	"github.com/codeforgood-org/community-issue-mapper/internal/models"
	"github.com/codeforgood-org/community-issue-mapper/internal/service"
)

type IssueHandler struct {
	service *service.IssueService
	config  *config.Config
}

func NewIssueHandler(service *service.IssueService, cfg *config.Config) *IssueHandler {
	return &IssueHandler{
		service: service,
		config:  cfg,
	}
}

func (h *IssueHandler) CreateIssue(w http.ResponseWriter, r *http.Request) {
	var req models.CreateIssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	issue, err := h.service.CreateIssue(&req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, issue)
}

func (h *IssueHandler) GetIssue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid issue ID")
		return
	}

	issue, err := h.service.GetIssue(id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, issue)
}

func (h *IssueHandler) ListIssues(w http.ResponseWriter, r *http.Request) {
	filter := models.IssueFilter{
		Category: r.URL.Query().Get("category"),
		Status:   r.URL.Query().Get("status"),
		Limit:    100, // default limit
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 1000 {
			filter.Limit = l
		}
	}

	if offset := r.URL.Query().Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
			filter.Offset = o
		}
	}

	// Bounding box filter
	if minLat := r.URL.Query().Get("min_lat"); minLat != "" {
		if lat, err := strconv.ParseFloat(minLat, 64); err == nil {
			filter.MinLat = lat
		}
	}
	if maxLat := r.URL.Query().Get("max_lat"); maxLat != "" {
		if lat, err := strconv.ParseFloat(maxLat, 64); err == nil {
			filter.MaxLat = lat
		}
	}
	if minLng := r.URL.Query().Get("min_lng"); minLng != "" {
		if lng, err := strconv.ParseFloat(minLng, 64); err == nil {
			filter.MinLng = lng
		}
	}
	if maxLng := r.URL.Query().Get("max_lng"); maxLng != "" {
		if lng, err := strconv.ParseFloat(maxLng, 64); err == nil {
			filter.MaxLng = lng
		}
	}

	issues, err := h.service.ListIssues(filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list issues")
		return
	}

	respondJSON(w, http.StatusOK, issues)
}

func (h *IssueHandler) UpdateIssue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid issue ID")
		return
	}

	var req models.UpdateIssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	issue, err := h.service.UpdateIssue(id, &req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, issue)
}

func (h *IssueHandler) DeleteIssue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid issue ID")
		return
	}

	if err := h.service.DeleteIssue(id); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *IssueHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid issue ID")
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(h.config.Upload.MaxSize); err != nil {
		respondError(w, http.StatusBadRequest, "File too large")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid image file")
		return
	}
	defer file.Close()

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		respondError(w, http.StatusBadRequest, "File must be an image")
		return
	}

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll(h.config.Upload.UploadDir, 0755); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create upload directory")
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", id, ext)
	filePath := filepath.Join(h.config.Upload.UploadDir, filename)

	// Create file
	dst, err := os.Create(filePath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to save file")
		return
	}
	defer dst.Close()

	// Copy file
	if _, err := io.Copy(dst, file); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to save file")
		return
	}

	// Update issue with image URL
	imageURL := fmt.Sprintf("/uploads/%s", filename)
	if err := h.service.UpdateImageURL(id, imageURL); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update issue")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"image_url": imageURL})
}

func (h *IssueHandler) VoteIssue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid issue ID")
		return
	}

	if err := h.service.VoteIssue(id); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	issue, err := h.service.GetIssue(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get updated issue")
		return
	}

	respondJSON(w, http.StatusOK, issue)
}

func (h *IssueHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetStats()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get stats")
		return
	}

	respondJSON(w, http.StatusOK, stats)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
