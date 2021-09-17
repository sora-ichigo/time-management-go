package di

import (
	"database/sql"
	"starter-restapi-golang/app/server"
)

type App struct {
	ContentHandler server.ContentHandler
	Db             *sql.DB
}
