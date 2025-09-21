package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the validator with custom rules
type Validator struct {
	validator *validator.Validate
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	v := validator.New()

	// Register custom tag name function
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{validator: v}
}

// ValidateStruct validates a struct
func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validator.Struct(s)
}

// ValidateVar validates a single variable
func (v *Validator) ValidateVar(field interface{}, tag string) error {
	return v.validator.Var(field, tag)
}

// GetValidationErrors returns detailed validation errors
func (v *Validator) GetValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()

			switch tag {
			case "required":
				errors[field] = fmt.Sprintf("%s is required", field)
			case "email":
				errors[field] = fmt.Sprintf("%s must be a valid email", field)
			case "min":
				errors[field] = fmt.Sprintf("%s must be at least %s characters", field, e.Param())
			case "max":
				errors[field] = fmt.Sprintf("%s must be at most %s characters", field, e.Param())
			default:
				errors[field] = fmt.Sprintf("%s is invalid", field)
			}
		}
	}

	return errors
}

// Global validator instance
var globalValidator = NewValidator()

// ValidateStruct validates a struct using the global validator
func ValidateStruct(s interface{}) error {
	return globalValidator.ValidateStruct(s)
}
