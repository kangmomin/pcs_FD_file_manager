package util

import "time"

type PostList struct {
	PostId   int
	UserName string
	Title    string
	Created  time.Time
}
