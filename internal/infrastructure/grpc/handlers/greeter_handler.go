package handlers

import (
	"context" // Import for fmt.Errorf

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"grpc-greeter/internal/application/usecases"
	"grpc-greeter/internal/domain/entities" // Import entities to create Person
	pb "grpc-greeter/pkg/proto/generated"
)

// GreeterHandler implements the gRPC GreeterServer interface.
type GreeterHandler struct {
	// The usecase field should hold the interface directly, not a pointer to it.
	usecase usecases.GreetingUseCase
	// This embeds the unimplemented server, which is good practice.
	pb.UnimplementedGreeterServer
}

// NewGreeterHandler creates and returns a new GreeterHandler instance.
// It accepts the GreetingUseCase interface.
func NewGreeterHandler(uc usecases.GreetingUseCase) *GreeterHandler {
	return &GreeterHandler{usecase: uc}
}

// SayHello implements the SayHello RPC method from the Greeter service.
func (h *GreeterHandler) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	// Create an entities.Person from the HelloRequest.
	// In a real application, you might have more fields in HelloRequest
	// to populate the Person entity more fully.
	person := entities.Person{
		FirstName: req.GetName(), // Assuming GetName() provides the first name
		LastName:  "",            // Placeholder, as HelloRequest only has Name
		Age:       0,             // Placeholder, as HelloRequest doesn't have age
	}

	// Call the Execute method on the usecase interface.
	// The usecase.Execute method expects an entities.Person.
	greeting, err := h.usecase.Execute(ctx, person)
	if err != nil {
		// Return a gRPC status error for better client error handling
		return nil, status.Errorf(codes.Internal, "failed to process greeting request: %v", err)
	}

	// Return the greeting message in the HelloReply.
	return &pb.HelloReply{Message: greeting.Message}, nil
}
