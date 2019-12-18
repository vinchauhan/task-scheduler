package main

import (
	"context"
	"fmt"
	"github.com/vinchauhan/task-scheduler/internal/db"
	"github.com/vinchauhan/task-scheduler/internal/handler"
	"github.com/vinchauhan/task-scheduler/internal/service"
	"log"
	"net/http"
	"os"
)

const webPort = 3000

func main()  {
	mongo := db.New("mongodb://"+os.Getenv("MONGO_URL"))
	ctx := context.Background()
	conn, err := mongo.Plug()

	if err != nil {
		log.Fatalf("Error plugging to the database %v\n", err)
	}

	//Ping the database
	err = conn.Ping(ctx,nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	//Initialize a service instance
	s := service.New(conn)
	//service := service.NewService(ctx, conn)
	//h := handler.SetUpRoutes(s)
	//Initialize a handler instance
	h := handler.NewRouter(s)
	//handler := handler.SetUpRoutes(service)
	address := fmt.Sprintf(":%d", webPort)
	if err := http.ListenAndServe(address, h); err != nil {
		log.Fatalf("Failed to start the server : %v\n", err)
	}
}
