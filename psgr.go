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
		// Query for write payment information
		// insertPaymentsQuery := `INSERT INTO payments (
		// 	transaction,
		// 	request_id,
		// 	currency,
		// 	provider,
		// 	amount,
		// 	payment_dt,
		// 	bank,
		// 	delivery_cost,
		// 	goods_total,
		// 	custom_fee
		// 	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
		insertDataQuery := `INSERT INTO jsonstore (id, jsondata) VALUES ($1, $2)`
		// Listen data channel to write payment info
		for data := range w.Datas {
			marshaledData, err := json.Marshal(data)
			if err != nil {
				w.Errors <- err
			}
			_, err = w.Db.Exec(insertDataQuery, data.Order_uid, []byte(marshaledData))
			if err != nil {
				w.Errors <- err
			}
		}
	}()
}

func (w *PostgresWriter) Backup(restoreCh chan<- Data) {
	insertStmt := "SELECT jsondata FROM jsonstore"
	rows, err := w.Db.Query(insertStmt)
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
