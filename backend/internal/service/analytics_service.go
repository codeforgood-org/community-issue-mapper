package service

import (
	"database/sql"
	"time"

	"github.com/codeforgood-org/community-issue-mapper/internal/models"
	"github.com/codeforgood-org/community-issue-mapper/internal/repository"
)

type AnalyticsService struct {
	db           *sql.DB
	userRepo     *repository.UserRepository
	activityRepo *repository.ActivityRepository
}

func NewAnalyticsService(db *sql.DB, userRepo *repository.UserRepository, activityRepo *repository.ActivityRepository) *AnalyticsService {
	return &AnalyticsService{
		db:           db,
		userRepo:     userRepo,
		activityRepo: activityRepo,
	}
}

func (s *AnalyticsService) GetAnalytics() (*models.AnalyticsData, error) {
	analytics := &models.AnalyticsData{
		IssuesByCategory: make(map[string]int),
		IssuesByStatus:   make(map[string]int),
	}

	// Total issues
	err := s.db.QueryRow("SELECT COUNT(*) FROM issues").Scan(&analytics.TotalIssues)
	if err != nil {
		return nil, err
	}

	// Total users
	err = s.db.QueryRow("SELECT COUNT(*) FROM users WHERE active = true").Scan(&analytics.TotalUsers)
	if err != nil {
		return nil, err
	}

	// Total comments
	err = s.db.QueryRow("SELECT COUNT(*) FROM comments").Scan(&analytics.TotalComments)
	if err != nil {
		return nil, err
	}

	// Issues this week
	weekAgo := time.Now().AddDate(0, 0, -7)
	err = s.db.QueryRow(
		"SELECT COUNT(*) FROM issues WHERE created_at >= $1",
		weekAgo,
	).Scan(&analytics.IssuesThisWeek)
	if err != nil {
		return nil, err
	}

	// Issues this month
	monthAgo := time.Now().AddDate(0, -1, 0)
	err = s.db.QueryRow(
		"SELECT COUNT(*) FROM issues WHERE created_at >= $1",
		monthAgo,
	).Scan(&analytics.IssuesThisMonth)
	if err != nil {
		return nil, err
	}

	// Resolved this week
	err = s.db.QueryRow(
		"SELECT COUNT(*) FROM issues WHERE status IN ('resolved', 'closed') AND updated_at >= $1",
		weekAgo,
	).Scan(&analytics.ResolvedThisWeek)
	if err != nil {
		return nil, err
	}

	// Resolved this month
	err = s.db.QueryRow(
		"SELECT COUNT(*) FROM issues WHERE status IN ('resolved', 'closed') AND updated_at >= $1",
		monthAgo,
	).Scan(&analytics.ResolvedThisMonth)
	if err != nil {
		return nil, err
	}

	// Average resolution time
	err = s.db.QueryRow(`
		SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (updated_at - created_at))/3600), 0)
		FROM issues
		WHERE status IN ('resolved', 'closed')
	`).Scan(&analytics.AverageResolutionTime)
	if err != nil {
		return nil, err
	}

	// Issues by category
	rows, err := s.db.Query("SELECT category, COUNT(*) FROM issues GROUP BY category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		var count int
		if err := rows.Scan(&category, &count); err != nil {
			return nil, err
		}
		analytics.IssuesByCategory[category] = count
	}

	// Issues by status
	rows, err = s.db.Query("SELECT status, COUNT(*) FROM issues GROUP BY status")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		analytics.IssuesByStatus[status] = count
	}

	// Issues trend (last 30 days)
	rows, err = s.db.Query(`
		SELECT DATE(created_at) as date, COUNT(*) as count
		FROM issues
		WHERE created_at >= NOW() - INTERVAL '30 days'
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	analytics.IssuesTrend = []models.TrendData{}
	for rows.Next() {
		var trend models.TrendData
		var date time.Time
		if err := rows.Scan(&date, &trend.Count); err != nil {
			return nil, err
		}
		trend.Date = date.Format("2006-01-02")
		analytics.IssuesTrend = append(analytics.IssuesTrend, trend)
	}

	// Top contributors
	contributors, err := s.userRepo.GetTopContributors(10)
	if err != nil {
		return nil, err
	}
	analytics.TopContributors = contributors

	// Recent activity
	activities, err := s.activityRepo.GetRecent(20)
	if err != nil {
		return nil, err
	}
	analytics.RecentActivity = activities

	return analytics, nil
}
