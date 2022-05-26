package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PostgresWriter struct {
	Datas    <-chan Data
	BackupCh chan<- Data
	Errors   chan<- error
	Db       *sql.DB
}

func (w *PostgresWriter) Write() {
	go func() {
		// Query for write payment information
		insertPaymentsQuery := `INSERT INTO payments (
			transaction,
			request_id,
			currency,
			provider,
			amount,
			payment_dt,
			bank,
			delivery_cost,
			goods_total,
			custom_fee
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
		// Listen data channel to write payment info
		for data := range w.Datas {
			_, err := w.Db.Exec(insertPaymentsQuery,
				data.Payment.Transaction,
				data.Payment.Request_id,
				data.Payment.Currency,
				data.Payment.Provider,
				data.Payment.Amount,
				data.Payment.Payment_dt,
				data.Payment.Bank,
				data.Payment.Delivery_cost,
				data.Payment.Goods_total,
				data.Payment.Custom_fee,
			)
			if err != nil {
				w.Errors <- err
			}
		}
	}()
}

func (w *PostgresWriter) Backup() {
	insertStmt := "SELECT * FROM payments"
	rows, err := w.Db.Query(insertStmt)
	if err != nil {
		w.Errors <- err
	}
	defer rows.Close()
	pmnt := Payment{}
	for rows.Next() {
		err := rows.Scan(
			&pmnt.Transaction, &pmnt.Request_id,
			&pmnt.Currency, &pmnt.Provider,
			&pmnt.Amount, &pmnt.Payment_dt,
			&pmnt.Bank, &pmnt.Delivery_cost,
			&pmnt.Goods_total, &pmnt.Custom_fee,
		)
		if err != nil {
			w.Errors <- err
		}
	}
	fmt.Printf("%v\n", pmnt)
}

func (w *PostgresWriter) CheckTablesExists() {

}

func NewPostgresWriter() (PostgresWriter, chan<- Data, <-chan error) {
	connStr := "user=service dbname=test_db password=servicepassword host=localhost sslmode=disable port=5455"

	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	datas := make(chan Data)
	backUpCh := make(chan Data)
	errors := make(chan error)
	return PostgresWriter{Datas: datas, BackupCh: backUpCh, Errors: errors, Db: Db}, datas, errors
}
