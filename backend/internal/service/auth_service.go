package service

import (
	"fmt"
	"regexp"

	"github.com/codeforgood-org/community-issue-mapper/internal/auth"
	"github.com/codeforgood-org/community-issue-mapper/internal/models"
	"github.com/codeforgood-org/community-issue-mapper/internal/repository"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type AuthService struct {
	userRepo   *repository.UserRepository
	jwtService *auth.JWTService
}

func NewAuthService(userRepo *repository.UserRepository, jwtService *auth.JWTService) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.LoginResponse, error) {
	// Validate input
	if err := s.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	// Check if user already exists
	existing, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("email already registered")
	}

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: passwordHash,
		Role:         models.RoleUser,
		Verified:     false,
		Active:       true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate token
	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check password
	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate token
	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) GetProfile(userID int64) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *AuthService) UpdateProfile(userID int64, req *models.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Bio != nil {
		user.Bio = *req.Bio
	}
	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) validateRegisterRequest(req *models.RegisterRequest) error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if !emailRegex.MatchString(req.Email) {
		return fmt.Errorf("invalid email format")
	}
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(req.Name) < 2 {
		return fmt.Errorf("name must be at least 2 characters")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(req.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	return nil
}

type CommentService struct {
	commentRepo  *repository.CommentRepository
	activityRepo *repository.ActivityRepository
}

func NewCommentService(commentRepo *repository.CommentRepository, activityRepo *repository.ActivityRepository) *CommentService {
	return &CommentService{
		commentRepo:  commentRepo,
		activityRepo: activityRepo,
	}
}

func (s *CommentService) CreateComment(issueID, userID int64, req *models.CreateCommentRequest) (*models.Comment, error) {
	if req.Content == "" {
		return nil, fmt.Errorf("comment content is required")
	}

	if len(req.Content) > 2000 {
		return nil, fmt.Errorf("comment too long (max 2000 characters)")
	}

	comment := &models.Comment{
		IssueID: issueID,
		UserID:  userID,
		Content: req.Content,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	// Log activity
	s.activityRepo.Create(&models.Activity{
		UserID:  userID,
		IssueID: issueID,
		Action:  "commented",
	})

	return comment, nil
}

func (s *CommentService) GetComments(issueID int64) ([]*models.Comment, error) {
	return s.commentRepo.GetByIssueID(issueID)
}

func (s *CommentService) DeleteComment(commentID, userID int64) error {
	return s.commentRepo.Delete(commentID, userID)
}
