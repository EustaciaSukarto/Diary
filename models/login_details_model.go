package models

import (
	"github.com/go-playground/validator/v10"
)

type LoginDetails struct {
	Detail1 string `json:"email" | "username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func ValidateLoginDetails(details LoginDetails) (msg string) {
	// validate input rules
	validate := validator.New()
	er := validate.Struct(details)
	if er != nil {
		msg = er.(validator.ValidationErrors).Error()
		return
	}
	msg = ""
	return
}