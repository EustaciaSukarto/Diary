package models

import (
	"github.com/go-playground/validator/v10"
	"books/db"
	"unicode"
)

type User struct {
	ID     		int  	`json:"id" validate:"isdefault"`
	Fullname   	string 	`json:"fullname"`
	Birthday  	string 	`json:"birthday"`
	Email 		string 	`json:"email" validate:"required,email"`
	Username 	string 	`json:"username" validate:"required"`
	Password 	string 	`json:"password" validate:"password-rules"`
}

func PasswordValidator(v validator.FieldLevel) bool {
	//		has: 	upper  lower  number char
	rules := []bool{false, false, false, false}
	pw := v.Field().String()
	if (len(pw) < 6 || len(pw) > 32) {
		return false
	}
	for _, character := range pw {
		if (unicode.IsUpper(character)) {
			rules[0] = true
		}
		if (unicode.IsLower(character)) {
			rules[1] = true
		}
		if (unicode.IsDigit(character)) {
			rules[2] = true
		}
		if (unicode.IsSymbol(character)) {
			rules[3] = true
		}
	}
	if (rules[0] == true && rules[1] == true && rules[2] == true && rules[3] == true) {
		return true
	}
	return false
}

func QueryUser(email string) (user User, err error) {
	err = db.DB.QueryRow("SELECT * FROM user WHERE Email=?", user.Email).Scan(&user.ID, &user.Fullname, &user.Birthday, &user.Email, &user.Username, &user.Password)
	return
}

func ValidateUser(user User) (msg string) {
	// check if user already exists
	dbuser, err := QueryUser(user.Email)
	if err != nil {
		return err.Error()
	}
	if dbuser.Email != "" {
		msg = "Email already in use."
		return
	}
	// validate input rules
	validate := validator.New()
	validate.RegisterValidation("password-rules", PasswordValidator)
	er := validate.Struct(user)
	if er != nil {
		msg = er.(validator.ValidationErrors).Error()
		return
	}
	msg = ""
	return
}



