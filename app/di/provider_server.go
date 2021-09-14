package di

import (
	"context"
	"database/sql"
	"starter-restapi-golang/app/server"

	"github.com/google/wire"
)

func provideContentHandler(ctx context.Context, db *sql.DB) server.ContentHandler {
	return server.NewContentHandler(ctx, db)
}

var ServerSet = wire.NewSet(
	provideContentHandler,
)
