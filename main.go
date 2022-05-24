package main

import (
	"encoding/json"
	"fmt"

	stan "github.com/nats-io/stan.go"
)

func main() {
	lstDataCh, lstErrCh, err := Listen()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	_, pgrsDataCh, pgrsErrCh := NewPostgresWriter()
	for {
		select {
		case data := <-lstDataCh:
			pgrsDataCh <- data
		case err := <-lstErrCh:
			fmt.Printf("Error: %v\n", err)
		case err := <-pgrsErrCh:
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func Listen() (<-chan Data, <-chan error, error) {
	sc, err := stan.Connect("test-cluster", "r9", stan.NatsURL("nats://localhost:4223"))
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting: %v", err)
	}
	data := make(chan Data)
	errors := make(chan error)
	var payment Payment
	_, err = sc.Subscribe("foo",
		func(m *stan.Msg) {
			// fmt.Printf("Got message. Data: %v\n", string(m.Data))
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
