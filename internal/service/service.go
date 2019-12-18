package service

import (
	"context"
	"github.com/vinchauhan/task-scheduler/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Service struct {
	db *mongo.Client
	ctx context.Context
	tasksCollection *mongo.Collection
	agentsCollection *mongo.Collection
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
			err := taskCursor.Decode(&taskForSkilledAgent)
			if err != nil {
				log.Fatalf("Error decoding the object from cursor %v\n", err)
			}
		}
	}
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}
}

func New(client *mongo.Client) *Service {
	return &Service{
		db:client,
		tasksCollection:client.Database("tasker").Collection("tasks"),
		agentsCollection:client.Database("tasker").Collection("agents"),
	}
}
