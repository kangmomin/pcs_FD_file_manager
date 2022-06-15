package router

import (
	"FD/util"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func WritePost(w http.ResponseWriter, r *http.Request) {
	var userId interface{}

	if userId = util.LoginCheck(w, r); userId == nil {
		util.LoginErr(w)
		return
	}

	// post data
	var pd util.WritePost

	err := json.NewDecoder(r.Body).Decode(&pd)
	if err != nil {
		util.GlobalErr("body data wrong", err, 400, w)
		return
	}

	postId, err := db.Exec(`INSERT INTO public.post(
		title, readme, file_path, created, user_id, club_id) VALUES ($1, $2, $3, $4, $5, $6);`,
		pd.Title, pd.Readme, pd.FilePath, pd.Created, userId, pd.ClubId)

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

func DeletePost(w http.ResponseWriter, r *http.Request) {
	var userId interface{}

	if userId = util.LoginCheck(w, r); userId == nil {
		util.LoginErr(w)
		return
	}
	postId := mux.Vars(r)["postId"]
	if numPostId, err := strconv.Atoi(postId); err != nil || numPostId < 1 {
		util.GlobalErr("wrong post_id", nil, 400, w)
		return
	}
	_, err := db.Exec(`DELETE from post WHERE user_id='$1' AND post_id=$2`, userId, postId)
	if err != nil {
		util.GlobalErr("cannot delete post", err, 500, w)
		return
	}

	resData, _ := json.Marshal(util.Res{
		Data: "delete success",
		Err:  false,
	})
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(resData))
}
