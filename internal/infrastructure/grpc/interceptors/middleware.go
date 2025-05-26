package interceptors

import (
	"context"
	"grpc-greeter/internal/infrastructure/logging" // Import the logger
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryValidationInterceptor performs basic request validation using domain IsValid/Validate methods if available.
func UnaryValidationInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// Validate the request if it implements the Validate() error interface
		if validatable, ok := req.(interface{ Validate() error }); ok {
			if err := validatable.Validate(); err != nil {
				// Use the injected logger instead of direct log.Printf
				logging.New().Error("Validation failed for %s: %v", info.FullMethod, err)
				return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
			}
		}

		// Otherwise check if it has IsValid() bool method (e.g., greeting)
		if checkable, ok := req.(interface{ IsValid() bool }); ok {
			if !checkable.IsValid() {
				logging.New().Error("Validation failed for %s: IsValid() returned false", info.FullMethod)
				return nil, status.Error(codes.InvalidArgument, "invalid request")
			}
		}

		// Proceed to actual handler
		return handler(ctx, req)
	}
}

// LoggingInterceptor logs incoming gRPC requests and their responses/errors.
func LoggingInterceptor(logger logging.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		logger.Info("Received gRPC request: Method=%s, Request=%+v", info.FullMethod, req)

		resp, err = handler(ctx, req)

		if err != nil {
			logger.Error("gRPC request failed: Method=%s, Error=%v", info.FullMethod, err)
		} else {
			logger.Info("gRPC request completed: Method=%s, Response=%+v", info.FullMethod, resp)
		}
		return resp, err
	}
}

// RecoveryInterceptor recovers from panics during gRPC request handling.
func RecoveryInterceptor(logger logging.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Recovered from panic during gRPC call: %v\nStack: %s", r, debug.Stack())
				err = status.Errorf(codes.Internal, "internal server error: %v", r)
			}
		}()
		return handler(ctx, req)
	}
}

// ChainUnaryServer is a helper function to chain multiple unary interceptors.
func ChainUnaryServer(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(interceptors...)
}
