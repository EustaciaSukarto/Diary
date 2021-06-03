package controllers

import (
	"net/http"
	"encoding/json"
	"strconv"
	"books/models"
	"books/db"
	"books/response"
)

func CreateEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := db.DB.Ping(); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	i := r.Header.Get("UserID")
	id, err := strconv.Atoi(i)
	var entry models.Entry
	_ = json.NewDecoder(r.Body).Decode(&entry)

	ok := models.ValidateEntry(entry)
	if ok != "" {
		response.RespondWithError(w, http.StatusBadRequest, ok)
		return
	}

	var content string
	err1 := db.DB.QueryRow("SELECT Content FROM entry WHERE (UserID=? AND Date=?)", id, entry.Date).Scan(&content)
	if err1 != nil {
		response.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if content == "" {
		new, err2 := db.DB.Prepare("INSERT INTO entry (Date, UserID, Content) VALUES (?,?,?)")
		if err2 != nil {
			response.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		_, err3 := new.Exec(entry.Date, entry.UserID, entry.Content)
		if err3 != nil {
			response.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		response.RespondWithJson(w, http.StatusOK, new)
	} else {
		update, err2 := db.DB.Prepare("UPDATE entry SET Content=? WHERE (Date=? AND UserID=?)")
		if err2 != nil {
			response.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		_, err3 := update.Exec(entry.Content, entry.Date, entry.UserID)
		if err3 != nil {
			response.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		response.RespondWithJson(w, http.StatusOK, update)
	}
}