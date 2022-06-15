package router

import (
	"FD/util"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-session/session/v3"
)

func WritePost(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(ctx, w, r)
	if err != nil {
		util.LoginErr(w)
		return
	}

	data, ok := store.Get("userId")
	if !ok {
		util.LoginErr(w)
		return
	}

	// post data
	var pd util.WritePost

	err = json.NewDecoder(r.Body).Decode(&pd)
	if err != nil {
		util.GlobalErr("body data wrong", err, 400, w)
		return
	}

	postId, err := db.Exec(`INSERT INTO public.post(
		title, readme, file_path, created, user_id, club_id) VALUES ($1, $2, $3, $4, $5, $6);`,
		pd.Title, pd.Readme, pd.FilePath, pd.Created, data, pd.ClubId)

	if err != nil {
		util.GlobalErr("inserting err", err, 500, w)
		return
	}
	resData, _ := json.Marshal(util.Res{
		Data: postId,
		Err:  false,
	})
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(resData))
}
