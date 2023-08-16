// Package config error translation validator.
package config

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

// CreateEnTranslation configuration errortranslation English.
// nolint: nolintlint, ireturn
func CreateEnTranslation(v *validator.Validate) ut.Translator {
	uni := ut.New(en.New(), en.New())

	trans, _ := uni.GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(v, trans)

	return trans
}
