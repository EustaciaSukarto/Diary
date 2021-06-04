package auth

import (
	"log"
	"net/http"
	"books/models"
	"books/response"
)

var AuthorizedID string

func Authenticate(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		clientToken := r.Header.Get("Authorization")
		if clientToken == "" {
			response.RespondWithError(w, http.StatusBadRequest, "No Authorization header provided.")
			return
		}
		id, err := models.ValidateToken(clientToken)
		if err != nil {
			log.Fatal(err)
			response.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		SetAuthorizedID(id)
		// w.Header().Set("UserID", id)
		handler.ServeHTTP(w, r)
	})
}

func SetAuthorizedID(id string) {
	AuthorizedID = id
}

func GetAuthorizedID() (string) {
	return AuthorizedID
}