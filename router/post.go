package router

import (
	"FD/util"
	"database/sql"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetPosts(w http.ResponseWriter, _ *http.Request) {
	var postList []util.PostList
	data, err := db.Query("SELECT post_id, u.user_name, title, created FROM post INNER JOIN \"user\" u ON u.user_id = post.user_id ORDER BY post_id DESC LIMIT 30;") // 추후 page searching도 만들어야함.
	if err != nil {
		if err == sql.ErrNoRows {
			util.GlobalErr("no data", nil, 404, w)
		} else {
			util.GlobalErr("select error", err, 400, w)
		}
		return
	}

	for data.Next() {
		var post util.PostList
		data.Scan(&post.PostId, &post.UserName, &post.Title, &post.Created)

		postList = append(postList, post)
	}

	resData, _ := json.Marshal(util.Res{
		Data: postList,
		Err:  false,
	})

	fmt.Fprint(w, string(resData))
}

func SearchPost(w http.ResponseWriter, r *http.Request) {
	query := "SELECT post_id, user_name, title, created FROM post WHERE "

	var searchSetting util.SearchBody
	err := json.NewDecoder(r.Body).Decode(&searchSetting)
	if err != nil {
		util.GlobalErr("body data wrong", err, 400, w)
		return
	}

	// word search
	if len(searchSetting.Word) < 3 {
		query += "title LIKE %" + searchSetting.Word + "% AND "
	}

	// club search
	if len(searchSetting.Club) > 0 {
		query += "club=" + searchSetting.Club + " AND "
	}

	// time search
	if len(searchSetting.StartDate) > 9 {
		// 2022-03-23
		query += "post_id=(SELECT post_id WHERE create BETWEEN" + searchSetting.StartDate
		if len(searchSetting.EndDate) > 9 {
			query += "AND " + searchSetting.EndDate
		}
		query += ")"
	}

	query += "ORDER BY post_id LIMIT 30;"

	data, err := db.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			util.GlobalErr("no data", nil, 404, w)
		} else {
			util.GlobalErr("select error", err, 400, w)
		}
		return
	}

	var postList []util.PostList
	for data.Next() {
		var row util.PostList
		data.Scan(&row.PostId, &row.UserName, &row.Title, &row.Created)
		postList = append(postList, row)
	}

	resData, _ := json.Marshal(util.Res{
		Data: postList,
		Err:  false,
	})
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(resData))
}

func PostDetail(w http.ResponseWriter, r *http.Request) {
	postId, ok := mux.Vars(r)["postId"]
	if !ok {
		util.GlobalErr("not enough params", nil, 500, w)
		return
	}

	var postDetail util.PostDetail
	err := db.QueryRow("SELECT u.user_name, p.club_id, p.title, p.readme, p.file_path, p.created FROM post p INNER JOIN \"user\" u ON u.user_id=p.user_id WHERE post_id=$1;", postId).
		Scan(&postDetail.WriterName, &postDetail.ClubId, &postDetail.Title, &postDetail.Readme, &postDetail.FilePath, &postDetail.Created)
	if err != nil {
		if err == sql.ErrNoRows {
			util.GlobalErr("no data", nil, 404, w)
		} else {
			util.GlobalErr("select error", err, http.StatusBadRequest, w)
		}
		return
	}

	resData, _ := json.Marshal(util.Res{
		Data: postDetail,
		Err:  false,
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(resData))
}
