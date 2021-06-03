package models

import (
	"github.com/go-playground/validator/v10"
)

type Entry struct {
	Date 	string 	`"json:"date" validate:"required"`
	UserID 	int 	`"json:userid" validate:"isdefault"`
	Content string 	`"json:content"`
}

func ValidateEntry(entry Entry) (msg string) {
	// validate input rules
	validate := validator.New()
	er := validate.Struct(entry)
	if er != nil {
		msg = er.(validator.ValidationErrors).Error()
		return
	}
	msg = ""
	return
}