package main

import (
	"github.com/gorilla/mux"

	"FD/router"
)

func main() {
	app := mux.NewRouter()

	// post
	app.HandleFunc("/post", router.GetPosts).Methods("GET")
}
