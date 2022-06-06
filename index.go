package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"FD/logger"
	"FD/router"
)

func main() {
	logger.LogSetting()
	app := mux.NewRouter()

	// post
	app.HandleFunc("/post", router.GetPosts).Methods("GET")

	http.ListenAndServe(":3000", app)
}
