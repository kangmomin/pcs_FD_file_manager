package util

import "time"

type PostList struct {
	PostId   int
	UserName string
	Title    string
	Created  time.Time
}

type SearchBody struct {
	Word      string `json:"word"`
	Club      string `json:"club"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}
