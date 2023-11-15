package token

import (
	"fmt"
	errorssys "myc-devices-simulator/business/sys/errors"

	"github.com/m1/go-generate-password/generator"
)

// Generate token with validations.
func Generate(length uint) (*string, error) {
	// Token Generation.
	validationConfig := new(generator.Config)
	validationConfig.Length = length
	validationConfig.IncludeSymbols = false
	validationConfig.IncludeNumbers = true
	validationConfig.IncludeLowercaseLetters = true
	validationConfig.IncludeUppercaseLetters = true
	validationConfig.ExcludeSimilarCharacters = false
	validationConfig.ExcludeAmbiguousCharacters = false

	// Generate token.
	gen, err := generator.New(validationConfig)
	if err != nil {
		return nil, fmt.Errorf("token.Generate.New(%v) - error: {%w}", validationConfig, errorssys.ErrTokenGenerating)
	}

	// Create token.
	token, err := gen.Generate()
	if err != nil {
		return nil, fmt.Errorf("token.Generate.Generate(-) - error: {%w}", errorssys.ErrTokenGenerating)
	}

	return token, nil
}
