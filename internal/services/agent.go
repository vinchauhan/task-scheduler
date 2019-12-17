package services

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrFailedToFetchAgents = errors.New("FAILED TO LIST ALL AGENTS")
	ErrFailedToDecodeAgents = errors.New("SYSTEM ERROR")
)

func (s *Service) GetAgents(ctx context.Context) ([] Agent, error) {
	//Get All Agents
	var agent Agent
	var agents []Agent
	cur, err := s.agentsCollection.Find(ctx, bson.D{})
	if err != nil {
		return agents, ErrFailedToFetchAgents
	}
	for cur.Next(ctx) {
		err := cur.Decode(&agent)
		if err != nil {
			return agents, ErrFailedToDecodeAgents
		}
		agents = append(agents, agent)
	}

	return agents, nil
}
