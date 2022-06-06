package router

import (
	"FD/util"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var db *sql.DB = util.DB

type res struct {
	Data any
	Err  bool `json:"error"`
}

func GetPosts(w http.ResponseWriter, _ *http.Request) {
	var postList []util.PostList
	data, err := db.Query("SELECT post_id, user_name, title, created FROM post ORDER BY post_id DESC LIMIT 30;") // 추후 page searching도 만들어야함.
	if err != nil {
		log.Println(err)
		resData, _ := json.Marshal(res{
			Data: nil,
			Err:  true,
		})
		fmt.Fprint(w, string(resData))
	}

	for data.Next() {
		var post util.PostList
		data.Scan(&post.PostId, &post.UserName, &post.Title, &post.Created)

		postList = append(postList, post)
	}

	resData, _ := json.Marshal(res{
		Data: postList,
		Err:  false,
	})

	fmt.Fprint(w, string(resData))
}
