package main

import (
	"encoding/json"
	"fmt"

	stan "github.com/nats-io/stan.go"
)

func main() {
	Listen()
}

func Listen() error {
	sc, err := stan.Connect("test-cluster", "r9", stan.NatsURL("nats://localhost:4223"))
	if err != nil {
		return fmt.Errorf("error connecting: %v", err)
	}
	var payment Payment
	sub, err := sc.Subscribe("foo",
		func(m *stan.Msg) {
			fmt.Printf("Got message. Data: %v\n", string(m.Data))
			err := json.Unmarshal(m.Data, &payment)
			if err != nil {
				fmt.Printf("Error unmarsh, wrong data in subscription. Error: %v\n", err)
			}
		},
		stan.StartWithLastReceived())
	defer sub.Unsubscribe()
	if err != nil {
		return fmt.Errorf("error subscribing: %v", err)
	}
	return nil
}
