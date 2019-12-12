package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	databaseURL = "postgresql://root@127.0.0.1:26257/tasker?sslmode=disable"
	webPort = 3000
)

func main() {

	address := fmt.Sprintf(":%d", webPort)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Failed to start the server : %v\n", err)
	}
}