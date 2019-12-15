package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/vinchauhan/task-scheduler/internal/handlers"
	"github.com/vinchauhan/task-scheduler/internal/services"
	"log"
	"net/http"
)

const (
	databaseURL = "postgresql://root@127.0.0.1:26257/tasker?sslmode=disable"
	webPort     = 3000
)

func main() {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Could not open connection to Database : ?? %v\n", err)
		return
	}

	defer db.Close()

	//ping the database
	if err = db.Ping(); err != nil {
		log.Fatalf("Could not ping to the Database : ?? %v\n", err)
	}

	log.Printf("Connected to the Database.!")

	//Add the db instance to the service struct needed to talk to the database
	s := services.New(db)
	h := handlers.SetUpRoutes(s)
	address := fmt.Sprintf(":%d", webPort)
	if err := http.ListenAndServe(address, h); err != nil {
		log.Fatalf("Failed to start the server : %v\n", err)
	}

}
