package utilities

import "fmt"

// This error is used to denote something risky and unexpected happening in the system.
// Ideally this should never throw, but if it does it means something really weird is going on.
// In most cases, BadErrors are untested.

type BadError struct {
	message string
}

func NewBadError(message string) *BadError {
	return &BadError{
		message: message,
	}
}

func (b *BadError) Error() string {
	return fmt.Sprintf("THIS IS BAD: %s", b.message)
}
