package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wonyus/backend-challenge/internal/application/dto"
	"github.com/wonyus/backend-challenge/internal/domain/entities"
	domainErrors "github.com/wonyus/backend-challenge/internal/domain/errors"
	"github.com/wonyus/backend-challenge/internal/infrastructure/auth"
	mock_repositories "github.com/wonyus/backend-challenge/mock/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestService_Auth_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repositories.NewMockUserRepository(ctrl)
	jwtService := auth.NewJWTService("cfg.JWTSecret", mockUserRepo)
	AuthService := NewAuthService(mockUserRepo, jwtService)

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

	mockUserEntity := &entities.User{
		ID:        id,
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "hashed_password",
		CreatedAt: now,
		UpdatedAt: now,
	}

	t.Run("Register user success", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmail(gomock.Any(), mockRequest.Email).Return(nil, nil).Times(1)
		mockUserRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		response, err := AuthService.Register(ctx, mockRequest)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		assert.Equal(t, "User registered successfully", response.Message)
	})

	t.Run("Register user already exists", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmail(gomock.Any(), mockRequest.Email).Return(mockUserEntity, nil).Times(1)
		response, err := AuthService.Register(ctx, mockRequest)
		assert.Nil(t, response)
		assert.Equal(t, "user already exists", err.Error())
	})

	t.Run("Register user creation error", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmail(gomock.Any(), mockRequest.Email).Return(nil, nil).Times(1)
		mockUserRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("failed to create user")).Times(1)
		response, err := AuthService.Register(ctx, mockRequest)
		assert.Nil(t, response)
		assert.Equal(t, "failed to create user", err.Error())
	})

	t.Run("Register user entity creation error", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmail(gomock.Any(), mockRequest.Email).Return(nil, nil).Times(1)
		mockRequest.Name = ""
		response, err := AuthService.Register(ctx, mockRequest)
		assert.Nil(t, response)
		assert.Equal(t, "name, email, and password are required", err.Error())
	})

	t.Run("Register user password hashing error", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmail(gomock.Any(), mockRequest.Email).Return(nil, nil).Times(1)
		mockRequest.Password = ""
		response, err := AuthService.Register(ctx, mockRequest)
		assert.Nil(t, response)
		assert.Equal(t, "password cannot be empty", err.Error())
	})

}

func TestService_Auth_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repositories.NewMockUserRepository(ctrl)
	jwtService := auth.NewJWTService("cfg.JWTSecret", mockUserRepo)
	AuthService := NewAuthService(mockUserRepo, jwtService)

	var (
		ctx          = context.Background()
		id           = primitive.NewObjectID()
		password     = "password"
		passwordHash = "$2a$06$R.ga34oljt5UqXmSgNR6ze4QpEbq8u9i0Fui/eG2WpZs/nCgjbT1e"
		now          = time.Now()
	)

	mockRequest := &dto.LoginRequest{
		Email:    "test@example.com",
		Password: password,
	}

	mockUserEntity := &entities.User{
		ID:        id,
		Name:      "Test User",
		Email:     mockRequest.Email,
		Password:  passwordHash,
		CreatedAt: now,
		UpdatedAt: now,
	}

	t.Run("Login Success", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmail(gomock.Any(), mockRequest.Email).Return(mockUserEntity, nil).Times(1)
		response, err := AuthService.Login(ctx, mockRequest)
		assert.NotNil(t, response)
		assert.NoError(t, err)
		fmt.Println(response)
	})

	t.Run("Login User Not Found", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmail(gomock.Any(), mockRequest.Email).Return(nil, errors.New("user not found")).Times(1)
		response, err := AuthService.Login(ctx, mockRequest)
		assert.Nil(t, response)
		assert.Equal(t, domainErrors.ErrInvalidCredentials.Error(), err.Error())
	})

	t.Run("Login Token Generation Error", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmail(gomock.Any(), mockRequest.Email).Return(mockUserEntity, nil).Times(1)
		jwtService = auth.NewJWTService("", mockUserRepo) // Empty secret to trigger error
		AuthService = NewAuthService(mockUserRepo, jwtService)
		response, err := AuthService.Login(ctx, mockRequest)
		assert.Nil(t, response)
		assert.Equal(t, domainErrors.ErrInvalidTokenSecret.Error(), err.Error())
	})

	t.Run("Login Invalid Password", func(t *testing.T) {
		mockUserEntity.Password = password
		mockUserRepo.EXPECT().GetByEmail(gomock.Any(), mockRequest.Email).Return(mockUserEntity, nil).Times(1)
		response, err := AuthService.Login(ctx, mockRequest)
		assert.Nil(t, response)
		assert.Equal(t, domainErrors.ErrInvalidCredentials.Error(), err.Error())
	})

}

func TestService_Auth_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repositories.NewMockUserRepository(ctrl)
	jwtService := auth.NewJWTService("cfg.JWTSecret", mockUserRepo)
	AuthService := NewAuthService(mockUserRepo, jwtService)

	var (
		ctx = context.Background()
		id  = primitive.NewObjectID()
		// password     = "password"
		passwordHash = "$2a$06$R.ga34oljt5UqXmSgNR6ze4QpEbq8u9i0Fui/eG2WpZs/nCgjbT1e"
		now          = time.Now()
	)

	mockUserEntity := &entities.User{
		ID:        id,
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  passwordHash,
		CreatedAt: now,
	}

	t.Run("Validate Token", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByID(gomock.Any(), mockUserEntity.ID).Return(mockUserEntity, nil).Times(1)
		token, err := jwtService.GenerateToken(ctx, mockUserEntity)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		response, err := AuthService.ValidateToken(ctx, token)
		assert.NoError(t, err)
		assert.NotNil(t, response)
	})

	t.Run("Validate Token Invalid Token", func(t *testing.T) {
		response, err := AuthService.ValidateToken(ctx, "invalid_token")
		assert.Nil(t, response)
		assert.Equal(t, domainErrors.ErrInvalidToken.Error(), err.Error())
	})

}
