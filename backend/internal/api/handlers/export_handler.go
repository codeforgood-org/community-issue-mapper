package handlers

import (
	"net/http"

	"github.com/codeforgood-org/community-issue-mapper/internal/export"
	"github.com/codeforgood-org/community-issue-mapper/internal/models"
	"github.com/codeforgood-org/community-issue-mapper/internal/service"
)

type ExportHandler struct {
	issueService  *service.IssueService
	exportService *export.ExportService
}

func NewExportHandler(issueService *service.IssueService) *ExportHandler {
	return &ExportHandler{
		issueService:  issueService,
		exportService: export.NewExportService(),
	}
}

func (h *ExportHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	filter := models.IssueFilter{
		Category: r.URL.Query().Get("category"),
		Status:   r.URL.Query().Get("status"),
		Limit:    10000, // Max for export
	}

	issues, err := h.issueService.ListIssues(filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch issues")
		return
	}

	csvData, err := h.exportService.ExportToCSV(issues)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate CSV")
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=issues.csv")
	w.Write(csvData)
}

func (h *ExportHandler) ExportPDF(w http.ResponseWriter, r *http.Request) {
	filter := models.IssueFilter{
		Category: r.URL.Query().Get("category"),
		Status:   r.URL.Query().Get("status"),
		Limit:    1000,
	}

	issues, err := h.issueService.ListIssues(filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch issues")
		return
	}

	pdfData, err := h.exportService.ExportToPDF(issues, "Community Issues Report")
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate PDF")
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=issues.pdf")
	w.Write(pdfData)
}

func (h *ExportHandler) ExportStatsPDF(w http.ResponseWriter, r *http.Request) {
	stats, err := h.issueService.GetStats()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch stats")
		return
	}

	pdfData, err := h.exportService.ExportStatsToPDF(stats)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate PDF")
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=stats.pdf")
	w.Write(pdfData)
}
