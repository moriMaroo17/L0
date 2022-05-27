package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Credentials struct {
	user     string
	dbname   string
	password string
	host     string
	sslmode  string
	port     string
}

type PostgresWriter struct {
	Datas  <-chan Data
	Errors chan<- error
	Db     *sql.DB
}

func (w *PostgresWriter) Write() {
	go func() {
		insertDataQuery := `INSERT INTO jsonstore (id, jsondata) VALUES ($1, $2)`
		// Listen data channel to write payment info
		for data := range w.Datas {
			marshaledData, err := json.Marshal(data)
			if err != nil {
				w.Errors <- err
			}
			_, err = w.Db.Exec(insertDataQuery, data.OrderUid, []byte(marshaledData))
			if err != nil {
				w.Errors <- err
			}
		}
	}()
}

func (w *PostgresWriter) Get(id string, output chan<- Data) {
	getDataQuery := `SELECT jsondata FROM jsonstore WHERE id = $1`
	rows, err := w.Db.Query(getDataQuery, id)
	if err != nil {
		w.Errors <- err
	}
	defer rows.Close()
	var marshData []byte
	var unmarshData Data
	for rows.Next() {
		err := rows.Scan(&marshData)
		if err != nil {
			w.Errors <- err
		} else {
			err := json.Unmarshal(marshData, &unmarshData)
			if err != nil {
				w.Errors <- err
			} else {
				output <- unmarshData
			}
		}
	}
	defer close(output)
}

func (w *PostgresWriter) Backup(restoreCh chan<- Data) {
	backupDataQuery := "SELECT jsondata FROM jsonstore"
	rows, err := w.Db.Query(backupDataQuery)
	if err != nil {
		w.Errors <- err
	}
	defer rows.Close()
	var marshData []byte
	var unmarshData Data
	for rows.Next() {
		err := rows.Scan(&marshData)
		if err != nil {
			w.Errors <- err
		} else {
			err := json.Unmarshal(marshData, &unmarshData)
			if err != nil {
				w.Errors <- err
			} else {
				restoreCh <- unmarshData
			}
		}
	}
	defer close(restoreCh)
}

func NewPostgresWriter(credentials Credentials) (PostgresWriter, chan<- Data, <-chan error) {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=%s port=%s",
		credentials.user, credentials.dbname, credentials.password, credentials.host, credentials.sslmode, credentials.port,
	)

	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	datas := make(chan Data)
	errors := make(chan error)
	return PostgresWriter{Datas: datas, Errors: errors, Db: Db}, datas, errors
}
