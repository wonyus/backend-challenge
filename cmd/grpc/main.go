package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wonyus/backend-challenge/internal/application/services"
	"github.com/wonyus/backend-challenge/internal/infrastructure/auth"
	"github.com/wonyus/backend-challenge/internal/infrastructure/config"
	grpcHandlers "github.com/wonyus/backend-challenge/internal/infrastructure/grpc/handlers"
	pb "github.com/wonyus/backend-challenge/internal/infrastructure/grpc/proto"
	"github.com/wonyus/backend-challenge/internal/infrastructure/persistence/mongodb"
	"github.com/wonyus/backend-challenge/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Initialize logger
	logger := logger.New()
	logger.Info("Starting gRPC server...")

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

	// Initialize gRPC handlers
	userGRPCHandler := grpcHandlers.NewUserGRPCHandler(userService)

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	pb.RegisterUserServiceServer(grpcServer, userGRPCHandler)

	// Enable reflection for testing with tools like grpcurl
	reflection.Register(grpcServer)

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

	// Create listener
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		logger.Error("Failed to listen on port", cfg.GRPCPort, ":", err)
		os.Exit(1)
	}

	// Start server in a goroutine
	go func() {
		logger.Info("gRPC server starting on port", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("Failed to serve gRPC server:", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down gRPC server...")

	// Graceful shutdown
	grpcServer.GracefulStop()

	logger.Info("gRPC server exited")
}
