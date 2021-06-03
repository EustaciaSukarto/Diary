package auth

import (
	"net/http"
	"books/models"
	"books/response"
)

func Authenticate(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		clientToken := r.Header.Get("Authorization")
		if clientToken == "" {
			response.RespondWithError(w, http.StatusBadRequest, "No Authorization header provided.")
			return
		}
		id, err := models.ValidateToken(clientToken)
		if err != nil {
			response.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		r.Header.Set("UserID", id)
		handler.ServeHTTP(w, r)
	})
}