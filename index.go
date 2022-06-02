package main

import (
	"github.com/gorilla/mux"

	"FD/logger"
	"FD/router"
)

func main() {
	logger.LogSetting()
	app := mux.NewRouter()

	// post
	app.HandleFunc("/post", router.GetPosts).Methods("GET")
}
