package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
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
	datas := make(chan Data)
	errors := make(chan error)
	go func() {
		for data := range datas {
			fmt.Printf("%v\n", data)
		}
	}()
	return PostgresWriter{Datas: datas, Errors: errors, Db: Db}, datas, errors
}
