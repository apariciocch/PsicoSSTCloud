package validator
package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Validator wrapper para validación
type Validator struct {
	validate *validator.Validate
}

// New crea un nuevo validador
func New() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// ValidateStruct valida una estructura
func (v *Validator) ValidateStruct(data interface{}) error {
	err := v.validate.Struct(data)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var messages []string
			for _, fieldError := range validationErrors {
				messages = append(messages, fmt.Sprintf(
					"Field '%s' failed validation: %s",
					fieldError.Field(),
					fieldError.Tag(),
				))
			}
			return fmt.Errorf("validation errors: %v", messages)
		}
	}
	return err
}

// ValidateEmail valida un email
func (v *Validator) ValidateEmail(email string) error {
	return v.validate.Var(email, "email")
}

// ValidateMinLength valida longitud mínima
func (v *Validator) ValidateMinLength(s string, min int) error {
	if len(s) < min {
		return fmt.Errorf("minimum length is %d", min)
	}
	return nil
}

// ValidateMaxLength valida longitud máxima
func (v *Validator) ValidateMaxLength(s string, max int) error {
	if len(s) > max {
		return fmt.Errorf("maximum length is %d", max)
	}
	return nil
}
