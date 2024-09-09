package utils

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
)

var decoder *form.Decoder

// DecodeAndValidateFormData is a utility function that parses form data from an HTTP request,
// populates the provided struct with the form values, and validates the data.
func DecodeAndValidateFormData(r *http.Request, data interface{}) error {
	// Parse the form data
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("failed to parse form: %w", err)
	}

	// Decode form data into the provided struct
	if err := decodeFormData(r.Form, data); err != nil {
		return fmt.Errorf("failed to decode form data: %w", err)
	}

	// Validate the decoded data
	if err := ValidateData(data); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}

// decodeFormData populates the provided struct with form values.
func decodeFormData(formData url.Values, data interface{}) error {
	decoder = form.NewDecoder()
	if err := decoder.Decode(data, formData); err != nil {
		return fmt.Errorf("failed to decode form data: %w", err)
	}

	return nil
}

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
