package di

import (
	"database/sql"
	"time_management_slackapp/app/server"

	"github.com/google/wire"
)

func provideTimePointServer(db *sql.DB) server.TimePointServer {
	return server.NewTimePointServer(db)
}

var ServerSet = wire.NewSet(
	provideTimePointServer,
)
