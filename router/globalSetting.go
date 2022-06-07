package router

import (
	"FD/util"
	"database/sql"
)

var db *sql.DB = util.DB
var log = util.Logger
