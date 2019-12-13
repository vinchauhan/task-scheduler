package main

import (
	"context"
	"fmt"
	"github.com/vinchauhan/task-scheduler/internal/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

const (
	databaseURL = "postgresql://root@127.0.0.1:26257/tasker?sslmode=disable"
	webPort = 3000
)

func main() {

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	s := services.New(databaseURL)
	address := fmt.Sprintf(":%d", webPort)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Failed to start the server : %v\n", err)
	}
}
