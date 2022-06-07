package util

import "time"

type PostList struct {
	PostId   int       `json:"post_id"`
	UserName string    `json:"user_name"`
	Title    string    `json:"title"`
	Created  time.Time `json:"created"`
}

type SearchBody struct {
	Word      string `json:"word"`
	Club      string `json:"club"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}
