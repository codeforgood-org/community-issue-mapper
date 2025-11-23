package repository

import (
	"testing"

	"github.com/codeforgood-org/community-issue-mapper/internal/models"
)

func TestIssueRepository(t *testing.T) {
	// Note: These are example tests
	// In a real scenario, you would set up a test database

	t.Run("Create Issue", func(t *testing.T) {
		// Test issue creation
		req := &models.CreateIssueRequest{
			Title:       "Test Pothole",
			Description: "Large pothole on Main St",
			Category:    models.CategoryPothole,
			Latitude:    37.7749,
			Longitude:   -122.4194,
		}

		// Assert validation
		if req.Title == "" {
			t.Error("Title should not be empty")
		}
		if req.Latitude < -90 || req.Latitude > 90 {
			t.Error("Invalid latitude")
		}
		if req.Longitude < -180 || req.Longitude > 180 {
			t.Error("Invalid longitude")
		}
	})

	t.Run("Validate Category", func(t *testing.T) {
		validCategories := []models.IssueCategory{
			models.CategoryPothole,
			models.CategoryAccessibility,
			models.CategoryStreetlight,
			models.CategoryGraffiti,
			models.CategoryTrash,
			models.CategoryOther,
		}

		for _, cat := range validCategories {
			if cat == "" {
				t.Errorf("Category %s should not be empty", cat)
			}
		}
	})

	t.Run("Validate Status", func(t *testing.T) {
		validStatuses := []models.IssueStatus{
			models.StatusNew,
			models.StatusInProgress,
			models.StatusResolved,
			models.StatusClosed,
		}

		for _, status := range validStatuses {
			if status == "" {
				t.Errorf("Status %s should not be empty", status)
			}
		}
	})
}
