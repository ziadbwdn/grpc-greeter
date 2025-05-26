package main

import (
	"fmt"
	"log" // Standard log for fatal errors before custom logger is initialized
	"net"
	"os" // For os.Exit

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"grpc-greeter/internal/application/usecases"
	"grpc-greeter/internal/config"
	"grpc-greeter/internal/domain/services"
	"grpc-greeter/internal/infrastructure/grpc/handlers"
	"grpc-greeter/internal/infrastructure/grpc/interceptors"
	"grpc-greeter/internal/infrastructure/logging"
	pb "grpc-greeter/pkg/proto/generated"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		// Use standard log.Fatalf here as the custom logger might not be ready yet
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize custom logger
	logger := logging.New()
	logger.Info("Starting gRPC server...")

	// Setup domain services and usecase
	// Pass the logger instance to service constructors as defined in their packages
	greetingService := services.NewGreetingService(logger)
	greetingPolicy := services.NewGreetingPolicy(logger)
	greetingAuditor := services.NewGreetingAuditor(logger)

	// Assuming NewGreetingUseCase and NewGreeterHandler are defined to accept these dependencies
	usecase := usecases.NewGreetingUseCase(greetingService, greetingPolicy, greetingAuditor)
	handler := handlers.NewGreeterHandler(usecase)

	// Create listener
	portStr := fmt.Sprintf(":%d", cfg.GRPCPort)
	lis, err := net.Listen("tcp", portStr)
	if err != nil {
		// Use custom logger for errors after it's initialized, then exit
		logger.Error("failed to listen: %v", err)
		os.Exit(1)
	}

	// Setup interceptors
	// ChainUnaryServer returns a grpc.ServerOption that bundles multiple interceptors
	unaryInterceptorOption := interceptors.ChainUnaryServer(
		interceptors.LoggingInterceptor(logger),   // Pass logger to LoggingInterceptor
		interceptors.RecoveryInterceptor(logger),  // Pass logger to RecoveryInterceptor
		interceptors.UnaryValidationInterceptor(), // Add the validation interceptor
	)

	// Create gRPC server with the chained interceptors
	// Pass the unaryInterceptorOption directly to grpc.NewServer
	grpcServer := grpc.NewServer(unaryInterceptorOption)

	// Register service
	pb.RegisterGreeterServer(grpcServer, handler)
	reflection.Register(grpcServer)

	// Start serving
	logger.Info("gRPC server listening on port %d", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		// Use custom logger for errors, then exit
		logger.Error("failed to serve: %v", err)
		os.Exit(1)
	}
}
