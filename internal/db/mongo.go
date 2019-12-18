package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//MongoDb struct will hold things needed to connect to mongo
type MongoDB struct {
	url  string
}

//Plug method provide implementation to connect to DB and return an handle
func (d *MongoDB) Plug() (*mongo.Client, error) {
	var con *mongo.Client
	clientOptions := options.Client().ApplyURI(d.url)

	//Connect to mongo db
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Connect to mongo db and return a connection/client object
	con, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return con, err
	}
	return con, err
}

//New function gives the instance of MongoDb
func New(url string) *MongoDB {
	return &MongoDB{
		url:url,
	}
}
