package models

import (
	// "log"
	"github.com/go-playground/validator/v10"
	"books/db"
)

type Entry struct {	
	ID 		int 	`"json:id" validate:"isdefault`
	Date 	string 	`"json:"date" validate:"required"`
	UserID 	int 	`"json:userid" validate:"isdefault"`
	Content string 	`"json:content"`
}

type EntryInput struct {
	Date string `"json:"date" validate:"required"`
	Content string 	`"json:content"`

}

func QueryEntry(userid int, date string) (entry Entry, err error) {
	err = db.DB.QueryRow("SELECT ID, Date, UserID, Content FROM entry WHERE (UserID=? AND Date=?)", userid, date).Scan(&entry.ID, &entry.Date, &entry.UserID, &entry.Content)
	return
}

func ValidateEntryInput(entry EntryInput) (msg string) {
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