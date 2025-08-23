package requests

import (
	"fmt"
	"strings"

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

// formatValidationError converts validator.ValidationErrors to user-friendly messages
func (vs *ValidationService) formatValidationError(err error) error {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return fmt.Errorf("validation failed: %s", err.Error())
	}

	var friendlyErrors []string
	
	for _, fieldError := range validationErrors {
		fieldName := strings.ToLower(fieldError.Field())
		
		// Convert field names to human-readable format
		switch fieldName {
		case "firstname":
			fieldName = "first name"
		case "lastname":
			fieldName = "last name"
		case "preferredagemin":
			fieldName = "preferred minimum age"
		case "preferredagemax":
			fieldName = "preferred maximum age"
		case "preferredgender":
			fieldName = "preferred gender"
		case "preferreddistance":
			fieldName = "preferred distance"
		case "aboutme":
			fieldName = "about me"
		case "lookingfor":
			fieldName = "looking for"
		case "musicpreferences":
			fieldName = "music preferences"
		case "foodpreferences":
			fieldName = "food preferences"
		case "communicationstyle":
			fieldName = "communication style"
		}

		var message string
		switch fieldError.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required", fieldName)
		case "email":
			message = fmt.Sprintf("%s must be a valid email address", fieldName)
		case "min":
			if fieldError.Kind().String() == "string" {
				message = fmt.Sprintf("%s must be at least %s characters long", fieldName, fieldError.Param())
			} else {
				message = fmt.Sprintf("%s must be at least %s", fieldName, fieldError.Param())
			}
		case "max":
			if fieldError.Kind().String() == "string" {
				message = fmt.Sprintf("%s must be no more than %s characters long", fieldName, fieldError.Param())
			} else {
				message = fmt.Sprintf("%s must be no more than %s", fieldName, fieldError.Param())
			}
		case "oneof":
			values := strings.ReplaceAll(fieldError.Param(), " ", ", ")
			message = fmt.Sprintf("%s must be one of: %s", fieldName, values)
		case "url":
			message = fmt.Sprintf("%s must be a valid URL", fieldName)
		case "dive":
			// For slice validation errors, provide generic message
			message = fmt.Sprintf("%s contains invalid values", fieldName)
		default:
			message = fmt.Sprintf("%s is invalid", fieldName)
		}
		
		friendlyErrors = append(friendlyErrors, message)
	}
	
	if len(friendlyErrors) == 1 {
		return fmt.Errorf("%s", friendlyErrors[0])
	}
	
	return fmt.Errorf("validation failed:\n• %s", strings.Join(friendlyErrors, "\n• "))
}

func (vs *ValidationService) ValidateRegisterUser(user RegisterUser) error {
	// Validate struct tags and convert to friendly errors
	if err := vs.validator.Struct(user); err != nil {
		return vs.formatValidationError(err)
	}

	return nil
}

func (vs *ValidationService) Validate(a any) error {
	// Validate struct tags and convert to friendly errors
	if err := vs.validator.Struct(a); err != nil {
		return vs.formatValidationError(err)
	}

	return nil
}

func (vs *ValidationService) ValidateUser(user UpdateUser) error {
	// Validate struct tags and convert to friendly errors
	if err := vs.validator.Struct(user); err != nil {
		return vs.formatValidationError(err)
	}

	if user.Bio != nil {
		// Then, validate UserBio with custom logic
		if err := ValidateUserBio(*user.Bio); err != nil {
			return err // ValidateUserBio already returns friendly errors
		}
	}

	return nil
}
