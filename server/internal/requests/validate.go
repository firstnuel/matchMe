package requests

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

func (vs *ValidationService) ValidateRegisterUser(user RegisterUser) error {
	// First, validate struct tags
	if err := vs.validator.Struct(user); err != nil {
		return fmt.Errorf("struct validation failed: %w", err)
	}

	return nil
}

func (vs *ValidationService) Validate(a any) error {
	// First, validate struct tags
	if err := vs.validator.Struct(a); err != nil {
		return fmt.Errorf("struct validation failed: %w", err)
	}

	return nil
}

func (vs *ValidationService) ValidateUser(user UpdateUser) error {
	// First, validate struct tags
	if err := vs.validator.Struct(user); err != nil {
		return fmt.Errorf("struct validation failed: %w", err)
	}

	if user.Bio != nil {
		// Then, validate UserBio with custom logic
		if err := ValidateUserBio(*user.Bio); err != nil {
			return fmt.Errorf("user bio validation failed: %w", err)
		}
	}

	return nil
}
