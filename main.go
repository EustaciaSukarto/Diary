package main

import (
	"log"
	"net/http"
	
	"github.com/gorilla/mux"
	
	"books/db"
	"books/controllers"
	"books/auth"
)

func main() {
	err := db.Init()
	if err != nil {
		return
	}
	defer db.CloseAll()
	// Init Router
	r := mux.NewRouter()
	// Routes
	r.HandleFunc("/api/register", controllers.CreateUser).Methods("POST")
	r.HandleFunc("/api/login", controllers.LoginUser).Methods("POST")
	r.HandleFunc("/api/user/entries/:q/:y", auth.Authenticate(http.HandlerFunc(controllers.GetEntries))).Methods("GET")
	r.HandleFunc("/api/user/newentry", auth.Authenticate(http.HandlerFunc(controllers.CreateEntry))).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", r))
}
