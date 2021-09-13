package di

import (
	"database/sql"
	"log"
	"os"

	"github.com/google/wire"
)

func provideDB() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("failed sql open : %v", err)
	}
	return db
}

var ConfigSet = wire.NewSet(
	provideDB,
)
