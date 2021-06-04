package controllers

import (
	"log"
	"net/http"
	"encoding/json"
	"strconv"
	"books/models"
	"books/db"
	"books/response"
	"books/auth"
)

func CreateEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := db.DB.Ping(); err != nil {
		log.Fatal(err)
		response.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	// i := r.Header.Get("UserID")
	i := auth.GetAuthorizedID()
	id, errr := strconv.Atoi(i)
	if errr != nil {
		log.Fatal(errr)
		response.RespondWithError(w, http.StatusBadRequest, errr.Error())
		return
	}
	var entry models.EntryInput
	_ = json.NewDecoder(r.Body).Decode(&entry)

	ok := models.ValidateEntryInput(entry)
	if ok != "" {
		response.RespondWithError(w, http.StatusBadRequest, ok)
		return
	}


	ent, err1 := models.QueryEntry(id, entry.Date)
	if err1.Error() == "sql: no rows in result set" {
		err1 = nil
	}
	if err1 != nil {
		log.Fatal(err1)
		response.RespondWithError(w, http.StatusBadRequest, err1.Error())
		return
	}
	if ent.Content == "" {
		new, err2 := db.DB.Prepare("INSERT INTO entry (Date, UserID, Content) VALUES (?,?,?)")
		if err2 != nil {
			log.Fatal(err2)
			response.RespondWithError(w, http.StatusBadRequest, err2.Error())
			return
		}
		_, err3 := new.Exec(entry.Date, id, entry.Content)
		if err3 != nil {
			log.Fatal(err3)
			response.RespondWithError(w, http.StatusBadRequest, err3.Error())
			return
		}
		display, err4 := models.QueryEntry(id, entry.Date)
		if err4 != nil {
			log.Fatal(err3)
			response.RespondWithError(w, http.StatusBadRequest, err4.Error())
			return
		}
		response.RespondWithJson(w, http.StatusOK, display)
	} else {
		update, err2 := db.DB.Prepare("UPDATE entry SET Content=? WHERE (Date=? AND UserID=?)")
		if err2 != nil {
			log.Fatal(err2)
			response.RespondWithError(w, http.StatusBadRequest, err2.Error())
			return
		}
		_, err3 := update.Exec(entry.Content, entry.Date, id)
		if err3 != nil {
			log.Fatal(err3)
			response.RespondWithError(w, http.StatusBadRequest, err3.Error())
			return
		}
		display, err4 := models.QueryEntry(id, entry.Date)
		if err4 != nil {
			log.Fatal(err3)
			response.RespondWithError(w, http.StatusBadRequest, err4.Error())
			return
		}
		response.RespondWithJson(w, http.StatusOK, display)
	}
}