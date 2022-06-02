package router

import (
	"fmt"
	"net/http"

	_ "github.com/lib/pg"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
}
