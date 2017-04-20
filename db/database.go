package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sqlx.DB
}

// Open ...
func Open(dbType, dbConn string) (*DB, error) {
	db, err := sqlx.Connect(dbType, dbConn)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// InitDbSchemas ...
func (db *DB) InitDbSchemas() error {
	_, err := db.Exec(schema)
	return err
}

const schema = `
CREATE TABLE IF NOT EXISTS cards (
	id integer primary key autoincrement,
	type tinyint not null,
	front text not null,
	back text not null,
	known boolean default 0
);
`
