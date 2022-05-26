package main

import (
	"encoding/json"
	"fmt"

	stan "github.com/nats-io/stan.go"
)

func main() {
	// Start listening nats-streaming service
	lstDataCh, lstErrCh, err := Listen()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	// Connect to postgres
	writer, pgrsDataCh, pgrsErrCh := NewPostgresWriter()
	writer.Write()
	// Listening errors and data channels
	for {
		select {
		case data := <-lstDataCh:
			pgrsDataCh <- data
		case lstErr := <-lstErrCh:
			fmt.Printf("Error while listening: %v\n", lstErr)
		case pgrsErr := <-pgrsErrCh:
			fmt.Printf("Error while working with pgrs: %v\n", pgrsErr)
		}
	}
	defer writer.Db.Close()
}

// Func, which starts listening nats-streaming service and returns two channels or error
func Listen() (<-chan Data, <-chan error, error) {
	// Connect to nats-streaming service
	sc, err := stan.Connect("test-cluster", "r9", stan.NatsURL("nats://localhost:4223"))
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting: %v", err)
	}
	// Created channels for sending data and errors
	data := make(chan Data)
	errors := make(chan error)
	var payment Payment
	_, err = sc.Subscribe("foo",
		func(m *stan.Msg) {
			err := json.Unmarshal(m.Data, &payment)
			if err != nil {
				formatedError := fmt.Errorf("error unmarsh, wrong data in subscription. error: %v", err)
				errors <- formatedError
			}
			data <- Data{payment}
		},
		stan.StartWithLastReceived())
	if err != nil {
		return nil, nil, fmt.Errorf("error subscribing: %v", err)
	}
	return data, errors, nil
}
