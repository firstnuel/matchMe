package models

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationService struct {
	validator *validator.Validate
}

// NewValidationService creates a new validation service
func NewValidationService() *ValidationService {
	return &ValidationService{
		validator: validator.New(),
	}
}

func (vs *ValidationService) ValidateUser(user User) error {
	// First, validate struct tags
	if err := vs.validator.Struct(user); err != nil {
		return fmt.Errorf("struct validation failed: %w", err)
	}

	// Then, validate UserBio with custom logic
	if err := ValidateUserBio(user.Bio); err != nil {
		return fmt.Errorf("user bio validation failed: %w", err)
	}

	return nil
}
