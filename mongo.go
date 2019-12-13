package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func main()  {

		// Set client options
		mongoUrl := "mongodb://"+os.Getenv("MONGO_URL")
		clientOptions := options.Client().ApplyURI(mongoUrl)

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
}
