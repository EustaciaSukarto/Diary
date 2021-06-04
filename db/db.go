package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:123abcDEF@tcp(127.0.0.1:3306)/diary")
	return db, err
}

func Init() (error) {
	db, err := InitDB()
	DB = db
	return err
}

func CloseAll() {
	if DB != nil {
		DB.Close()
	}
}

// func GetDB() {
// 	db, err := sql.Open("mysql", "root:1234567890@tcp(127.0.0.1:3306)/diary")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	DB = db
// }

