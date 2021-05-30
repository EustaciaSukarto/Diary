package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	// "github.com/go-sql-driver/mysql"
	"encoding/json"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"github.com/go-playground/validator/v10"
)

func GetDB() (db *sql.DB, err error, w http.ResponseWriter) {
	db, err = sql.Open("mysql", "root:1234567890@tcp(127.0.0.1:3306)/diary")
	
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	return
}

type User struct {
	ID     		int  	`json:"id" validate:"isdefault"`
	Fullname   	string 	`json:"fullname"`
	Birthday  	string 	`json:"birthday"`
	Email 		string 	`json:"email" validate:"required,email"` // fix validation later
	Username 	string 	`json:"username" validate: "required"` // fix validation later
	Password 	string 	`json:"password" validate: "required,min=6,max=32"` // fix validation later
}

type Entry struct {
	Date 	string 	`"json:"date" validate:"required"`
	UserID 	int 	`"json:userid" validate:"isdefault"`
	Content string 	`"json:content" validate:"required"`
}

type LoginDetails struct {
	Detail1 string `json:"email/username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func QueryUser(email string) User {
	db, err, w := GetDB()
	defer db.Close()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	var user User
	err2 := db.QueryRow("SELECT * FROM user WHERE Email=?", user.Email).Scan(&user.ID, &user.Fullname, &user.Birthday, &user.Email, &user.Username, &user.Password)
	if err2 != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	return user
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	
	db, err, w := GetDB()
	defer db.Close()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	validate := validator.New()
	er := validate.Struct(user)
	if er != nil {
		respondWithError(w, http.StatusBadRequest, er.(validator.ValidationErrors).Error())
		return
	}
	// check if email is already used
	dbuser := QueryUser(user.Email)
	if dbuser.Email != "" {
		respondWithError(w, http.StatusBadRequest, "Email already in use.")
		return
	}

	hashedPassword, err1 := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err1 != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(hashedPassword) != 0 {
		insert, err2 := db.Prepare("INSERT INTO user (Fullname, Birthday, Email, Username, Password) VALUES (?,?,?,?,?)")
		if err2 != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		_, err3 := insert.Exec(user.Fullname, user.Birthday, user.Email, user.Username, hashedPassword)
		if err3 != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	} else {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	display := QueryUser(user.Email)
	respondWithJson(w, http.StatusOK, display)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err, w := GetDB()
	defer db.Close()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	var details LoginDetails
	_ = json.NewDecoder(r.Body).Decode(&details)
	var user User
	validate := validator.New()
	er := validate.Struct(details)
	if er != nil {
		respondWithError(w, http.StatusBadRequest, er.(validator.ValidationErrors).Error())
		return
	}
	err1 := db.QueryRow("SELECT * FROM user WHERE (Email=? OR Username=?)", details.Detail1, details.Detail1).Scan(&user.ID, &user.Fullname, &user.Birthday, &user.Email, &user.Username, &user.Password)
	if err1 != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if user.Password == "" {
		respondWithError(w, http.StatusBadRequest, "User not reqistered.")
		return
	}
	check := CheckPasswordHash(details.Password, user.Password)
	if !check {
		respondWithError(w, http.StatusBadRequest, "Password is incorrect.")
		return
	}
	// generate Jwt token
	// user = QueryUser(user.Email)
	claims := jwt.MapClaims{}
	claims["UserID"] = user.ID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err2 := token.SignedString([]byte("secret"))
	if err2 != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, signedToken)
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Auth(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		clientToken := r.Header.Get("Authorization")
		if clientToken == "" {
			respondWithError(w, http.StatusBadRequest, "No Authorization header provided.")
			return
		}
		id, err := ValidateToken(clientToken)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		r.Header.Set("UserID", id)
		handler.ServeHTTP(w, r)
	})
}

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

func GetEntries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err, w := GetDB()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer db.Close()
	i := r.Header.Get("UserID")
	id, err := strconv.Atoi(i)
	quarter := r.URL.Query().Get("q")
	year := r.URL.Query().Get("y")

	entries, err2 := db.Query("SELECT (Date, UserID, Content) FROM entry WHERE (UserID=? AND QUARTER(Date) AND YEAR(Date)) ORDER BY Date ASC", id, quarter, year)
	if err2 != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	entry := Entry{}
	display := []Entry{}
	for entries.Next() {
		err3 := entries.Scan(&entry.Date, &entry.UserID, &entry.Content)
		if err3 != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		display = append(display, entry)
	}
	respondWithJson(w, http.StatusOK, display)
}

func CreateEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err, w := GetDB()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer db.Close()
	i := r.Header.Get("UserID")
	id, err := strconv.Atoi(i)
	var entry Entry
	_ = json.NewDecoder(r.Body).Decode(&entry)
	validate := validator.New()
	er := validate.Struct(entry)
	if er != nil {
		respondWithError(w, http.StatusBadRequest, er.(validator.ValidationErrors).Error())
		return
	}
	var content string
	err1 := db.QueryRow("SELECT Content FROM entry WHERE (UserID=? AND Date=?)", id, entry.Date).Scan(&content)
	if err1 != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if content == "" {
		new, err2 := db.Prepare("INSERT INTO entry (Date, UserID, Content) VALUES (?,?,?)")
		if err2 != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		_, err3 := new.Exec(entry.Date, entry.UserID, entry.Content)
		if err3 != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithJson(w, http.StatusOK, new)
	} else {
		update, err2 := db.Prepare("UPDATE entry SET Content=? WHERE (Date=? AND UserID=?)")
		if err2 != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		_, err3 := update.Exec(entry.Content, entry.Date, entry.UserID)
		if err3 != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithJson(w, http.StatusOK, update)
	}
}

func main() {
	// Init Router
	r := mux.NewRouter()
	// Routes
	r.HandleFunc("/api/register", CreateUser).Methods("POST")
	r.HandleFunc("/api/login", LoginUser).Methods("POST")
	r.HandleFunc("/api/user/entries/:q/:y", Auth(http.HandlerFunc(GetEntries))).Methods("GET")
	r.HandleFunc("/api/user/newentry", Auth(http.HandlerFunc(CreateEntry))).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", r))
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}
func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
