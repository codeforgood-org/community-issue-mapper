package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/codeforgood-org/community-issue-mapper/internal/models"
)

type IssueRepository struct {
	db *sql.DB
}

func NewIssueRepository(db *sql.DB) *IssueRepository {
	return &IssueRepository{db: db}
}

func (r *IssueRepository) Create(issue *models.CreateIssueRequest) (*models.Issue, error) {
	query := `
		INSERT INTO issues (title, description, category, latitude, longitude, address, reporter_name, reporter_email, status, votes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	result := &models.Issue{
		Title:         issue.Title,
		Description:   issue.Description,
		Category:      issue.Category,
		Latitude:      issue.Latitude,
		Longitude:     issue.Longitude,
		Address:       issue.Address,
		ReporterName:  issue.ReporterName,
		ReporterEmail: issue.ReporterEmail,
		Status:        models.StatusNew,
		Votes:         0,
	}

	err := r.db.QueryRow(
		query,
		issue.Title,
		issue.Description,
		issue.Category,
		issue.Latitude,
		issue.Longitude,
		issue.Address,
		issue.ReporterName,
		issue.ReporterEmail,
		models.StatusNew,
		0,
		now,
		now,
	).Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create issue: %w", err)
	}

	return result, nil
}

func (r *IssueRepository) GetByID(id int64) (*models.Issue, error) {
	query := `
		SELECT id, title, description, category, status, latitude, longitude,
		       address, image_url, reporter_name, reporter_email, votes, created_at, updated_at
		FROM issues
		WHERE id = $1
	`

	issue := &models.Issue{}
	err := r.db.QueryRow(query, id).Scan(
		&issue.ID,
		&issue.Title,
		&issue.Description,
		&issue.Category,
		&issue.Status,
		&issue.Latitude,
		&issue.Longitude,
		&issue.Address,
		&issue.ImageURL,
		&issue.ReporterName,
		&issue.ReporterEmail,
		&issue.Votes,
		&issue.CreatedAt,
		&issue.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}

	return issue, nil
}

func (r *IssueRepository) List(filter models.IssueFilter) ([]*models.Issue, error) {
	query := `
		SELECT id, title, description, category, status, latitude, longitude,
		       address, image_url, reporter_name, reporter_email, votes, created_at, updated_at
		FROM issues
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 1

	if filter.Category != "" {
		query += fmt.Sprintf(" AND category = $%d", argCount)
		args = append(args, filter.Category)
		argCount++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, filter.Status)
		argCount++
	}

	if filter.MinLat != 0 && filter.MaxLat != 0 && filter.MinLng != 0 && filter.MaxLng != 0 {
		query += fmt.Sprintf(" AND latitude BETWEEN $%d AND $%d", argCount, argCount+1)
		args = append(args, filter.MinLat, filter.MaxLat)
		argCount += 2

		query += fmt.Sprintf(" AND longitude BETWEEN $%d AND $%d", argCount, argCount+1)
		args = append(args, filter.MinLng, filter.MaxLng)
		argCount += 2
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filter.Limit)
		argCount++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, filter.Offset)
		argCount++
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list issues: %w", err)
	}
	defer rows.Close()

	issues := []*models.Issue{}
	for rows.Next() {
		issue := &models.Issue{}
		err := rows.Scan(
			&issue.ID,
			&issue.Title,
			&issue.Description,
			&issue.Category,
			&issue.Status,
			&issue.Latitude,
			&issue.Longitude,
			&issue.Address,
			&issue.ImageURL,
			&issue.ReporterName,
			&issue.ReporterEmail,
			&issue.Votes,
			&issue.CreatedAt,
			&issue.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan issue: %w", err)
		}
		issues = append(issues, issue)
	}

	return issues, nil
}

func (r *IssueRepository) Update(id int64, req *models.UpdateIssueRequest) (*models.Issue, error) {
	updates := []string{}
	args := []interface{}{}
	argCount := 1

	if req.Title != nil {
		updates = append(updates, fmt.Sprintf("title = $%d", argCount))
		args = append(args, *req.Title)
		argCount++
	}

	if req.Description != nil {
		updates = append(updates, fmt.Sprintf("description = $%d", argCount))
		args = append(args, *req.Description)
		argCount++
	}

	if req.Category != nil {
		updates = append(updates, fmt.Sprintf("category = $%d", argCount))
		args = append(args, *req.Category)
		argCount++
	}

	if req.Status != nil {
		updates = append(updates, fmt.Sprintf("status = $%d", argCount))
		args = append(args, *req.Status)
		argCount++
	}

	if req.Address != nil {
		updates = append(updates, fmt.Sprintf("address = $%d", argCount))
		args = append(args, *req.Address)
		argCount++
	}

	if len(updates) == 0 {
		return r.GetByID(id)
	}

	updates = append(updates, fmt.Sprintf("updated_at = $%d", argCount))
	args = append(args, time.Now())
	argCount++

	args = append(args, id)

	query := fmt.Sprintf(
		"UPDATE issues SET %s WHERE id = $%d",
		strings.Join(updates, ", "),
		argCount,
	)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update issue: %w", err)
	}

	return r.GetByID(id)
}

func (r *IssueRepository) Delete(id int64) error {
	query := "DELETE FROM issues WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete issue: %w", err)
	}
	return nil
}

func (r *IssueRepository) UpdateImageURL(id int64, imageURL string) error {
	query := "UPDATE issues SET image_url = $1, updated_at = $2 WHERE id = $3"
	_, err := r.db.Exec(query, imageURL, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update image URL: %w", err)
	}
	return nil
}

func (r *IssueRepository) IncrementVotes(id int64) error {
	query := "UPDATE issues SET votes = votes + 1, updated_at = $1 WHERE id = $2"
	_, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to increment votes: %w", err)
	}
	return nil
}

func (r *IssueRepository) GetStats() (*models.Stats, error) {
	stats := &models.Stats{
		IssuesByStatus:   make(map[string]int),
		IssuesByCategory: make(map[string]int),
	}

	// Get total count
	err := r.db.QueryRow("SELECT COUNT(*) FROM issues").Scan(&stats.TotalIssues)
	if err != nil {
		return nil, fmt.Errorf("failed to get total issues: %w", err)
	}

	// Get count by status
	rows, err := r.db.Query("SELECT status, COUNT(*) FROM issues GROUP BY status")
	if err != nil {
		return nil, fmt.Errorf("failed to get issues by status: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		stats.IssuesByStatus[status] = count
	}

	// Get count by category
	rows, err = r.db.Query("SELECT category, COUNT(*) FROM issues GROUP BY category")
	if err != nil {
		return nil, fmt.Errorf("failed to get issues by category: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		var count int
		if err := rows.Scan(&category, &count); err != nil {
			return nil, err
		}
		stats.IssuesByCategory[category] = count
	}

	return stats, nil
}
