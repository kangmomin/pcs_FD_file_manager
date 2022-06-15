package main

import (
	"net/http"

	"github.com/go-session/redis/v3"
	"github.com/go-session/session/v3"
	"github.com/gorilla/mux"

	"FD/router"
)

func main() {
	session.InitManager(
		session.SetStore(redis.NewRedisStore(&redis.Options{
			Addr: "localhost:6379",
			DB:   15,
		})),
	)

	app := mux.NewRouter()

	// post
	app.HandleFunc("/post", router.GetPosts).Methods("GET")
	app.HandleFunc("/post/{postId}", router.PostDetail).Methods("GET")
	app.HandleFunc("/post/search", router.SearchPost).Methods("GET")

	// admin rights
	app.HandleFunc("/admin/post", router.WritePost).Methods("POST")
	app.HandleFunc("/admin/post/{postId}", router.DeletePost).Methods("DELETE")

	// account
	app.HandleFunc("/login", router.Login).Methods("POST")
	app.HandleFunc("/sign-up", router.SignUp).Methods("POST")

	// file download
	app.HandleFunc("/file/{path}", router.DownloadFile).Methods("POST")

	http.ListenAndServe(":8080", app)
}
