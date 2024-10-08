package db

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgconn"

	_ "github.com/jackc/pgx/v4"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func DatabaseConn() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=localhost port=5432 dbname=todo_list user=postgres password=raptor3796")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return db, nil
}