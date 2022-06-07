package router

import (
	"FD/logger"
	"FD/util"

	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

var db *sql.DB = util.DB
var log = logger.Logger

func GetPosts(w http.ResponseWriter, _ *http.Request) {
	var postList []util.PostList
	data, err := db.Query("SELECT post_id, user_name, title, created FROM post ORDER BY post_id DESC LIMIT 30;") // 추후 page searching도 만들어야함.
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		resData, _ := json.Marshal(util.Res{
			Data: nil,
			Err:  true,
		})
		fmt.Fprint(w, string(resData))
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
	sql := "SELECT post_id, user_name, title, created FROM post WHERE "

	var searchSetting util.SearchBody
	err := json.NewDecoder(r.Body).Decode(&searchSetting)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		resData, _ := json.Marshal(util.Res{
			Data: nil,
			Err:  true,
		})
		fmt.Fprint(w, string(resData))
		return
	}

	// word search
	if len(searchSetting.Word) < 3 {
		sql += "title LIKE %" + searchSetting.Word + "% AND "
	}

	// club search
	if len(searchSetting.Club) > 0 {
		sql += "club=" + searchSetting.Club + " AND "
	}

	// time search
	if len(searchSetting.StartDate) > 9 {
		// 2022-03-23
		sql += "post_id=(SELECT post_id WHERE create BETWEEN" + searchSetting.StartDate
		if len(searchSetting.EndDate) > 9 {
			sql += "AND " + searchSetting.EndDate
		}
		sql += ")"
	}

	sql += "ORDER BY post_id LIMIT 30;"

	data, err := db.Query(sql)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		resData, _ := json.Marshal(util.Res{
			Data: nil,
			Err:  true,
		})
		fmt.Fprint(w, string(resData))
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
	var postId struct {
		id int `json:"post_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&postId.id)
	if err != nil {
		resData, _ := json.Marshal(util.Res{
			Data: nil,
			Err:  true,
		})
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string(resData))
		return
	}

	var postDetail util.PostDetail
	err = db.QueryRow("SELECT club_id, title, readme, file_path, created FROMM post WHERE id=?;", postId.id).
		Scan(&postDetail.ClubId, &postDetail.Title, &postDetail.FilePath, &postDetail.Created)
	if err != nil {
		log.Println(err)
		resData, _ := json.Marshal(util.Res{
			Data: nil,
			Err:  true,
		})
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string(resData))
		return
	}

	err = db.QueryRow("SELECT user_name FROM user WHERE user_id=(SELECT writer_id FROM post WHERE post_id=?);", postId.id).
		Scan(&postDetail.WriterName)
	if err != nil {
		log.Println(err)
		resData, _ := json.Marshal(util.Res{
			Data: nil,
			Err:  true,
		})
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string(resData))
		return
	}

	resData, _ := json.Marshal(util.Res{
		Data: postDetail,
		Err:  false,
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(resData))
}
