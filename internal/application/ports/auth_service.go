package ports

import (
	"context"

	"github.com/wonyus/backend-challenge/internal/application/dto"
)

type AuthService interface {
	Register(ctx context.Context, req *dto.CreateUserRequest) (*dto.RegisterResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	ValidateToken(ctx context.Context, token string) (*dto.UserResponse, error)
}
