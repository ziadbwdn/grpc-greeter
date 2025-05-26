package services

import (
	"context"
	"fmt"
	"grpc-greeter/internal/domain/entities"
	"grpc-greeter/internal/infrastructure/logging" // Import the logger
)

// GreetingService defines the interface for creating greetings.
type GreetingService interface {
	CreateGreeting(ctx context.Context, person entities.Person) (entities.Greeting, error)
}

// GreetingFormatter defines the interface for formatting greeting messages.
type GreetingFormatter interface {
	FormatGreeting(person entities.Person) string
}

// greetingService implements the GreetingService and GreetingFormatter interfaces.
type greetingService struct {
	logger logging.Logger // Dependency on the logger
}

// NewGreetingService creates and returns a new GreetingService instance.
func NewGreetingService(logger logging.Logger) GreetingService {
	return &greetingService{
		logger: logger,
	}
}

// CreateGreeting implements the GreetingService interface.
// It generates a greeting for the given person.
func (s *greetingService) CreateGreeting(ctx context.Context, person entities.Person) (entities.Greeting, error) {
	s.logger.Info("Creating greeting for person: %s %s", person.FirstName, person.LastName)

	// In a real application, this might involve more complex logic,
	// database interaction, or external API calls.
	greetingText := s.FormatGreeting(person)

	greeting := entities.Greeting{
		Message: greetingText,
		Person:  person,
	}

	s.logger.Info("Generated greeting: %s", greeting.Message)
	return greeting, nil
}

// FormatGreeting implements the GreetingFormatter interface.
// It formats a greeting string.
func (s *greetingService) FormatGreeting(person entities.Person) string {
	return fmt.Sprintf("Hello, %s %s!", person.FirstName, person.LastName)
}
