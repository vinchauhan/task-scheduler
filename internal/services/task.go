package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/sanity-io/litter"
	"github.com/vinchauhan/task-scheduler/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var (
	//ErrTaskCannotBeAssigned is for Task that cannot be assigned because of rules
	ErrTaskCannotBeAssigned = errors.New("TASK CANNOT BE ASSIGNED")
	ErrFailedToFindAnyAgentWithSkill = errors.New("COULD NOT FIND ANY SKILLED AGENT")
	ErrSystemFindingAgent = errors.New("SYSTEM ERROR OCCURRED WHILE FINDING AGENT")
	ErrTaskCouldNotBeCreated = errors.New("ERROR SAVING THE TASK TO DB")
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
func (s *Service) CreateTask(ctx context.Context, priority string, skills []string) (TaskOutput , error) {
	log.Printf("Create Task Called")
	var taskOut TaskOutput
	//If skills for the passed task is nil which should be not allowed from UI but still The system cannot decide which agent to be assigned.
	if len(skills) == 0 {
		return taskOut, ErrTaskCannotBeAssigned
	}

    //Filter to get any agent that matches the skill and doesnt have a task assigned yet. It is most effective
    // for the system to find an agent which doesnt have a task and matches the skill rather than assigning task to
    // an agent who already has something on his plate.
	agent, err := findAgentWithSkillAndNoTasks(ctx, s.agentsCollection, skills)
	if err != nil {
		return taskOut, ErrSystemFindingAgent
	}
	//If the agent with right skill was found and had no task to him
	if agent.Id != "" {
		log.Printf("Agent.Id != '' so an skilled agent is found with no task")
		//Create the task in the task collection and then add the id to Agents > task column.
		out, err := s.createNewTaskOnDatabase(ctx, priority, skills, agent.Id)
		log.Printf("Task Created with taskId %s", taskOut.Id)
		if err !=nil {
			return taskOut, ErrTaskCouldNotBeCreated
		}
		//Assign taskOut to out
		taskOut = out
	} else {
		//Find Agents with Skill but already having tasks to work on.
		if priority != "high" {
			//which means incoming task is low priority
			//assuming that there would be matching agent but since they would be working on low or high the incoming
			//task cannot be scheduled.
			log.Printf("Incoming task priority is low and cannot be assigned")
			return taskOut, ErrTaskCannotBeAssigned
		} else {
			//incoming task is high and can be assigned to agent working on low and matching skill
			//Find one agent who matches skill and is working on a low priority task.
			//s.findSkilledAgentWorkingOnLowPriority(ctx, skills)
		}
		log.Printf("Finding Skilled Agents already having task on their plate")
		//s.findSkilledAgentsWithExistingTasks(ctx, skills, priority)
	}
	log.Printf("Returning task with taskId %s", taskOut.Id)
	return taskOut, err
}

func (s *Service) CompleteTask(ctx context.Context, id string) (interface{}, error) {
	//To mark task as complete remove the taskId from the agents collection
	//Update the tasks[] field in agents collection so that the task is not longer assigned to agent
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}
	res, err := s.tasksCollection.UpdateOne(ctx, bson.M{"_id": bson.M{"$eq": objID}},bson.M{"$set": bson.M{"enddatetime": time.Now(),
		                                                                                                   "status":"Completed"}})
	if err != nil {
		log.Fatal(err)
	}
	return res.UpsertedID, nil
}

func (s *Service) BeginTask(ctx context.Context, id string) (string, error) {
	//To mark task as complete remove the taskId from the agents collection
	//Update the tasks[] field in agents collection so that the task is not longer assigned to agent

	//objID, err := primitive.ObjectIDFromHex(id)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//s.tasksCollection.UpdateOne(ctx, bson.M{"_id": bson.M{"$eq": ohanbjID}}, bson.M{"$push" : bson.M{"tasks" : res.InsertedID}})
	return "", nil
}

func (s *Service) findSkilledAgentsWithExistingTasks(ctx context.Context, skills []string, priority string) {
	onlySkillFilter := bson.D{{"skills",
		bson.D{{
			"$all",
			bson.A{skills},
		}},
	}}
	var onlySkillFilterAgents Agent
	cur, err := s.agentsCollection.Find(ctx, onlySkillFilter)
	for cur.Next(ctx) {
		err := cur.Decode(&onlySkillFilterAgents)
		if err != nil {
			log.Fatalf("Error decoding the object from cursor %v\n", err)
		}
		//Convert the slice of string objectId to slice of hex form as Object("AAAAA")
		tasksObjectIds, err := util.ObjectIDArrayFromHex(onlySkillFilterAgents.Tasks)
		if err != nil {
			log.Fatalf("Error occured %v", err)
		}
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
		fmt.Errorf("Error in finding matching skill and low agents %v\n",err)
	}
}

func (s *Service) createNewTaskOnDatabase(ctx context.Context, priority string, skills []string, agentId string) (TaskOutput, error){
	var taskOut TaskOutput
	res, err := s.tasksCollection.InsertOne(ctx, bson.D{{"priority",priority},
		{"skills",skills},
		{"owner", agentId},
		{"status", "Assigned"},
		{"startdatetime",""},
		{"enddatetime", ""}})
	log.Printf("Task inserted in the tasks collection %s", res.InsertedID)

	if err != nil {
		log.Fatal(err)
	}

	//Update the agents collection with received ObjectId of the task
	objID, err := primitive.ObjectIDFromHex(agentId)
	agentColUpdtRes, err := s.agentsCollection.UpdateOne(ctx, bson.M{"_id": bson.M{"$eq": objID}}, bson.M{"$push" : bson.M{"tasks" : res.InsertedID}})
	log.Printf("%d record of agentCollection was updated with taskId %s", agentColUpdtRes.MatchedCount, agentId)

	if err != nil {
		log.Fatal(err)
	}
	taskID := res.InsertedID
	if oid, ok := taskID.(primitive.ObjectID); ok {
		log.Printf("Creating a new task with TaskId: %s", oid.String())
		return TaskOutput{
			Id:       oid.String(),
			Priority: priority,
			Skills:   skills,
			Status:   "Assigned",
			Owner:    agentId,
		}, err
	} else {
		return taskOut, err
	}
}

func findAgentWithSkillAndNoTasks(ctx context.Context, collection *mongo.Collection, skills []string) (Agent, error) {
	var fetchedAgent Agent

	filter := bson.D{{"$and", []bson.D{
		bson.D{{"skills",
			bson.D{{
				"$all",
				bson.A{skills},
			}},
		}},
		bson.D{{"tasks", bson.D{{"$size", 0}}}},
	}}}
	cursor, err := collection.Find(ctx, filter, options.Find().SetLimit(1))
	for cursor.Next(ctx) {
		err := cursor.Decode(&fetchedAgent)
		if err != nil {
			log.Fatal(err)
		}
	}
	if err != nil {
		return fetchedAgent, err
	}
	litter.Dump(fetchedAgent)
	return fetchedAgent, nil
}

func (a *Agent) isAgentStructEmpty() bool {
	return a.Id == ""
}
