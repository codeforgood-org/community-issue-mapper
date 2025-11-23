package models

import (
	"time"
)

type UserRole string

const (
	RoleUser      UserRole = "user"
	RoleModerator UserRole = "moderator"
	RoleAdmin     UserRole = "admin"
)

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	Avatar       string    `json:"avatar,omitempty"`
	Bio          string    `json:"bio,omitempty"`
	Verified     bool      `json:"verified"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type UpdateUserRequest struct {
	Name   *string `json:"name,omitempty"`
	Bio    *string `json:"bio,omitempty"`
	Avatar *string `json:"avatar,omitempty"`
}

type Comment struct {
	ID        int64     `json:"id"`
	IssueID   int64     `json:"issue_id"`
	UserID    int64     `json:"user_id"`
	User      *User     `json:"user,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCommentRequest struct {
	Content string `json:"content"`
}

type Activity struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	IssueID   int64     `json:"issue_id"`
	Action    string    `json:"action"` // created, updated, commented, voted, resolved
	Metadata  string    `json:"metadata,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Notification struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"` // issue_update, comment, mention
	IssueID   *int64    `json:"issue_id,omitempty"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

type AnalyticsData struct {
	TotalIssues       int                       `json:"total_issues"`
	TotalUsers        int                       `json:"total_users"`
	TotalComments     int                       `json:"total_comments"`
	IssuesThisWeek    int                       `json:"issues_this_week"`
	IssuesThisMonth   int                       `json:"issues_this_month"`
	ResolvedThisWeek  int                       `json:"resolved_this_week"`
	ResolvedThisMonth int                       `json:"resolved_this_month"`
	AverageResolutionTime float64              `json:"avg_resolution_time_hours"`
	IssuesByCategory  map[string]int            `json:"issues_by_category"`
	IssuesByStatus    map[string]int            `json:"issues_by_status"`
	IssuesTrend       []TrendData               `json:"issues_trend"`
	TopContributors   []ContributorData         `json:"top_contributors"`
	RecentActivity    []Activity                `json:"recent_activity"`
}

type TrendData struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type ContributorData struct {
	User         *User `json:"user"`
	IssueCount   int   `json:"issue_count"`
	CommentCount int   `json:"comment_count"`
	VoteCount    int   `json:"vote_count"`
}
