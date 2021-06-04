package controllers

import (
	"log"
	"net/http"
	"encoding/json"
	"books/models"
	"books/db"
	"books/response"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
)

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := db.DB.Ping(); err != nil {
		log.Fatal(err)
		response.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	var details models.LoginDetails
	_ = json.NewDecoder(r.Body).Decode(&details)
	ok := models.ValidateLoginDetails(details)
	if ok != "" {
		response.RespondWithError(w, http.StatusBadRequest, ok)
		return
	}
	var user models.User
	err1 := db.DB.QueryRow("SELECT * FROM user WHERE (Email=? OR Username=?)", details.Detail1, details.Detail1).Scan(&user.ID, &user.Fullname, &user.Birthday, &user.Email, &user.Username, &user.Password)
	if err1 != nil {
		log.Fatal(err1)
		response.RespondWithError(w, http.StatusBadRequest, err1.Error())
		return
	}
	if user.Password == "" {
		response.RespondWithError(w, http.StatusBadRequest, "User not reqistered.")
		return
	}
	check := CheckPasswordHash(details.Password, user.Password)
	if !check {
		response.RespondWithError(w, http.StatusBadRequest, "Password is incorrect.")
		return
	}
	claims := jwt.MapClaims{}
	claims["UserID"] = user.ID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err2 := token.SignedString([]byte("secret"))
	if err2 != nil {
		log.Fatal(err2)
		response.RespondWithError(w, http.StatusBadRequest, err2.Error())
		return
	}
	response.RespondWithJson(w, http.StatusOK, "AuthToken:" + signedToken)
}