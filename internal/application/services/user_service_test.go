package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wonyus/backend-challenge/internal/application/dto"
	"github.com/wonyus/backend-challenge/internal/domain/entities"
	"github.com/wonyus/backend-challenge/internal/infrastructure/auth"
	mock_repositories "github.com/wonyus/backend-challenge/mock/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestService_User_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repositories.NewMockUserRepository(ctrl)
	jwtService := auth.NewJWTService("cfg.JWTSecret", mockUserRepo)
	userService := NewUserService(mockUserRepo, jwtService)

	var (
		ctx = context.Background()
		id  = primitive.NewObjectID()
		now = time.Now()
	)

	mockRequest := &dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password",
	}

	mockResponse := &dto.UserResponse{
		ID:        id,
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: now,
	}

	mockUserEntity := &entities.User{
		ID:        id,
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "hashed_password",
		CreatedAt: now,
		UpdatedAt: now,
	}

	t.Run("Create user success", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmail(gomock.Any(), mockRequest.Email).Return(nil, nil).Times(1)
		mockUserRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		response, err := userService.CreateUser(ctx, mockRequest)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, mockResponse.Email, response.Email)
	})

	t.Run("Create user failed", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmail(gomock.Any(), mockRequest.Email).Return(nil, nil).Times(1)
		mockUserRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("failed to create user")).Times(1)
		response, err := userService.CreateUser(ctx, mockRequest)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "failed to create user", err.Error())
	})

	t.Run("User already exists", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmail(gomock.Any(), mockRequest.Email).Return(mockUserEntity, nil).Times(1)
		response, err := userService.CreateUser(ctx, mockRequest)
		assert.Nil(t, response)
		assert.Equal(t, "user already exists", err.Error())
	})

	t.Run("User entity creation error", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmail(gomock.Any(), mockRequest.Email).Return(nil, nil).Times(1)
		mockRequest.Name = ""
		response, err := userService.CreateUser(ctx, mockRequest)
		assert.Nil(t, response)
		assert.Equal(t, "name, email, and password are required", err.Error())
	})
	t.Run("password hashing error", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmail(gomock.Any(), mockRequest.Email).Return(nil, nil).Times(1)
		mockRequest.Password = ""
		response, err := userService.CreateUser(ctx, mockRequest)
		assert.Nil(t, response)
		assert.Equal(t, "password cannot be empty", err.Error())
	})
}

func TestService_User_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repositories.NewMockUserRepository(ctrl)
	jwtService := auth.NewJWTService("cfg.JWTSecret", mockUserRepo)
	userService := NewUserService(mockUserRepo, jwtService)

	var (
		ctx = context.Background()
		id  = primitive.NewObjectID()
		now = time.Now()
	)

	mockUserEntity := &entities.User{
		ID:        id,
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "hashed_password",
		CreatedAt: now,
		UpdatedAt: now,
	}

	t.Run("Get user by ID", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(gomock.Any(), id).Return(mockUserEntity, nil).Times(1)
		response, err := userService.GetUserByID(ctx, id)
		assert.NoError(t, err)
		assert.NotNil(t, response)
	})

	t.Run("Get user by ID failed", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, errors.New("failed to get user")).Times(1)
		response, err := userService.GetUserByID(ctx, id)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "failed to get user", err.Error())
	})
}

func TestService_User_GetAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repositories.NewMockUserRepository(ctrl)
	jwtService := auth.NewJWTService("cfg.JWTSecret", mockUserRepo)
	userService := NewUserService(mockUserRepo, jwtService)

	var (
		ctx = context.Background()
		now = time.Now()
	)

	mockUsers := []*entities.User{
		{
			ID:        primitive.NewObjectID(),
			Name:      "Test User 1",
			Email:     "test1@example.com",
			Password:  "hashed_password",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        primitive.NewObjectID(),
			Name:      "Test User 2",
			Email:     "test2@example.com",
			Password:  "hashed_password",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	t.Run("Get all users", func(t *testing.T) {
		mockUserRepo.EXPECT().GetAll(gomock.Any()).Return(mockUsers, nil).Times(1)
		response, err := userService.GetAllUsers(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Users, 2)
	})
	t.Run("Get all users failed", func(t *testing.T) {
		mockUserRepo.EXPECT().GetAll(gomock.Any()).Return(nil, errors.New("failed to get users")).Times(1)
		response, err := userService.GetAllUsers(ctx)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "failed to get users", err.Error())
	})
}

