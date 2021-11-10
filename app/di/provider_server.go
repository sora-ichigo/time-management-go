package di

import (
	"context"
	"database/sql"
	"time_management_slackapp/app/server"
	"time_management_slackapp/app/validator"

	"github.com/google/wire"
)

func provideContentServer(ctx context.Context, db *sql.DB, v validator.ContentValidator) server.ContentServer {
	return server.NewContentServer(ctx, db, v)
}

var ServerSet = wire.NewSet(
	provideContentServer,
)
