package internal

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/apex/log"
	_ "github.com/lib/pq"
)

func NewPostgresConn(database *Database) *sql.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		database.DatabaseUsername,
		database.DatabasePassword,
		database.DatabaseHost,
		database.DatabasePort,
		database.DatabaseName,
		database.DatabaseSslMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("[NewDatabase] error opening connection to database: %v\n", err)
	}

	db.SetMaxOpenConns(database.ConnectionPool.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(database.ConnectionPool.MaxLifetime) * time.Second)
	db.SetMaxIdleConns(database.ConnectionPool.MaxIdleConns)
	db.SetConnMaxIdleTime(time.Duration(database.ConnectionPool.MaxIdleTime) * time.Second)
	return db
}
