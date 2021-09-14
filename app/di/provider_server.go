package di

import (
	"context"
	"database/sql"
	"starter-restapi-golang/app/server"
	"starter-restapi-golang/app/validator"

	"github.com/google/wire"
)

func provideContentHandler(ctx context.Context, db *sql.DB, v validator.ContentValidator) server.ContentHandler {
	return server.NewContentHandler(ctx, db, v)
}

var ServerSet = wire.NewSet(
	provideContentHandler,
)
