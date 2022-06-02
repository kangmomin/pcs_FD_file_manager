package router

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
	godotenv.Load("./.env")
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		log.Println(err)
	}
	defer db.Close()

}
