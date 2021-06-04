package controllers

import (
	"log"
	"net/http"
	"strconv"
	"books/models"
	"books/db"
	"books/response"
	"books/auth"
	"github.com/gorilla/mux"
)

func GetEntries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := db.DB.Ping(); err != nil {
		log.Fatal(err)
		response.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	// i := r.Header.Get("UserID")
	i := auth.GetAuthorizedID()
	id, _ := strconv.Atoi(i)

	params := mux.Vars(r)
	quarter := params["q"]
	year := params["y"]

	entries, err2 := db.DB.Query("SELECT Date, UserID, Content FROM entry WHERE (UserID=? AND QUARTER(Date)=? AND YEAR(Date)=?) ORDER BY Date ASC", id, quarter, year)

	if err2 != nil {
		log.Fatal(err2)
		response.RespondWithError(w, http.StatusBadRequest, err2.Error())
		return
	}
	var entry models.Entry
	var display []models.Entry
	for entries.Next() {
		err3 := entries.Scan(&entry.Date, &entry.UserID, &entry.Content)
		if err3 != nil {
			log.Fatal(err3)
			response.RespondWithError(w, http.StatusBadRequest, err3.Error())
			return
		}
		display = append(display, entry)
	}
	if len(display) == 0 {
		response.RespondWithJson(w, http.StatusOK, "No entries available for this date.")
		return
	}
	response.RespondWithJson(w, http.StatusOK, display)
}