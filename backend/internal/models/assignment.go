package models

import "time"

type Assignment struct {
	ID          int64     `json:"id"`
	IssueID     int64     `json:"issue_id"`
	AssignedTo  int64     `json:"assigned_to"`
	AssignedBy  int64     `json:"assigned_by"`
	Department  string    `json:"department,omitempty"`
	Priority    string    `json:"priority"` // low, medium, high, critical
	DueDate     time.Time `json:"due_date,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	Status      string    `json:"status"` // assigned, in_progress, completed
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type IssueTemplate struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Category    IssueCategory `json:"category"`
	Description string    `json:"description"`
	Fields      string    `json:"fields"` // JSON string
	Active      bool      `json:"active"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BulkOperation struct {
	IssueIDs  []int64     `json:"issue_ids"`
	Operation string      `json:"operation"` // update_status, assign, delete
	Data      interface{} `json:"data"`
}

type BulkUpdateStatus struct {
	Status IssueStatus `json:"status"`
}

type BulkAssign struct {
	AssignedTo int64  `json:"assigned_to"`
	Department string `json:"department"`
	Priority   string `json:"priority"`
}

type SLAConfig struct {
	Category        IssueCategory `json:"category"`
	Priority        string        `json:"priority"`
	ResponseTime    int           `json:"response_time_hours"`    // Hours to first response
	ResolutionTime  int           `json:"resolution_time_hours"`  // Hours to resolution
}

type SLAStatus struct {
	IssueID           int64     `json:"issue_id"`
	ResponseDeadline  time.Time `json:"response_deadline"`
	ResolutionDeadline time.Time `json:"resolution_deadline"`
	ResponseMet       bool      `json:"response_met"`
	ResolutionMet     bool      `json:"resolution_met"`
	ResponseTime      int       `json:"response_time_hours"`
	ResolutionTime    int       `json:"resolution_time_hours,omitempty"`
}
