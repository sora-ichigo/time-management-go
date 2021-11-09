package di

import (
	"database/sql"
	"starter-restapi-golang/app/server"
)

type App struct {
	ContentServer server.ContentServer
	Db            *sql.DB
}
