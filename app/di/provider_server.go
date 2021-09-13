package di

import (
	"context"
	"database/sql"
	"starter-restapi-golang/app/server"

	"github.com/google/wire"
)

func provideUserHandler(ctx context.Context, db *sql.DB) server.UserHandler {
	return server.NewUserHandler(ctx, db)
}

var ServerSet = wire.NewSet(
	provideUserHandler,
)
