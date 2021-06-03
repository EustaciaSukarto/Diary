package models

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

func ValidateToken(signedToken string) (id string, err error){
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims (signedToken, claims, func(token *jwt.Token) (interface{}, error) {return []byte("secret"), nil})
	if err != nil {
		return
	}
	if token == nil {
			return
	}	
	id = fmt.Sprintf("%v", claims["UserID"])
	return
}