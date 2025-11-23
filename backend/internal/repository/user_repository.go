package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/codeforgood-org/community-issue-mapper/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (email, name, password_hash, role, verified, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	err := r.db.QueryRow(
		query,
		user.Email,
		user.Name,
		user.PasswordHash,
		user.Role,
		user.Verified,
		user.Active,
		now,
		now,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	query := `
		SELECT id, email, name, password_hash, role, avatar, bio, verified, active, created_at, updated_at
		FROM users
		WHERE id = $1 AND active = true
	`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.Role,
		&user.Avatar,
		&user.Bio,
		&user.Verified,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, name, password_hash, role, avatar, bio, verified, active, created_at, updated_at
		FROM users
		WHERE email = $1 AND active = true
	`

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.Role,
		&user.Avatar,
		&user.Bio,
		&user.Verified,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET name = $1, avatar = $2, bio = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.Exec(query, user.Name, user.Avatar, user.Bio, time.Now(), user.ID)
	return err
}

func (r *UserRepository) GetTopContributors(limit int) ([]models.ContributorData, error) {
	query := `
		SELECT
			u.id, u.email, u.name, u.avatar,
			COUNT(DISTINCT i.id) as issue_count,
			COUNT(DISTINCT c.id) as comment_count,
			COUNT(DISTINCT v.id) as vote_count
		FROM users u
		LEFT JOIN issues i ON i.user_id = u.id
		LEFT JOIN comments c ON c.user_id = u.id
		LEFT JOIN issue_votes v ON v.user_id = u.id
		WHERE u.active = true
		GROUP BY u.id
		ORDER BY (COUNT(DISTINCT i.id) + COUNT(DISTINCT c.id) + COUNT(DISTINCT v.id)) DESC
		LIMIT $1
	`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contributors := []models.ContributorData{}
	for rows.Next() {
		user := &models.User{}
		var issueCount, commentCount, voteCount int
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.Avatar,
			&issueCount,
			&commentCount,
			&voteCount,
		)
		if err != nil {
			return nil, err
		}

		contributors = append(contributors, models.ContributorData{
			User:         user,
			IssueCount:   issueCount,
			CommentCount: commentCount,
			VoteCount:    voteCount,
		})
	}

	return contributors, nil
}

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *models.Comment) error {
	query := `
		INSERT INTO comments (issue_id, user_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	return r.db.QueryRow(query, comment.IssueID, comment.UserID, comment.Content, now, now).
		Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
}

func (r *CommentRepository) GetByIssueID(issueID int64) ([]*models.Comment, error) {
	query := `
		SELECT c.id, c.issue_id, c.user_id, c.content, c.created_at, c.updated_at,
		       u.id, u.email, u.name, u.avatar, u.role
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.issue_id = $1
		ORDER BY c.created_at ASC
	`

	rows, err := r.db.Query(query, issueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*models.Comment{}
	for rows.Next() {
		comment := &models.Comment{User: &models.User{}}
		err := rows.Scan(
			&comment.ID,
			&comment.IssueID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.User.ID,
			&comment.User.Email,
			&comment.User.Name,
			&comment.User.Avatar,
			&comment.User.Role,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *CommentRepository) Delete(id, userID int64) error {
	query := "DELETE FROM comments WHERE id = $1 AND user_id = $2"
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("comment not found or unauthorized")
	}

	return nil
}

type ActivityRepository struct {
	db *sql.DB
}

func NewActivityRepository(db *sql.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func (r *ActivityRepository) Create(activity *models.Activity) error {
	query := `
		INSERT INTO activities (user_id, issue_id, action, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	return r.db.QueryRow(
		query,
		activity.UserID,
		activity.IssueID,
		activity.Action,
		activity.Metadata,
		time.Now(),
	).Scan(&activity.ID, &activity.CreatedAt)
}

func (r *ActivityRepository) GetRecent(limit int) ([]models.Activity, error) {
	query := `
		SELECT id, user_id, issue_id, action, metadata, created_at
		FROM activities
		ORDER BY created_at DESC
		LIMIT $1
	`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activities := []models.Activity{}
	for rows.Next() {
		activity := models.Activity{}
		err := rows.Scan(
			&activity.ID,
			&activity.UserID,
			&activity.IssueID,
			&activity.Action,
			&activity.Metadata,
			&activity.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	return activities, nil
}
