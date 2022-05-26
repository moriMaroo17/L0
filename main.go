package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	stan "github.com/nats-io/stan.go"
)

var cacher Cache

func initService() {
	creds := Credentials{
		user:     "service",
		dbname:   "test_db",
		password: "servicepassword",
		host:     "localhost",
		sslmode:  "disable",
		port:     "5455",
	}
	// Start listening nats-streaming service
	lstDataCh, lstErrCh, err := Listen()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	// Connect to postgres
	writer, pgrsDataCh, pgrsErrCh := NewPostgresWriter(creds)
	writer.Write()
	// Create cacher and restore data from postgres
	cacher = NewCache()
	if cacher.CheckEmpty() {
		go cacher.Restore(&writer)
	}
	// Listening errors and data channels
	for {
		select {
		case data := <-lstDataCh:
			go cacher.Put(data.Order_uid, data)
			pgrsDataCh <- data
		case lstErr := <-lstErrCh:
			fmt.Printf("Error while listening: %v\n", lstErr)
		case pgrsErr := <-pgrsErrCh:
			fmt.Printf("Error while working with pgrs: %v\n", pgrsErr)
		}
	}
	defer writer.Db.Close()
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	data, err := cacher.Get("b563feb7b2b84b6test")
	fmt.Printf("data: %v\n", cacher.memoryCache)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%v\n", data)
}

func main() {
	go initService()

	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Func, which starts listening nats-streaming service and returns two channels or error
func Listen() (<-chan Data, <-chan error, error) {
	// Connect to nats-streaming service
	sc, err := stan.Connect("test-cluster", "r9", stan.NatsURL("nats://localhost:4223"))
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting: %v", err)
	}
	// Created channels for sending data and errors
	dataCh := make(chan Data)
	errors := make(chan error)
	var data Data
	_, err = sc.Subscribe("foo",
		func(m *stan.Msg) {
			err := json.Unmarshal(m.Data, &data)
			if err != nil {
				formatedError := fmt.Errorf("error unmarsh, wrong data in subscription. error: %v", err)
				errors <- formatedError
			} else {
				dataCh <- data
			}
		},
		stan.StartWithLastReceived())
	if err != nil {
		return nil, nil, fmt.Errorf("error subscribing: %v", err)
	}
	return dataCh, errors, nil
}
