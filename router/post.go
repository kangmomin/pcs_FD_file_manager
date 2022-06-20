package router

import (
	"FD/util"
	"database/sql"
	"strconv"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetPosts(w http.ResponseWriter, _ *http.Request) {
	var postList []util.PostList
	data, err := db.Query("SELECT post_id, u.user_name, title, created FROM post INNER JOIN \"user\" u ON u.user_id = post.user_id ") // 추후 page searching도 만들어야함.
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
	query := `
SELECT p.post_id, u.user_name, p.title, p.created FROM post p INNER JOIN "user" u ON u.user_id=p.user_id 
`

	var (
		searchSetting util.SearchBody
		queryParams   []any
		err           error
	)

	params := r.URL.Query()
	searchSetting.Club = params.Get("club")
	searchSetting.Word = params.Get("word")
	searchSetting.StartDate = params.Get("startDate")
	searchSetting.EndDate = params.Get("endDate")

	if len(searchSetting.Word) > 3 ||
		len(searchSetting.Club) > 0 ||
		len(searchSetting.StartDate) > 9 ||
		len(searchSetting.EndDate) > 9 {
		query += "WHERE "
	}
	dataIdx := 1

	// word search
	if len(searchSetting.Word) > 3 {
		query += "p.title ILIKE '%' || " + strconv.Itoa(dataIdx) + "|| '%' "
		queryParams = append(queryParams, searchSetting.Word)
		dataIdx++
	}

	// club search
	if len(searchSetting.Club) > 0 {
		if dataIdx > 1 {
			query += "AND "
		}
		query += "p.club=" + strconv.Itoa(dataIdx)
		searchSetting.ClubId, err = strconv.Atoi(searchSetting.Club)
		if err != nil {
			util.GlobalErr("Club id is not number", nil, 400, w)
			return
		}
		queryParams = append(queryParams, searchSetting.ClubId)
		dataIdx++
	}

	// time search
	if len(searchSetting.StartDate) > 9 {
		if dataIdx > 1 {
			query += " AND "
		}
		// 2022-03-23
		query += "post_id=(SELECT post_id FROM post WHERE created BETWEEN " + strconv.Itoa(dataIdx) + " AND " + strconv.Itoa(dataIdx+1) + ")"
		queryParams = append(queryParams, searchSetting.StartDate, searchSetting.EndDate)
	}
	query += " ORDER BY post_id DESC LIMIT 30;"

	data, err := db.Query(query, queryParams...)
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
