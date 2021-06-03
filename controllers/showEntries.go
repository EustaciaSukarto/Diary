package controllers

import (
	"net/http"
	"strconv"
	"books/models"
	"books/db"
	"books/response"
)

func GetEntries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := db.DB.Ping(); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	i := r.Header.Get("UserID")
	id, err := strconv.Atoi(i)
	quarter := r.URL.Query().Get("q")
	year := r.URL.Query().Get("y")

	entries, err2 := db.DB.Query("SELECT (Date, UserID, Content) FROM entry WHERE (UserID=? AND QUARTER(Date)=? AND YEAR(Date)=?) ORDER BY Date ASC", id, quarter, year)

	if err2 != nil {
		response.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	var entry models.Entry
	var display []models.Entry
	for entries.Next() {
		err3 := entries.Scan(&entry.Date, &entry.UserID, &entry.Content)
		if err3 != nil {
			response.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		display = append(display, entry)
	}
	response.RespondWithJson(w, http.StatusOK, display)
}