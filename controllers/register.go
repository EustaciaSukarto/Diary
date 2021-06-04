package controllers

import (
	"log"
	"net/http"
	"encoding/json"
	"books/models"
	"books/db"
	"books/response"
	"golang.org/x/crypto/bcrypt"
	// "github.com/dgrijalva/jwt-go"
)


func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.UserInput
	_ = json.NewDecoder(r.Body).Decode(&user)
	
	if err := db.DB.Ping(); err != nil {
		log.Fatal(err)
		response.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	ok := models.ValidateUserInput(user)
	if ok != "" {
		response.RespondWithError(w, http.StatusBadRequest, ok)
		return
	}
	
	hashedPassword, err1 := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err1 != nil {
		log.Fatal(err1)
		response.RespondWithError(w, http.StatusBadRequest, err1.Error())
		return
	}

	if len(hashedPassword) != 0 {
		insert, err2 := db.DB.Prepare("INSERT INTO user (Fullname, Birthday, Email, Username, Password) VALUES (?,?,?,?,?)")
		if err2 != nil {
			log.Fatal(err2)
			response.RespondWithError(w, http.StatusBadRequest, err2.Error())
			return
		}
		_, err3 := insert.Exec(user.Fullname, user.Birthday, user.Email, user.Username, hashedPassword)
		if err3 != nil {
			log.Fatal(err3)
			response.RespondWithError(w, http.StatusBadRequest, err3.Error())
			return
		}
	} else {
		response.RespondWithError(w, http.StatusBadRequest, "Error in hashing password.")
		return
	}
	display, err3 := models.QueryUser(user.Email)
	if err3 != nil {
			response.RespondWithError(w, http.StatusBadRequest, err3.Error())
			return
	}
	// log.Print(display.ID)
	// log.Print(display.Email)
	response.RespondWithJson(w, http.StatusOK, display)
}