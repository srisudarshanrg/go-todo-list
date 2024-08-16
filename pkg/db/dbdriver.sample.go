package db

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgconn"

	_ "github.com/jackc/pgx/v4"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func DatabaseConnSample() (*sql.DB, error) {
	dbSample, err := sql.Open("pgx", "host= port= dbname= user= password=")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return dbSample, nil
}
