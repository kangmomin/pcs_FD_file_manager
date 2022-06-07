package util

import "time"

type PostList struct {
	PostId   int       `json:"post_id"`
	UserName string    `json:"user_name"`
	Title    string    `json:"title"`
	Created  time.Time `json:"created"`
}

type PostDetail struct {
	WriterName string    `json:"user_name"`
	ClubId     int       `json:"club_id"`
	Title      string    `json:"title"`
	Readme     string    `json:"readme"`
	FilePath   []string  `json:"file_path"`
	Created    time.Time `json:"created"`
}

type SearchBody struct {
	Word      string `json:"word"`
	Club      string `json:"club"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}
