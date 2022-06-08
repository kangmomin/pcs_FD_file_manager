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
	app.HandleFunc("/post/{postId}", router.PostDetail).Methods("GET")
	app.HandleFunc("/post/search", router.SearchPost).Methods("GET")

	// account
	app.HandleFunc("/login", router.Login).Methods("POST")
	app.HandleFunc("/sign-up", router.SignUp).Methods("POST")

	http.ListenAndServe(":8080", app)
}