func TestService_User_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repositories.NewMockUserRepository(ctrl)
	jwtService := auth.NewJWTService("cfg.JWTSecret", mockUserRepo)
	userService := NewUserService(mockUserRepo, jwtService)

	var (
		ctx   = context.Background()
		id    = primitive.NewObjectID()
		idNew = primitive.NewObjectID()
		now   = time.Now()
	)

	mockRequest := &dto.UpdateUserRequest{
		Name:  "Updated User",
		Email: "updated@example.com",
	}

	mockRequestEntity := &entities.User{
		ID:        idNew,
		Name:      mockRequest.Name,
		Email:     mockRequest.Email,
		Password:  "hashed_password",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockUserEntity := &entities.User{
		ID:        id,
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "hashed_password",
		CreatedAt: now,
		UpdatedAt: now,
	}
	t.Run("Update user success", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(gomock.Any(), id).Return(mockUserEntity, nil).Times(1)
		mockUserRepo.EXPECT().GetByEmail(gomock.Any(), mockRequest.Email).Return(nil, errors.New("email not found")).Times(1)
		mockUserRepo.EXPECT().Update(gomock.Any(), mockUserEntity).Return(nil).Times(1)
		response, err := userService.UpdateUser(ctx, id, mockRequest)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "Updated User", response.Name)
	})

	t.Run("Update user failed - user not found", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, errors.New("user not found")).Times(1)
		response, err := userService.UpdateUser(ctx, id, mockRequest)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "user not found", err.Error())
	})

	t.Run("Update user failed - email already exists", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(gomock.Any(), id).Return(mockUserEntity, nil).Times(1)
		mockUserRepo.EXPECT().GetByEmail(gomock.Any(), mockRequest.Email).Return(mockRequestEntity, nil).Times(1)
		response, err := userService.UpdateUser(ctx, id, mockRequest)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "user already exists", err.Error())
	})

	t.Run("Update user failed", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(gomock.Any(), id).Return(mockUserEntity, nil).Times(1)
		mockUserRepo.EXPECT().GetByEmail(gomock.Any(), mockRequest.Email).Return(nil, errors.New("email not found")).Times(1)
		mockUserRepo.EXPECT().Update(gomock.Any(), mockUserEntity).Return(errors.New("failed to update user")).Times(1)
		response, err := userService.UpdateUser(ctx, id, mockRequest)
		assert.Nil(t, response)
		assert.Equal(t, "failed to update user", err.Error())
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repositories.NewMockUserRepository(ctrl)
	jwtService := auth.NewJWTService("cfg.JWTSecret", mockUserRepo)
	userService := NewUserService(mockUserRepo, jwtService)

	var (
		ctx = context.Background()
		id  = primitive.NewObjectID()
		now = time.Now()
	)

	mockUserEntity := &entities.User{
		ID:        id,
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "hashed_password",
		CreatedAt: now,
		UpdatedAt: now,
	}

	t.Run("Delete user success", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(gomock.Any(), id).Return(mockUserEntity, nil).Times(1)
		mockUserRepo.EXPECT().Delete(gomock.Any(), id).Return(nil).Times(1)
		err := userService.DeleteUser(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("Delete user failed - user not found", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, errors.New("user not found")).Times(1)
		err := userService.DeleteUser(ctx, id)
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})

	t.Run("Delete user failed", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(gomock.Any(), id).Return(mockUserEntity, nil).Times(1)
		mockUserRepo.EXPECT().Delete(gomock.Any(), id).Return(errors.New("failed to delete user")).Times(1)
		err := userService.DeleteUser(ctx, id)
		assert.Error(t, err)
		assert.Equal(t, "failed to delete user", err.Error())
	})
}

func TestUserService_GetUserCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repositories.NewMockUserRepository(ctrl)
	jwtService := auth.NewJWTService("cfg.JWTSecret", mockUserRepo)
	userService := NewUserService(mockUserRepo, jwtService)

	var (
		ctx = context.Background()
	)

	t.Run("Get user count success", func(t *testing.T) {
		mockUserRepo.EXPECT().Count(gomock.Any()).Return(int64(1), nil).Times(1)
		count, err := userService.GetUserCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	t.Run("Get user count failed", func(t *testing.T) {
		mockUserRepo.EXPECT().Count(gomock.Any()).Return(int64(0), errors.New("failed to get user count")).Times(1)
		count, err := userService.GetUserCount(ctx)
		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
		assert.Equal(t, "failed to get user count", err.Error())
	})
}
