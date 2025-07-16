package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/wonyus/backend-challenge/internal/application/dto"
	mock_ports "github.com/wonyus/backend-challenge/mock/port"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestHandler_User_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_ports.NewMockUserService(ctrl)

	var (
		ctx = context.Background()
		id  = primitive.NewObjectID()
	)

	mockRequest := &dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockResponse := &dto.UserResponse{
		ID:    id,
		Name:  "Test User",
		Email: "test@example.com",
	}

	executeWithRequest := func(method string, jsonRequestBody []byte) *httptest.ResponseRecorder {
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(method, "/api/users", strings.NewReader(string(jsonRequestBody)))
		userHandler := NewUserHandler(mockUserService)
		mux := http.NewServeMux()
		mux.HandleFunc("/api/users", userHandler.CreateUser)
		mux.ServeHTTP(response, req)
		return response
	}

	t.Run("Success", func(t *testing.T) {
		jsonBody := []byte(`{
			"name": "Test User",
			"email": "test@example.com",
			"password": "password123"
		}`)
		mockUserService.EXPECT().CreateUser(ctx, mockRequest).Return(mockResponse, nil)
		response := executeWithRequest(http.MethodPost, jsonBody)
		assert.Equal(t, http.StatusCreated, response.Code)
	})

	t.Run("User Already Exists", func(t *testing.T) {
		jsonBody := []byte(`{
			"name": "Test User",
			"email": "test@example.com",
			"password": "password123"
		}`)
		mockUserService.EXPECT().CreateUser(ctx, mockRequest).Return(nil, errors.New("user already exists"))
		response := executeWithRequest(http.MethodPost, jsonBody)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		jsonBody := []byte(`{
			"name": "Test User",
			"email": "test@example.com"
		}`)
		response := executeWithRequest(http.MethodPost, jsonBody)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Bad Request", func(t *testing.T) {
		jsonBody := []byte(`{
			"name": "Test User",
			"email": "test@example.com",
		}`)
		response := executeWithRequest(http.MethodPost, jsonBody)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

}

func TestHandler_User_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_ports.NewMockUserService(ctrl)

	var (
		// ctx = context.Background()
		id = primitive.NewObjectID()
	)

	mockResponse := &dto.UserResponse{
		ID:    id,
		Name:  "Test User",
		Email: "test@example.com",
	}

	executeWithRequest := func(method string, userID string, vars map[string]string) *httptest.ResponseRecorder {
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(method, fmt.Sprintf("/api/users/%s", userID), nil)
		userHandler := NewUserHandler(mockUserService)
		s := http.NewServeMux()
		s.HandleFunc("/api/users/{id}", userHandler.GetUser)
		req = mux.SetURLVars(req, vars)
		s.ServeHTTP(response, req)
		return response
	}
	t.Run("Success", func(t *testing.T) {
		mockUserService.EXPECT().GetUserByID(gomock.Any(), id).Return(mockResponse, nil)
		vars := map[string]string{"id": id.Hex()}
		response := executeWithRequest(http.MethodGet, id.Hex(), vars)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Invalid User ID", func(t *testing.T) {
		vars := map[string]string{"ids": "invalid-id"}
		response := executeWithRequest(http.MethodGet, "invalid-id", vars)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockUserService.EXPECT().GetUserByID(gomock.Any(), id).Return(nil, errors.New("user not found"))
		vars := map[string]string{"id": id.Hex()}
		response := executeWithRequest(http.MethodGet, id.Hex(), vars)
		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

func TestHandler_User_GetAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_ports.NewMockUserService(ctrl)

	var (
		ctx = context.Background()
	)
	mockResponse := &dto.UsersListResponse{
		Users: []dto.UserResponse{
			{
				ID:    primitive.NewObjectID(),
				Name:  "Test User 1",
				Email: "test1@example.com",
			},
			{
				ID:    primitive.NewObjectID(),
				Name:  "Test User 2",
				Email: "test2@example.com",
			},
		},
		Total: 2,
	}

	executeWithRequest := func(method string) *httptest.ResponseRecorder {
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(method, fmt.Sprintf("/api/users"), nil)
		userHandler := NewUserHandler(mockUserService)
		s := http.NewServeMux()
		s.HandleFunc("/api/users", userHandler.GetAllUsers)
		s.ServeHTTP(response, req)
		return response
	}

	t.Run("Success", func(t *testing.T) {
		mockUserService.EXPECT().GetAllUsers(ctx).Return(mockResponse, nil)
		response := executeWithRequest(http.MethodGet)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockUserService.EXPECT().GetAllUsers(ctx).Return(nil, errors.New("internal server error"))
		response := executeWithRequest(http.MethodGet)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestHandler_User_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_ports.NewMockUserService(ctrl)

	var (
		id         = primitive.NewObjectID()
		idNotFound = primitive.NewObjectID()
	)
	mockRequest := &dto.UpdateUserRequest{
		Name:  "Updated User",
		Email: "Test@update.com",
	}

	mockResponse := &dto.UserResponse{
		ID:    id,
		Name:  "Updated User",
		Email: "Test@update.com",
	}

	executeWithRequest := func(method string, userID string, vars map[string]string, jsonBody []byte) *httptest.ResponseRecorder {
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(method, fmt.Sprintf("/api/users/%s", userID), strings.NewReader(string(jsonBody)))
		userHandler := NewUserHandler(mockUserService)
		s := http.NewServeMux()
		s.HandleFunc("/api/users/{id}", userHandler.UpdateUser)
		req = mux.SetURLVars(req, vars)
		s.ServeHTTP(response, req)
		return response
	}
	t.Run("Success", func(t *testing.T) {
		jsonBody := []byte(`{
			"name": "Updated User",
			"email": "Test@update.com"
		}`)

		mockUserService.EXPECT().UpdateUser(gomock.Any(), id, mockRequest).Return(mockResponse, nil)
		vars := map[string]string{"id": id.Hex()}
		response := executeWithRequest(http.MethodGet, id.Hex(), vars, jsonBody)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Invalid id", func(t *testing.T) {
		jsonBody := []byte(`{
			"name": "Updated User",
			"email": "Test@update.com"
		}`)
		vars := map[string]string{"ids": "invalid-id"}
		response := executeWithRequest(http.MethodGet, "invalid-id", vars, jsonBody)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Invalid body", func(t *testing.T) {
		jsonBody := []byte(`{
			"name": "Updated User",
			"email": "Test@update.com",
		}`)
		vars := map[string]string{"id": id.Hex()}
		response := executeWithRequest(http.MethodGet, id.Hex(), vars, jsonBody)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		jsonBody := []byte(`{
			"name": "U",
			"email": "Test@update.com"
		}`)

		vars := map[string]string{"id": id.Hex()}
		response := executeWithRequest(http.MethodGet, id.Hex(), vars, jsonBody)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("User Not Found", func(t *testing.T) {
		jsonBody := []byte(`{
			"name": "Updated User",
			"email": "Test@update.com"
		}`)

		mockUserService.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), mockRequest).Return(nil, errors.New("user not found"))
		vars := map[string]string{"id": idNotFound.Hex()}
		response := executeWithRequest(http.MethodGet, idNotFound.Hex(), vars, jsonBody)
		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

func TestHandler_User_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_ports.NewMockUserService(ctrl)

	var (
		id         = primitive.NewObjectID()
		idNotFound = primitive.NewObjectID()
	)

	executeWithRequest := func(method string, userID string, vars map[string]string) *httptest.ResponseRecorder {
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(method, fmt.Sprintf("/api/users/%s", userID), nil)
		userHandler := NewUserHandler(mockUserService)
		s := http.NewServeMux()
		s.HandleFunc("/api/users/{id}", userHandler.DeleteUser)
		req = mux.SetURLVars(req, vars)
		s.ServeHTTP(response, req)
		return response
	}
	t.Run("Success", func(t *testing.T) {
		mockUserService.EXPECT().DeleteUser(gomock.Any(), id).Return(nil)
		vars := map[string]string{"id": id.Hex()}
		response := executeWithRequest(http.MethodDelete, id.Hex(), vars)
		assert.Equal(t, http.StatusNoContent, response.Code)
	})

	t.Run("Invalid Parameter", func(t *testing.T) {
		vars := map[string]string{"ids": "invalid-id"}
		response := executeWithRequest(http.MethodDelete, "invalid-id", vars)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockUserService.EXPECT().DeleteUser(gomock.Any(), idNotFound).Return(errors.New("user not found"))
		vars := map[string]string{"id": idNotFound.Hex()}
		response := executeWithRequest(http.MethodDelete, idNotFound.Hex(), vars)
		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}
