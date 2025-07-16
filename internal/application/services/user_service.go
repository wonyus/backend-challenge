package services

import (
	"context"

	"github.com/wonyus/backend-challenge/internal/application/dto"
	"github.com/wonyus/backend-challenge/internal/application/ports"
	"github.com/wonyus/backend-challenge/internal/domain/entities"
	domainErrors "github.com/wonyus/backend-challenge/internal/domain/errors"
	"github.com/wonyus/backend-challenge/internal/domain/repositories"
	domainServices "github.com/wonyus/backend-challenge/internal/domain/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userService struct {
	userRepo    repositories.UserRepository
	authService domainServices.AuthService
}

func NewUserService(userRepo repositories.UserRepository, authService domainServices.AuthService) ports.UserService {
	return &userService{
		userRepo:    userRepo,
		authService: authService,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
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

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *userService) GetUserByID(ctx context.Context, id primitive.ObjectID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
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
func (s *userService) GetAllUsers(ctx context.Context) (*dto.UsersListResponse, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		}
	}
	return &dto.UsersListResponse{
		Users: userResponses,
		Total: len(userResponses),
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, id primitive.ObjectID, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		user.UpdateName(req.Name)
	}

	if req.Email != "" {
		// Check if email is already taken by another user
		existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, domainErrors.ErrUserAlreadyExists
		}

		user.UpdateEmail(req.Email)
	}

	// Save updated user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *userService) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.userRepo.Delete(ctx, user.ID)
}

func (s *userService) GetUserCount(ctx context.Context) (int64, error) {
	return s.userRepo.Count(ctx)
}
