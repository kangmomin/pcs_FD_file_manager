package router

import (
	"FD/util"
	"context"
	"database/sql"
)

var (
	db          *sql.DB = util.DB
	log                 = util.Logger
	adminLog            = util.AdminLogger
	argonConfig         = util.ArgonConfig{
		Time:   10,
		Memory: 64 * 1024,
		Thread: 4,
		KeyLen: 32,
	}

	ctx = context.Background()
)
