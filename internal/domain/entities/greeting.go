package entities

import "fmt"

// Person represents the person to be greeted.
type Person struct {
	FirstName string
	LastName  string
	Age       int64 // Added for potential validation examples
}

// Validate implements a simple validation for Person.
func (p Person) Validate() error {
	if p.FirstName == "" {
		return fmt.Errorf("first name cannot be empty")
	}
	if p.LastName == "" {
		return fmt.Errorf("last name cannot be empty")
	}
	if p.Age <= 0 {
		return fmt.Errorf("age must be a positive number")
	}
	return nil
}

// IsValid implements a simple validity check for Person.
func (p Person) IsValid() bool {
	return p.FirstName != "" && p.LastName != "" && p.Age > 0
}

// Greeting represents a generated greeting message.
type Greeting struct {
	Message string
	Person  Person
}
