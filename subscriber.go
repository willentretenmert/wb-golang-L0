package main

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/nats-io/stan.go"
	"log"
	"time"
)

const (
	//natsURL = "nats://localhost:4222"
	natsURL     = "nats://stan-container:4222"
	clusterID   = "test-cluster"
	clientID    = "your-client-id"
	channelName = "test-subject"
)

func SubscribeToNATS(conn *pgx.Conn) {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL), stan.ConnectWait(5*time.Second))
	if err != nil {
		log.Printf("Can't connect to NATS Streaming: %v\n", err)
		return
	}
	defer sc.Close()

	_, err = sc.Subscribe(channelName, func(m *stan.Msg) {
		//fmt.Printf("Recieved a message from %s: %s\n", channelName, m)
		data := m.Data
		var order Order

		err1 := json.Unmarshal(data, &order)
		if err1 != nil {
			fmt.Printf("error %v\n", err1)
			return
		}

		err2 := insertOrder(conn, order)
		if err2 != nil {
			fmt.Printf("inserting error: %v\n", err2)
			return
		}

		loadCacheFromDB(conn, &cache)

		fmt.Printf("Received NATS message: %+v\n", order)
	}, stan.StartWithLastReceived())

	if err != nil {
		fmt.Printf("Error on NATS subscription: %v\n", err)
	}

	select {}
}
