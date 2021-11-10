package di

import (
	"database/sql"
	"time_management_slackapp/app/server"
)

type App struct {
	ContentServer server.ContentServer
	Db            *sql.DB
}
