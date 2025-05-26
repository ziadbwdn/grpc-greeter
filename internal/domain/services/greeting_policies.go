package services

import (
	"context"
	"fmt"
	"grpc-greeter/internal/domain/entities"
	"grpc-greeter/internal/infrastructure/logging" // Import the logger
)

// GreetingPolicy defines the interface for checking if a greeting is allowed.
type GreetingPolicy interface {
	IsGreetingAllowed(ctx context.Context, person entities.Person) (bool, error)
}

// GreetingAuditor defines the interface for recording greeting events.
type GreetingAuditor interface {
	RecordGreeting(ctx context.Context, greeting entities.Greeting) error
}

// greetingPolicy implements the GreetingPolicy interface.
type greetingPolicy struct {
	logger logging.Logger
}

// NewGreetingPolicy creates and returns a new GreetingPolicy instance.
func NewGreetingPolicy(logger logging.Logger) GreetingPolicy {
	return &greetingPolicy{
		logger: logger,
	}
}

// IsGreetingAllowed implements the GreetingPolicy interface.
// This is a placeholder policy; in a real app, it might check blacklists, rate limits, etc.
func (p *greetingPolicy) IsGreetingAllowed(ctx context.Context, person entities.Person) (bool, error) {
	p.logger.Debug("Checking greeting policy for person: %s %s", person.FirstName, person.LastName)
	// For demonstration, always allow greetings.
	return true, nil
}

// greetingAuditor implements the GreetingAuditor interface.
type greetingAuditor struct {
	logger logging.Logger
}

// NewGreetingAuditor creates and returns a new GreetingAuditor instance.
func NewGreetingAuditor(logger logging.Logger) GreetingAuditor {
	return &greetingAuditor{
		logger: logger,
	}
}

// RecordGreeting implements the GreetingAuditor interface.
// This is a placeholder; in a real app, it would persist the greeting to a database or log system.
func (a *greetingAuditor) RecordGreeting(ctx context.Context, greeting entities.Greeting) error {
	a.logger.Info("Auditing greeting: %s for %s %s", greeting.Message, greeting.Person.FirstName, greeting.Person.LastName)
	// Simulate saving to a persistent store
	fmt.Printf("AUDIT: Recorded greeting '%s' for %s %s\n", greeting.Message, greeting.Person.FirstName, greeting.Person.LastName)
	return nil
}
