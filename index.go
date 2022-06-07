package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"FD/router"
)

func main() {
	app := mux.NewRouter()

	// post
	app.HandleFunc("/post", router.GetPosts).Methods("GET")

	http.ListenAndServe(":8080", app)
}
