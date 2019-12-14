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
	mongoClient := services.NewMongoClient(ctx, client)
	//h := handlers.SetUpRoutes(s)
	mongoHandler := handlers.SetUpRoutesForMongo(mongoClient)
	address := fmt.Sprintf(":%d", 3000)
	if err := http.ListenAndServe(address, mongoHandler); err != nil {
		log.Fatalf("Failed to start the server : %v\n", err)
	}
}
