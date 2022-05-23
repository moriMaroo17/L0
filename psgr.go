package main

import (
	"database/sql"
	"fmt"
	"log"
)

type PostgresWriter struct {
	Datas  <-chan Data
	Errors chan<- error
	Db     *sql.DB
}

func (w *PostgresWriter) Write() {

}

func NewPostgresWriter() (PostgresWriter, chan<- Data, <-chan error) {
	connStr := "user=service dbname=test_db password=servicepassword host=localhost sslmode=disable"

	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", Db)
	datas := make(chan Data)
	errors := make(chan error)
	return PostgresWriter{Datas: datas, Errors: errors, Db: Db}, datas, errors
}
