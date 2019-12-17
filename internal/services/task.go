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
	ErrTaskDoesNotHaveSkills = errors.New("TASK REQUIRES SKILLS")
	ErrFailedToFindAnyAgentWithSkill = errors.New("COULD NOT FIND ANY SKILLED AGENT")
	ErrSystemFindingAgent = errors.New("SYSTEM ERROR OCCURRED WHILE FINDING AGENT")
	ErrTaskCouldNotBeCreated = errors.New("ERROR SAVING THE TASK TO DB")
	ErrTaskCannotBeAssigned = errors.New("TASK COULD NOT BE ASSIGNED")
)
type TaskOutput struct {
	Id string `json:"id" bson:"_id"`
	Priority string
	Skills []string
	Status string
	AgentId string
}

type Task struct {
	Id string `json:"id" bson:"_id"`
	Priority string
	Skills []string
	Status string
	AgentId string
	StartDateTime time.Time
	EndDateTime time.Time
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
		return taskOut, ErrTaskDoesNotHaveSkills
	}

    //Filter to get any agent that matches the skill and doesnt have a task assigned yet.
    //It is most effective for the system to find an agent who doesnt have a task and matches the skill
    //rather than assigning task to an agent who already has something on his plate.
	agent, err := findAgentWithSkillAndNoTasks(ctx, s.agentsCollection, skills)
	if err != nil {
		return taskOut, ErrSystemFindingAgent
	}
	//If the agent with right skill was found and had no task to him
	if agent.Id != "" {
		log.Printf("Found skilled agent is found with no task")
		//Create the task in the task collection and then add the id to Agents > task column.
		taskOut, err = s.createNewTaskOnDatabase(ctx, priority, skills, agent.Id)
		log.Printf("Task Created with taskId %s", taskOut.Id)
		if err !=nil {
			return taskOut, ErrTaskCouldNotBeCreated
		}
		//Assign taskOut to out
		//taskOut = out
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
			//Find agents who all matches the skills
			//s.findSkilledAgentWorkingOnLowPriority(ctx, skills)
			agents, err := s.findSkilledAgentsIds(ctx,skills)
			if err != nil {
				return taskOut, err
			}
			log.Printf("Found %d agents which matches the skill", len(agents))

			//Find the task for all of these agentsIds that has a prority low and startdatetime initialized ( Only low can be picked )
			//Task cannot be assigned to agent working on same or higher priority
			//Application level join needed to get to all tasks
			s.findSkilledAgentJustStartedWithLowPriorityTask(ctx, agents)


		}
		log.Printf("All Skilled agents are busy - Trying to schedule task")
		//s.findSkilledAgentsWithExistingTasks(ctx, skills, priority)
	}
	log.Printf("Returning task with taskId %s", taskOut.Id)
	return taskOut, err
}

func(s *Service) findSkilledAgentJustStartedWithLowPriorityTask(ctx context.Context, agents []string) (Agent, error) {

	var agent Agent
	var task Task
	var minTimeSince time.Duration
	minTimeSince = 0
	filter := bson.D{{"agentid",
			bson.D{{
				"$in",
				bson.A{agents},
			}},
		}}

	cur, err := s.tasksCollection.Find(ctx, filter)
	if err != nil {
		return agent, err
	}
	//Loop through the found tasks to see which one has low priority and lowest time since starttime
	for cur.Next(ctx) {
		err := cur.Decode(&task)
		if err != nil {
			return agent,err
		}
		if task.Status != "high" {
			if time.Since(task.StartDateTime) > minTimeSince {
				minTimeSince = time.Since(task.StartDateTime)
			}
		}

	}
}
//	filter := bson.M{"agentid": bson.M{"$in": agents}},
//			  bson.M{"$elemMatch": bson.M{priority: "low"}},
//	cur, err := s.tasksCollection.Find(ctx, bson.M{"agentid": bson.M{"$in": agents}},
//		bson.M{"$elemMatch": bson.M{priority: "low"}}	)
//}

func (s *Service) CompleteTask(ctx context.Context, id string) (string, error) {
	var task TaskOutput
	//To mark task as complete remove the taskId from the agents collection
	//Update task's status and enddatetime in the tasks collection
	taskObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Marking the task with Id : {%s} complete", id)
	err = s.tasksCollection.FindOneAndUpdate(ctx, bson.M{"_id": bson.M{"$eq": taskObjID}},
		bson.M{"$set": bson.M{"enddatetime": time.Now(),"status":"Completed"}}).Decode(&task)
	if err != nil {
		log.Fatal(err)
	}
	agentObjectID , err := primitive.ObjectIDFromHex(task.AgentId)
	//Remove the taskId from the list of tasks for Agents
	log.Printf("Removing the task with Id : {%s} from the Agent Id : {%s} tasks list", id, task.AgentId)
	_, err = s.agentsCollection.UpdateOne(ctx, bson.M{"_id": bson.M{"$eq": agentObjectID}},
		bson.M{"$pull": bson.M{"tasks": taskObjID}})

	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("Task Id : %s Completed", task.Id), nil
}

func (s *Service) BeginTask(ctx context.Context, id string) (string, error) {
	//To Mark a task started - We need up update the status and add a startdatetime
	var task TaskOutput
	taskObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}

	err = s.tasksCollection.FindOneAndUpdate(ctx, bson.M{"_id": bson.M{"$eq": taskObjID}},
		bson.M{"$set": bson.M{"startdatetime": time.Now(),"status":"InProgress"}}).Decode(&task)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("Task Id : %s Started", task.Id), nil
}

func (s *Service) findSkilledAgentsIds(ctx context.Context, skills []string) ([]string, error) {
	var listOfAgentIds []string
	var agent Agent
	filter := bson.D{{"skills",
		bson.D{{
			"$all",
			bson.A{skills},
		}},
	}}
	cur, err := s.agentsCollection.Find(ctx, filter)
	if err != nil {
		return listOfAgentIds, err
	}

	for cur.Next(ctx) {
		err := cur.Decode(&agent)
		if err != nil {
			return listOfAgentIds,err
		}
		listOfAgentIds = append(listOfAgentIds, agent.Id)
	}

	return listOfAgentIds, nil
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
		{"agentid", agentId},
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
			AgentId:    agentId,
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
