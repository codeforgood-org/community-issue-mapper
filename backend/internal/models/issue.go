package models

import (
	"time"
)

type IssueStatus string
type IssueCategory string

const (
	StatusNew        IssueStatus = "new"
	StatusInProgress IssueStatus = "in_progress"
	StatusResolved   IssueStatus = "resolved"
	StatusClosed     IssueStatus = "closed"
)

const (
	CategoryPothole       IssueCategory = "pothole"
	CategoryAccessibility IssueCategory = "accessibility"
	CategoryStreetlight   IssueCategory = "streetlight"
	CategoryGraffiti      IssueCategory = "graffiti"
	CategoryTrash         IssueCategory = "trash"
	CategoryOther         IssueCategory = "other"
)

type Issue struct {
	ID          int64         `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Category    IssueCategory `json:"category"`
	Status      IssueStatus   `json:"status"`
	Latitude    float64       `json:"latitude"`
	Longitude   float64       `json:"longitude"`
	Address     string        `json:"address,omitempty"`
	ImageURL    string        `json:"image_url,omitempty"`
	ReporterName string       `json:"reporter_name,omitempty"`
	ReporterEmail string      `json:"reporter_email,omitempty"`
	Votes       int           `json:"votes"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type CreateIssueRequest struct {
	Title         string        `json:"title"`
	Description   string        `json:"description"`
	Category      IssueCategory `json:"category"`
	Latitude      float64       `json:"latitude"`
	Longitude     float64       `json:"longitude"`
	Address       string        `json:"address,omitempty"`
	ReporterName  string        `json:"reporter_name,omitempty"`
	ReporterEmail string        `json:"reporter_email,omitempty"`
}

type UpdateIssueRequest struct {
	Title       *string        `json:"title,omitempty"`
	Description *string        `json:"description,omitempty"`
	Category    *IssueCategory `json:"category,omitempty"`
	Status      *IssueStatus   `json:"status,omitempty"`
	Address     *string        `json:"address,omitempty"`
}

type IssueFilter struct {
	Category  string
	Status    string
	MinLat    float64
	MaxLat    float64
	MinLng    float64
	MaxLng    float64
	Limit     int
	Offset    int
}

type Stats struct {
	TotalIssues      int            `json:"total_issues"`
	IssuesByStatus   map[string]int `json:"issues_by_status"`
	IssuesByCategory map[string]int `json:"issues_by_category"`
}
