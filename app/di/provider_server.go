package di

import (
	"context"
	"database/sql"
	"starter-restapi-golang/app/server"
	"starter-restapi-golang/app/validator"

	"github.com/google/wire"
)

func provideContentServer(ctx context.Context, db *sql.DB, v validator.ContentValidator) server.ContentServer {
	return server.NewContentServer(ctx, db, v)
}

var ServerSet = wire.NewSet(
	provideContentServer,
)
