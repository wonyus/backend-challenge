package ports

import (
	"context"

	"github.com/wonyus/backend-challenge/internal/application/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*dto.UserResponse, error)
	GetAllUsers(ctx context.Context) (*dto.UsersListResponse, error)
	UpdateUser(ctx context.Context, id primitive.ObjectID, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
	GetUserCount(ctx context.Context) (int64, error)
}
