package main

import (
	"context"
	"fmt"
	"github.com/vinchauhan/task-scheduler/internal/handlers"
	"github.com/vinchauhan/task-scheduler/internal/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	databaseURL = "postgresql://root@127.0.0.1:26257/tasker?sslmode=disable"
	webPort     = 3000
)

func main()  {
	//mongoUrl := "mongodb://"+os.Getenv("MONGO_URL")
	// Set client options
	mongoUrl := "mongodb://"+os.Getenv("MONGO_URL")
	clientOptions := options.Client().ApplyURI(mongoUrl)

	// Connect to MongoDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	//s := services.New(db)
	service := services.NewClient(ctx, client)
	//h := handlers.SetUpRoutes(s)
	handler := handlers.SetUpRoutes(service)
	address := fmt.Sprintf(":%d", webPort)
	if err := http.ListenAndServe(address, handler); err != nil {
		log.Fatalf("Failed to start the server : %v\n", err)
	}
}
