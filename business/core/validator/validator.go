// Package validator validation structs core models.
package validator

import (
	"fmt"
	"myc-devices-simulator/business/core/validator/errortranslation"
	"myc-devices-simulator/business/sys/errors"

	"github.com/go-playground/validator/v10"
)

// IsValid validates data structures and returns validation errors.
func IsValid(data interface{}) ([]string, error) {
	validate := newValidator()

	if err := validate.Struct(data); err != nil {
		messageError := errortranslation.TranslateEnError(err, validate)

		return messageError, fmt.Errorf("core.validator.IsValid(%v): %w - list: %v",
			data, errors.ErrValidatorInvalidCoreModel, messageError)
	}

	return nil, nil
}

// newValidator create a new validator.
func newValidator() *validator.Validate {
	newValidator := validator.New()

	return newValidator
}
