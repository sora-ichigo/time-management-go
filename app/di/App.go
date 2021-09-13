package di

import (
	"database/sql"
	"starter-restapi-golang/app/server"
)

type App struct {
	UserHandler server.UserHandler
	Db          *sql.DB
}
