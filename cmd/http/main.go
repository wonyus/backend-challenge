package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wonyus/backend-challenge/internal/application/services"
	"github.com/wonyus/backend-challenge/internal/infrastructure/auth"
	"github.com/wonyus/backend-challenge/internal/infrastructure/config"
	"github.com/wonyus/backend-challenge/internal/infrastructure/http/handlers"
	"github.com/wonyus/backend-challenge/internal/infrastructure/http/middleware"
	"github.com/wonyus/backend-challenge/internal/infrastructure/http/router"
	"github.com/wonyus/backend-challenge/internal/infrastructure/persistence/mongodb"
	"github.com/wonyus/backend-challenge/pkg/logger"
)

func main() {
	// Initialize logger
	logger := logger.New()
	logger.Info("Starting HTTP server...")

	// Load configuration
	cfg := config.Load()

	// Connect to MongoDB
	mongoClient, err := mongodb.NewConnection(cfg.MongoURI)
	if err != nil {
		logger.Error("Failed to connect to MongoDB:", err)
		os.Exit(1)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			logger.Error("Failed to disconnect from MongoDB:", err)
		}
	}()

	db := mongoClient.Database(cfg.DatabaseName)
	logger.Info("Connected to MongoDB successfully")

	// Initialize repositories
	userRepo := mongodb.NewUserRepository(db)

	// Initialize services
	jwtService := auth.NewJWTService(cfg.JWTSecret, userRepo)
	userService := services.NewUserService(userRepo, jwtService)
	authService := services.NewAuthService(userRepo, jwtService)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Initialize logging middleware
	loggingMiddleware := middleware.NewLoggingMiddleware(logger)

	// Initialize router
	r := router.NewRouter(userHandler, authHandler, authMiddleware, loggingMiddleware)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start background goroutine for user count logging
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				count, err := userService.GetUserCount(context.Background())
				if err != nil {
					logger.Error("Failed to get user count:", err)
				} else {
					logger.Info("Current user count:", count)
				}
			}
		}
	}()

	// Start server in a goroutine
	go func() {
		logger.Info("HTTP server starting on port", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server:", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown:", err)
		os.Exit(1)
	}

	logger.Info("Server exited")
}
