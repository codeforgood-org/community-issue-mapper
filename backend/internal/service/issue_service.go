package service

import (
	"fmt"

	"github.com/codeforgood-org/community-issue-mapper/internal/models"
	"github.com/codeforgood-org/community-issue-mapper/internal/repository"
)

type IssueService struct {
	repo *repository.IssueRepository
}

func NewIssueService(repo *repository.IssueRepository) *IssueService {
	return &IssueService{repo: repo}
}

func (s *IssueService) CreateIssue(req *models.CreateIssueRequest) (*models.Issue, error) {
	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	return s.repo.Create(req)
}

func (s *IssueService) GetIssue(id int64) (*models.Issue, error) {
	issue, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if issue == nil {
		return nil, fmt.Errorf("issue not found")
	}
	return issue, nil
}

func (s *IssueService) ListIssues(filter models.IssueFilter) ([]*models.Issue, error) {
	return s.repo.List(filter)
}

func (s *IssueService) UpdateIssue(id int64, req *models.UpdateIssueRequest) (*models.Issue, error) {
	// Check if issue exists
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("issue not found")
	}

	return s.repo.Update(id, req)
}

func (s *IssueService) DeleteIssue(id int64) error {
	// Check if issue exists
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("issue not found")
	}

	return s.repo.Delete(id)
}

func (s *IssueService) UpdateImageURL(id int64, imageURL string) error {
	return s.repo.UpdateImageURL(id, imageURL)
}

func (s *IssueService) VoteIssue(id int64) error {
	// Check if issue exists
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("issue not found")
	}

	return s.repo.IncrementVotes(id)
}

func (s *IssueService) GetStats() (*models.Stats, error) {
	return s.repo.GetStats()
}

func (s *IssueService) validateCreateRequest(req *models.CreateIssueRequest) error {
	if req.Title == "" {
		return fmt.Errorf("title is required")
	}
	if req.Description == "" {
		return fmt.Errorf("description is required")
	}
	if req.Category == "" {
		return fmt.Errorf("category is required")
	}
	if req.Latitude < -90 || req.Latitude > 90 {
		return fmt.Errorf("invalid latitude")
	}
	if req.Longitude < -180 || req.Longitude > 180 {
		return fmt.Errorf("invalid longitude")
	}

	// Validate category
	validCategories := map[models.IssueCategory]bool{
		models.CategoryPothole:       true,
		models.CategoryAccessibility: true,
		models.CategoryStreetlight:   true,
		models.CategoryGraffiti:      true,
		models.CategoryTrash:         true,
		models.CategoryOther:         true,
	}

	if !validCategories[req.Category] {
		return fmt.Errorf("invalid category")
	}

	return nil
}
