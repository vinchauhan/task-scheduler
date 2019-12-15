package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/sanity-io/litter"
	"github.com/vinchauhan/task-scheduler/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	//ErrTaskCannotBeAssigned is for Task that cannot be assigned because of rules
	ErrTaskCannotBeAssigned = errors.New("TASK CANNOT BE ASSIGNED")
	ErrFailedToFindAnyAgentWithSkill = errors.New("COULD NOT FIND ANY SKILLED AGENT")
	ErrSystemFindingAgent = errors.New("SYSTEM ERROR OCCURRED WHILE FINDING AGENT")
)
type TaskOutput struct {
	Id string `json:"id" bson:"_id"`
	Priority string
	Skills []string
	Status string
	Owner string
}

type TaskId struct {
	Id string `json:"id" bson:"_id"`
}

type Agent struct {
	Id string  `json:"id" bson:"_id"`
	AgentId string
	Skills []string
	Tasks []string
}

//CreateTask method creates a task on the database
func (s *Service) CreateTask(ctx context.Context, taskID string, priority string, skills []string) error {

	//If skills for the passed task is nil > The system cannot decide which agent to be assigned.
	if len(skills) == 0 {
		return ErrTaskCannotBeAssigned
	}

	//Connect to agents collection
	agentsCollection := s.mongoClient.Database("tasker").Collection("agents")
	tasksCollection := s.mongoClient.Database("tasker").Collection("tasks")
	var fetchedAgent Agent
    //Filter to get any agent that matches the skill and doesnt have a task assigned yet. It is most effective
    // for the system to find an agent which doesnt have a task and matches the skill rather than assigning task to
    // an agent who already has something on his plate.
	filter := bson.D{{"$and", []bson.D{
		bson.D{{"skills",
			bson.D{{
				"$all",
				bson.A{skills},
			}},
		}},
		bson.D{{"tasks", bson.D{{"$size", 0}}}},
	}}}
	cursor, err := agentsCollection.Find(ctx, filter, options.Find().SetLimit(1))
	for cursor.Next(ctx) {
		err := cursor.Decode(&fetchedAgent)
		if err != nil {
			log.Fatal(err)
		}
	}
	if err != nil {
		 return ErrSystemFindingAgent
	}
	litter.Dump(fetchedAgent)

	//If no agent found with condition of matching skills && tasks = 0 then the decoded struct will be empty on fields
	//The system should find agents which were probably skillful but didnt get picked because they had tasks with them
	if fetchedAgent.Id == "" {
		onlySkillFilter := bson.D{{"skills",
			bson.D{{
				"$all",
				bson.A{skills},
			}},
		}}
		var onlySkillFilterAgents Agent
		cur, err := agentsCollection.Find(ctx, onlySkillFilter)
		for cur.Next(ctx) {
			err := cur.Decode(&onlySkillFilterAgents)
			if err != nil {
				log.Fatalf("Error decoding the object from cursor %v\n", err)
			}
			//log.Printf("Got agent matching skill with Mongo document Id %s", onlySkillFilterAgents.Id)
			//log.Printf("Got agent matching skill with Agent Id %s", onlySkillFilterAgents.AgentId)
			//log.Printf("Got agent matching skill with Task Id %s", onlySkillFilterAgents.Tasks)

			//Convert the slice of string objectId to slice of hex form as Object("AAAAA")
			//var objectIDArray []primitive.ObjectID
			//for _, v := range onlySkillFilterAgents.Tasks {
			//	objectID, err := primitive.ObjectIDFromHex(v)
			//	if err != nil {
			//		log.Fatalf("Could not get objectId from string")
			//	}
			//	objectIDArray = append(objectIDArray, objectID)
			//}
			tasksObjectIds, err := util.ObjectIDArrayFromHex(onlySkillFilterAgents.Tasks)
			if err != nil {
				log.Fatalf("Error occured %v", err)
			}
			//User application-level join to get agent's tasks
			skillAgntCur, err := tasksCollection.Find(ctx, bson.M{"_id": bson.M{"$in": tasksObjectIds}})
			for skillAgntCur.Next(ctx) {
				var taskForSkilledAgent TaskOutput
				log.Printf("Looping cursor object for Tasks")
				err := skillAgntCur.Decode(&taskForSkilledAgent)
				if err != nil {
					log.Fatalf("Error decoding the object from cursor %v\n", err)
				}
				log.Printf("Got task for agent matching skill with Mongo document Id %s", taskForSkilledAgent.Id)
				log.Printf("Got task for agent matching skill with Priority Id %s", taskForSkilledAgent.Priority)
				log.Printf("Got task for agent matching skill with status %s", taskForSkilledAgent.Skills)
			}

		}

		if err != nil {
			fmt.Errorf("Error in finding matching skill and low agents %v\n",err)
		}

		//return ErrFailedToFindAnyAgentWithSkill
	} else {
		//Create the task in the task collection and then add the id to Agents > task column.
		taskCollection := s.mongoClient.Database("tasker").Collection("tasks")
		res, err := taskCollection.InsertOne(ctx, bson.D{{"priority",priority},
			{"skills",skills},
			{"owner", fetchedAgent.AgentId}})
		log.Printf("Task inserted in the tasks collection %s", res.InsertedID)
		if err != nil {
			log.Fatal(err)
		}

		//Update the agents tasks.id with received Id from inserting the above task
		objID, err := primitive.ObjectIDFromHex(fetchedAgent.Id)
		agentColUpdtRes, err := agentsCollection.UpdateOne(ctx, bson.M{"_id": bson.M{"$eq": objID}}, bson.M{"$push" : bson.M{"tasks" : res.InsertedID}})
		log.Printf("%d record of agentCollection was updated with taskId %s", agentColUpdtRes.MatchedCount, fetchedAgent.Id)

	}
	return nil
}

func (a *Agent) isAgentStructEmpty() bool {
	return a.Id == ""
}
