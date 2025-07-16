package router

import (
	"github.com/gorilla/mux"
	"github.com/wonyus/backend-challenge/internal/infrastructure/http/handlers"
	"github.com/wonyus/backend-challenge/internal/infrastructure/http/middleware"
)

func NewRouter(userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler, authMiddleware *middleware.AuthMiddleware, loggingMiddleware *middleware.LoggingMiddleware) *mux.Router {
	r := mux.NewRouter()

	// Apply logging middleware to all routes
	r.Use(loggingMiddleware.Middleware)

	// API prefix
	api := r.PathPrefix("/api").Subrouter()

	// Auth routes (public)
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", authHandler.Register).Methods("POST")
	auth.HandleFunc("/login", authHandler.Login).Methods("POST")

	// User routes (protected)
	users := api.PathPrefix("/users").Subrouter()
	users.Use(authMiddleware.Authenticate)
	users.HandleFunc("", userHandler.CreateUser).Methods("POST")
	users.HandleFunc("", userHandler.GetAllUsers).Methods("GET")
	users.HandleFunc("/{id}", userHandler.GetUser).Methods("GET")
	users.HandleFunc("/{id}", userHandler.UpdateUser).Methods("PUT")
	users.HandleFunc("/{id}", userHandler.DeleteUser).Methods("DELETE")

	return r
}
