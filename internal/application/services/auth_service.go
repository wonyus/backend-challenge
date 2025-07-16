package services

import (
	"context"
	"fmt"

	"github.com/wonyus/backend-challenge/internal/application/dto"
	"github.com/wonyus/backend-challenge/internal/application/ports"
	"github.com/wonyus/backend-challenge/internal/domain/entities"
	domainErrors "github.com/wonyus/backend-challenge/internal/domain/errors"
	"github.com/wonyus/backend-challenge/internal/domain/repositories"
	domainServices "github.com/wonyus/backend-challenge/internal/domain/services"
)

type authService struct {
	userRepo    repositories.UserRepository
	authService domainServices.AuthService
}

func NewAuthService(userRepo repositories.UserRepository, service domainServices.AuthService) ports.AuthService {
	return &authService{
		userRepo:    userRepo,
		authService: service,
	}
}
func (s *authService) Register(ctx context.Context, req *dto.CreateUserRequest) (*dto.RegisterResponse, error) {
	fmt.Println(req)
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, domainErrors.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := s.authService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user entity
	user, err := entities.NewUser(req.Name, req.Email, hashedPassword)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &dto.RegisterResponse{
		ID:      user.ID,
		Message: "User registered successfully",
	}, nil
}

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, domainErrors.ErrInvalidCredentials
	}

	// Compare password
	if err := s.authService.ComparePassword(user.Password, req.Password); err != nil {
		return nil, domainErrors.ErrInvalidCredentials
	}

	// Generate token
	token, err := s.authService.GenerateToken(ctx, user)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*dto.UserResponse, error) {
	user, err := s.authService.ValidateToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}
