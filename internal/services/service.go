package services

import (
	"context"
	"database/sql"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	db *sql.DB
	mongoClient *mongo.Client
	ctx context.Context
}

func New(db *sql.DB) *Service {
	return &Service{db: db}
}

func NewMongoClient(ctx context.Context, mongoClient *mongo.Client) *Service  {
	return &Service{
		db: nil,
		mongoClient:mongoClient,
		ctx:ctx,
	}
}
