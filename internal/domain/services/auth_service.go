package services

import (
	"context"

	"github.com/wonyus/backend-challenge/internal/domain/entities"
)

type AuthService interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
	GenerateToken(ctx context.Context, user *entities.User) (string, error)
	ValidateToken(ctx context.Context, token string) (*entities.User, error)
}
