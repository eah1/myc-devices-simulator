// Package errortranslation contains a error translation errors to text.
package errortranslation

import (
	"errors"
	"myc-devices-simulator/business/core/validator/config"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// TranslateEnError translate english error of validator structs.
func TranslateEnError(err error, v *validator.Validate) (errs []string) {
	trans := config.CreateEnTranslation(v)

	return translateError(err, trans)
}

// translateError translate error of validator structs.
func translateError(err error, trans ut.Translator) (errs []string) {
	if err != nil {
		var validatorErrs validator.ValidationErrors

		if errors.As(err, &validatorErrs) {
			for _, e := range validatorErrs {
				//nolint:goerr113
				translatedErr := errors.New(e.Translate(trans))
				errs = append(errs, translatedErr.Error())
			}
		}

		return errs
	}

	return nil
}
