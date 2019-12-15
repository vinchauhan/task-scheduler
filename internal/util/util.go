package util

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

func ObjectIDArrayFromHex(objectIDs []string) ([]primitive.ObjectID, error) {
	var objectIDArray []primitive.ObjectID
	for _ , v := range objectIDs {
		objectID, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			log.Fatalf("Could not get objectId from string")
			return nil, err
		}
		objectIDArray = append(objectIDArray, objectID)
	}
	return objectIDArray, nil
}
