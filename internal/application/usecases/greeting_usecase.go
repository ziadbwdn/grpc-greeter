package usecases

import (
	"context"
	"fmt"
	"grpc-greeter/internal/domain/entities"
	"grpc-greeter/internal/domain/services"
	"grpc-greeter/internal/infrastructure/logging" // Assuming logger might be useful here for debugging/logging use case logic
)

// GreetingUseCase defines the interface for the greeting use case.
type GreetingUseCase interface {
	Execute(ctx context.Context, person entities.Person) (entities.Greeting, error)
}

// greetingUseCase implements the GreetingUseCase interface.
type greetingUseCase struct {
	greetingService services.GreetingService
	greetingPolicy  services.GreetingPolicy
	greetingAuditor services.GreetingAuditor
	logger          logging.Logger // Optional: add logger for use case specific logging
}

// NewGreetingUseCase creates and returns a new GreetingUseCase instance.
func NewGreetingUseCase(
	greetingService services.GreetingService,
	greetingPolicy services.GreetingPolicy,
	greetingAuditor services.GreetingAuditor,
) GreetingUseCase {
	// You might want to pass the logger here from main.go if you need use case specific logging
	return &greetingUseCase{
		greetingService: greetingService,
		greetingPolicy:  greetingPolicy,
		greetingAuditor: greetingAuditor,
		logger:          logging.New(), // Using the singleton logger for simplicity
	}
}

// Execute implements the GreetingUseCase interface.
// It orchestrates the creation, policy check, and auditing of a greeting.
func (uc *greetingUseCase) Execute(ctx context.Context, person entities.Person) (entities.Greeting, error) {
	uc.logger.Info("Executing greeting use case for person: %s %s", person.FirstName, person.LastName)

	// 1. Apply greeting policy
	allowed, err := uc.greetingPolicy.IsGreetingAllowed(ctx, person)
	if err != nil {
		uc.logger.Error("Error checking greeting policy: %v", err)
		return entities.Greeting{}, err
	}
	if !allowed {
		uc.logger.Info("Greeting not allowed for person: %s %s", person.FirstName, person.LastName)
		return entities.Greeting{}, fmt.Errorf("greeting not allowed for this person")
	}

	// 2. Create the greeting
	greeting, err := uc.greetingService.CreateGreeting(ctx, person)
	if err != nil {
		uc.logger.Error("Error creating greeting: %v", err)
		return entities.Greeting{}, err
	}

	// 3. Record the greeting for auditing
	err = uc.greetingAuditor.RecordGreeting(ctx, greeting)
	if err != nil {
		uc.logger.Error("Error recording greeting for auditing: %v", err)
		// Decide if this error should block the response or just be logged
		// For now, we'll return it.
		return entities.Greeting{}, err
	}

	uc.logger.Info("Greeting use case completed successfully for person: %s %s", person.FirstName, person.LastName)
	return greeting, nil
}
