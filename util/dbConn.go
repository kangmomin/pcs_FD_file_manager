package util

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB = connDB()

func connDB() *sql.DB {
	godotenv.Load("./.env")
	dbinfo := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PWD"), os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		log.Println(err)
	}

	return db
}
