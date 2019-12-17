package services

import (
	"context"
	"github.com/vinchauhan/task-scheduler/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Service struct {
	mongoClient *mongo.Client
	ctx context.Context
	tasksCollection *mongo.Collection
	agentsCollection *mongo.Collection
}

//Designed using interface to support mocking
type ServiceLayer interface {
	GetService(ctx context.Context, client *mongo.Client) *Service
}

func (s* Service) GetService(ctx context.Context, client *mongo.Client) *Service {
	return &Service{
		mongoClient:client,
		ctx:ctx,
		tasksCollection:client.Database("tasker").Collection("tasks"),
		agentsCollection:client.Database("tasker").Collection("agents"),
	}
}

func (s *Service) findSkilledAgentWorkingOnLowPriority(ctx context.Context, skills []string) {

	filter := bson.D{{"skills",
		bson.D{{
			"$all",
			bson.A{skills},
		}},
	}}
	var agent Agent
	cur, err := s.agentsCollection.Find(ctx, filter)
	for cur.Next(ctx) {
		err := cur.Decode(&agent)
		if err != nil {
			log.Fatalf("Error decoding the object from cursor %v\n", err)
		}
		//Convert the slice of string objectId to slice of hex form as Object("AAAAA")
		tasksObjectIds, err := util.ObjectIDArrayFromHex(agent.Tasks)
		//User application-level join to get agent's tasks
		taskCursor, err := s.tasksCollection.Find(ctx, bson.M{"_id": bson.M{"$in": tasksObjectIds}})
		for taskCursor.Next(ctx) {
			var taskForSkilledAgent TaskOutput
			log.Printf("Looping cursor object for Tasks")
			err := taskCursor.Decode(&taskForSkilledAgent)
			if err != nil {
				log.Fatalf("Error decoding the object from cursor %v\n", err)
			}
			log.Printf("Got task for agent matching skill with Mongo document Id %s", taskForSkilledAgent.Id)
			log.Printf("Got task for agent matching skill with Priority Id %s", taskForSkilledAgent.Priority)
			log.Printf("Got task for agent matching skill with status %s", taskForSkilledAgent.Skills)
		}
	}
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}
}

func NewService(ctx context.Context, mongoClient *mongo.Client) *Service  {
	return &Service{
		mongoClient:mongoClient,
		ctx:ctx,
		tasksCollection:mongoClient.Database("tasker").Collection("tasks"),
		agentsCollection:mongoClient.Database("tasker").Collection("agents"),
	}
}
