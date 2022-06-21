package router

import (
	"FD/util"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func WritePost(w http.ResponseWriter, r *http.Request) {
	var (
		userId interface{} // test
		ok     bool
	)

	if userId, ok = util.AdminCheck(r); !ok {
		util.GlobalErr("not admin", nil, http.StatusForbidden, w)
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
	var (
		userId interface{} // test
		ok     bool
	)

	if userId, ok = util.AdminCheck(r); !ok {
		util.GlobalErr("not admin", nil, http.StatusForbidden, w)
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

func UserList(w http.ResponseWriter, r *http.Request) {
	var (
		userId interface{}
		ok     bool
	)

	if userId, ok = util.AdminCheck(r); !ok {
		util.GlobalErr("not admin", nil, http.StatusForbidden, w)
		return
	}

	var userList []util.UserList

	data, err := db.Query(`SELECT user_id, club_id, user_name, email, phone_num FROM user;`)
	if err != nil {
		if err == sql.ErrNoRows {
			util.GlobalErr("not found", nil, 404, w)
		} else {
			util.GlobalErr("SELECT error", err, 500, w)
		}
		return
	}

	for data.Next() {
		var user util.UserList
		if err := data.Scan(&user.UserId, &user.ClubId, &user.UserName, &user.Email, &user.PhoneNum); err != nil {
			continue
		}

		userList = append(userList, user)
	}

	resData, _ := json.Marshal(util.Res{
		Data: userList,
		Err:  false,
	})
	adminLog.Println("[" + fmt.Sprintf("%v", userId) + "] inquire user information")
	w.WriteHeader(200)
	fmt.Fprint(w, string(resData))
}

func ApplyAdminList(w http.ResponseWriter, r *http.Request) {
	var (
		userId interface{}
		ok     bool
	)

	if userId, ok = util.AdminCheck(r); !ok {
		util.GlobalErr("not admin", nil, http.StatusForbidden, w)
		return
	}

	var applyList []util.ApplyAdmin

	data, err := db.Query(`SELECT user_id, user_name, club_id FROM admin a INNER JOIN user u ON u.user_id=a.user_id WHERE a.accept=false;`)
	if err != nil {
		if err == sql.ErrNoRows {
			util.GlobalErr("data not found", nil, 404, w)
		} else {
			util.GlobalErr("SELECT error", err, 500, w)
		}
		return
	}

	for data.Next() {
		var user util.ApplyAdmin
		data.Scan(&user.UserId, &user.UserName, &user.ClubId)
		applyList = append(applyList, user)
	}

	resData, _ := json.Marshal(util.Res{
		Data: applyList,
		Err:  false,
	})
	adminLog.Println("[" + fmt.Sprintf("%v", userId) + "] was inquire user list that apply for admin")
	w.WriteHeader(200)
	fmt.Fprint(w, string(resData))
}

func AcceptUser(w http.ResponseWriter, r *http.Request) {
	var ok bool

	if _, ok = util.AdminCheck(r); !ok {
		util.GlobalErr("not admin", nil, http.StatusForbidden, w)
		return
	}

	userId := mux.Vars(r)["userId"]

	_, err := db.Exec(`UPDATE TO admin SET accept=true WHERE user_id=$1`, userId)
	if err != nil {
		util.GlobalErr("cannot update", err, 500, w)
		return
	}

	resData, _ := json.Marshal(util.Res{
		Data: "update success",
		Err:  false,
	})
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(resData))
}
