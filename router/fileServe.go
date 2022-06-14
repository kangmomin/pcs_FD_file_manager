package router

import (
	"FD/util"
	"net/http"
	"os"

	"github.com/go-session/session/v3"
	"github.com/gorilla/mux"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(ctx, w, r)
	if err != nil {
		util.SessionErr(w)
		return
	}

	userId, ok := store.Get("userId")
	if !ok {
		util.LoginErr(w)
		return
	}

	log.Println(userId)
	path := "default/file/path" + mux.Vars(r)["path"]
	if _, err := os.Stat(path); os.IsNotExist(err) {
		util.GlobalErr("file dose not exist", err, 404, w)
		return
	}

	http.ServeFile(w, r, path)
}
