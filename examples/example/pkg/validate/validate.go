package validate

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"go.uber.org/multierr"
)

func RegisterValidation(valid *validator.Validate) error {
	err := multierr.Combine(
		valid.RegisterValidation("mobile", ValidateMobile),
	)
	if err != nil {
		return fmt.Errorf("validator: register validation failed, %w", err)
	}
	return nil
}
