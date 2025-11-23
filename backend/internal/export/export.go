package export

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

	"github.com/codeforgood-org/community-issue-mapper/internal/models"
	"github.com/jung-kurt/gofpdf"
)

type ExportService struct{}

func NewExportService() *ExportService {
	return &ExportService{}
}

// ExportToCSV exports issues to CSV format
func (s *ExportService) ExportToCSV(issues []*models.Issue) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{
		"ID", "Title", "Description", "Category", "Status",
		"Latitude", "Longitude", "Address", "Reporter Name",
		"Reporter Email", "Votes", "Created At", "Updated At",
	}
	if err := writer.Write(header); err != nil {
		return nil, err
	}

	// Write data
	for _, issue := range issues {
		record := []string{
			strconv.FormatInt(issue.ID, 10),
			issue.Title,
			issue.Description,
			string(issue.Category),
			string(issue.Status),
			strconv.FormatFloat(issue.Latitude, 'f', 6, 64),
			strconv.FormatFloat(issue.Longitude, 'f', 6, 64),
			issue.Address,
			issue.ReporterName,
			issue.ReporterEmail,
			strconv.Itoa(issue.Votes),
			issue.CreatedAt.Format(time.RFC3339),
			issue.UpdatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ExportToPDF exports issues to PDF format
func (s *ExportService) ExportToPDF(issues []*models.Issue, title string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, title)
	pdf.Ln(15)

	// Subtitle
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(190, 5, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(5)
	pdf.Cell(190, 5, fmt.Sprintf("Total Issues: %d", len(issues)))
	pdf.Ln(15)

	// Table header
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(37, 99, 235)
	pdf.SetTextColor(255, 255, 255)

	colWidths := []float64{15, 60, 30, 25, 30, 30}
	headers := []string{"ID", "Title", "Category", "Status", "Reporter", "Created"}

	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 8, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Table data
	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(0, 0, 0)
	fill := false

	for _, issue := range issues {
		if fill {
			pdf.SetFillColor(240, 240, 240)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		// Truncate long titles
		title := issue.Title
		if len(title) > 40 {
			title = title[:37] + "..."
		}

		data := []string{
			strconv.FormatInt(issue.ID, 10),
			title,
			string(issue.Category),
			string(issue.Status),
			issue.ReporterName,
			issue.CreatedAt.Format("2006-01-02"),
		}

		for i, text := range data {
			pdf.CellFormat(colWidths[i], 7, text, "1", 0, "L", true, 0, "")
		}
		pdf.Ln(-1)
		fill = !fill
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ExportStatsToPDF exports statistics to PDF
func (s *ExportService) ExportStatsToP DF(stats *models.Stats) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(190, 10, "Community Issue Mapper - Statistics Report")
	pdf.Ln(15)

	// Date
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(190, 5, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(15)

	// Total Issues
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 8, "Overview")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(100, 7, fmt.Sprintf("Total Issues: %d", stats.TotalIssues))
	pdf.Ln(10)

	// Issues by Status
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 8, "Issues by Status")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	for status, count := range stats.IssuesByStatus {
		pdf.Cell(100, 6, fmt.Sprintf("%s: %d", status, count))
		pdf.Ln(7)
	}
	pdf.Ln(5)

	// Issues by Category
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 8, "Issues by Category")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	for category, count := range stats.IssuesByCategory {
		pdf.Cell(100, 6, fmt.Sprintf("%s: %d", category, count))
		pdf.Ln(7)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
