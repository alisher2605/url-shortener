package validator

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	Validator *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	return &Validator{
		Validator: v,
	}
}
