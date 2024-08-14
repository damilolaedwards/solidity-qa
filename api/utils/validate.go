package utils

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func ValidateData(v interface{}) error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(v)

	if err != nil {
		var validationErrors validator.ValidationErrors
		ok := errors.As(err, &validationErrors)
		if !ok {
			// Handle unexpected validation error type
			return fmt.Errorf("unexpected validation error: %v", err)
		}

		// Construct an error message with validation errors
		var errorMsgs []string
		for _, e := range validationErrors {
			errorMsgs = append(errorMsgs, fmt.Sprintf("%s: %s", e.Field(), e.Tag()))
		}
		return fmt.Errorf("validation errors: %s", strings.Join(errorMsgs, "; "))
	}

	return nil
}
