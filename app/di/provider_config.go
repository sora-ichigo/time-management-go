package di

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/google/wire"
)

func DSN() string {
	dsn := os.Getenv("DSN")

	if dsn != "" {
		return dsn
	}

	return "admin:admin@(127.0.0.1:3307)/api-server?parseTime=true"
}

func provideDB() *sql.DB {
	db, err := sql.Open("mysql", DSN())
	if err != nil {
		log.Fatalf("failed sql open : %v", err)
	}
	return db
}

var ConfigSet = wire.NewSet(
	provideDB,
)
