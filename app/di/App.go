package di

import (
	"database/sql"
	"time_management_slackapp/app/server"
)

type App struct {
	TimePointServer server.TimePointServer
	Db              *sql.DB
}
