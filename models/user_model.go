package models

import (
	"log"
	"github.com/go-playground/validator/v10"
	"books/db"
	// "unicode"
)

type User struct {
	ID     		int  	`json:"id" validate:"isdefault"`
	Fullname   	string 	`json:"fullname"`
	Birthday  	string 	`json:"birthday"`
	Email 		string 	`json:"email" validate:"required,email"`
	Username 	string 	`json:"username" validate:"required"`
	// Password 	string 	`json:"password" validate:"password-rules"`
	Password 	string 	`json:"password"`
}

type UserInput struct {
	Fullname   	string 	`json:"fullname"`
	Birthday  	string 	`json:"birthday"`
	Email 		string 	`json:"email" validate:"required,email"`
	Username 	string 	`json:"username" validate:"required"`
	// Password 	string 	`json:"password" validate:"password-rules"`
	Password 	string 	`json:"password"`
}

// func PasswordValidator(v validator.FieldLevel) bool {
func PasswordValidator(pw string) string {
	//		has: 	upper  lower  number char
	rules := []bool{false, false, false, false}
	if (len(pw) < 6 || len(pw) > 32) {
		return "Password must be 6 - 32 characters."
	}
	for _, character := range pw {
		if (character >= 65 && character <= 90) {
			log.Print(character)
			rules[0] = true
		}
		if (character >= 97 && character <= 122) {
			log.Print(character)
			rules[1] = true
		}
		if (character >= 48 && character <= 57) {
			log.Print(character)
			rules[2] = true
		}
		if ((character >= 32 && character <= 47) || (character >= 58 && character <= 64) || (character >= 91 && character <= 96) || (character >= 123 && character <= 126)) {
			log.Print(character)
			rules[3] = true
		}
	}
	if (rules[0] == true && rules[1] == true && rules[2] == true && rules[3] == true) {
		return ""
	}
	return "Password must contain at least 1 uppercase letter, 1 lowercase letter, 1 number, and 1 special character."
}

func QueryUser(email string) (user User, err error) {
	err = db.DB.QueryRow("SELECT ID, Fullname, Birthday, Email, Username, Password FROM user WHERE Email=?", email).Scan(&user.ID, &user.Fullname, &user.Birthday, &user.Email, &user.Username, &user.Password)
	// log.Print(user.ID)
	// log.Print(user.Email)
	return
}

func ValidateUserInput(user UserInput) (msg string) {
	// check if user already exists
	dbuser, err := QueryUser(user.Email)
	if dbuser.Email == "" {
		err = nil
	}
	if err != nil {
		return err.Error()
	}
	if dbuser.Email != "" {
		msg = "Email already in use."
		return
	}
	// validate input rules
	validate := validator.New()
	// validate.RegisterValidation("password-rules", PasswordValidator)
	er := validate.Struct(user)
	if er != nil {
		msg = er.(validator.ValidationErrors).Error()
		return
	}
	msg = PasswordValidator(user.Password)
	return
}



