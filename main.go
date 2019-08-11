package main

import (
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
)

func main() {

	// Connect to a Postgresql DB called notifications
	listener := pq.NewListener("postgres://tom@localhost:5432/notifications?sslmode=disable",
		1*time.Second, 10*time.Second, failureCallback) // TODO: Check these

	listener.Listen("notes")

	for {
		select {
		case <-listener.Notify:
			fmt.Println("received notification, new work available")
		case <-time.After(90 * time.Second):
			go listener.Ping()
			// Check if there's more work available, just in case it takes
			// a while for the Listener to notice connection loss and
			// reconnect.
			fmt.Println("received no work for 90 seconds, checking for new work")
		}
	}
}

func failureCallback(event pq.ListenerEventType, err error) {
	if err != nil {
		log.Fatal(err)
	}
}
