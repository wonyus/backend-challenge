package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wonyus/backend-challenge/internal/application/dto"
	mock_ports "github.com/wonyus/backend-challenge/mock/port"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestHandler_Auth_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_ports.NewMockAuthService(ctrl)

	var (
		ctx = context.Background()
		id  = primitive.NewObjectID()
	)

	mockRequest := &dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockResponse := &dto.RegisterResponse{
		ID:      id,
		Message: "User registered successfully",
	}
	executeWithRequest := func(method string, jsonRequestBody []byte) *httptest.ResponseRecorder {
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(method, "/api/auth/register", strings.NewReader(string(jsonRequestBody)))
		authHandler := NewAuthHandler(mockAuthService)
		mux := http.NewServeMux()
		mux.HandleFunc("/api/auth/register", authHandler.Register)
		mux.ServeHTTP(response, req)
		return response
	}

	t.Run("Success", func(t *testing.T) {
		jsonBody := []byte(`{
			"email": "test@example.com",
			"password": "password123",
			"name": "Test User"
		}`)

		mockAuthService.EXPECT().Register(ctx, mockRequest).Return(mockResponse, nil)
		response := executeWithRequest(http.MethodPost, jsonBody)
		assert.Equal(t, http.StatusCreated, response.Code)
	})

	t.Run("User Already Exists", func(t *testing.T) {
		jsonBody := []byte(`{
			"email": "test@example.com",
			"password": "password123",
			"name": "Test User"
		}`)

		mockAuthService.EXPECT().Register(ctx, mockRequest).Return(nil, errors.New("user already exists"))
		response := executeWithRequest(http.MethodPost, jsonBody)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		jsonBody := []byte(`{
			"email": "test@example.com",
			"password": "password123",
		}`)

		response := executeWithRequest(http.MethodPost, jsonBody)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Bad Request", func(t *testing.T) {
		jsonBody := []byte(`{}`)

		response := executeWithRequest(http.MethodPost, jsonBody)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

}

func TestHandler_Auth_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_ports.NewMockAuthService(ctrl)

	var (
		ctx = context.Background()
		id  = primitive.NewObjectID()
	)

	mockRequest := &dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	mockResponse := &dto.LoginResponse{
		Token: "mocked_token",
		User: dto.UserResponse{
			ID:    id,
			Email: "test@example.com",
			Name:  "Test User",
		},
	}

	executeWithRequest := func(method string, jsonRequestBody []byte) *httptest.ResponseRecorder {
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(method, "/api/auth/login", strings.NewReader(string(jsonRequestBody)))
		authHandler := NewAuthHandler(mockAuthService)
		mux := http.NewServeMux()
		mux.HandleFunc("/api/auth/login", authHandler.Login)
		mux.ServeHTTP(response, req)
		return response
	}

	t.Run("Success", func(t *testing.T) {
		jsonBody := []byte(`{
			"email": "test@example.com",
			"password": "password123"
		}`)

		mockAuthService.EXPECT().Login(ctx, mockRequest).Return(mockResponse, nil)
		response := executeWithRequest(http.MethodPost, jsonBody)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		jsonBody := []byte(`{
			"email": "test@example.com",
			"password": "password123"
		}`)

		mockAuthService.EXPECT().Login(ctx, mockRequest).Return(nil, errors.New("invalid credentials"))
		response := executeWithRequest(http.MethodPost, jsonBody)
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		jsonBody := []byte(`{
			"email": "test@example.com",
		}`)

		response := executeWithRequest(http.MethodPost, jsonBody)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Bad Request", func(t *testing.T) {
		jsonBody := []byte(`{}`)

		response := executeWithRequest(http.MethodPost, jsonBody)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

}
