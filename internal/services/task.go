package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/sanity-io/litter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	//ErrTaskCannotBeAssigned is for Task that cannot be assigned because of rules
	ErrTaskCannotBeAssigned = errors.New("TASK CANNOT BE ASSIGNED")
	ErrFailedToFindAnyAgentWithSkill = errors.New("COULD NOT FIND ANY SKILLED AGENT")
	ErrSystemFindingAgent = errors.New("SYSTEM ERROR OCCURRED WHILE FINDING AGENT")
)
type Task struct {
	TaskId string
	Priority string
	Status string
	StartDateTime string
}
type Agent struct {
	Id string  `json:"id" bson:"_id"`
	AgentId string
	Skills []string
	Tasks []Task
}

//CreateTask method creates a task on the database
func (s *Service) CreateTask(ctx context.Context, taskID string, priority string, skills []string) error {

	//If skills for the passed task is nil > The system cannot decide which agent to be assigned.
	if len(skills) == 0 {
		return ErrTaskCannotBeAssigned
	}

	//Check the if all the incoming skills match any of the agent.
	//for i, skill := range skills {
	//	subquery = fmt.Sprintf("$%d", i+1) + subquery
	//	fmt.Printf("index , value %d %s", i, skill)
	//}
	//query := "SELECT * from skillmapping where skill in(" + subquery
	//fmt.Printf("query is %s", query)
	//Connect to agents collection
	collection := s.mongoClient.Database("tasker").Collection("agents")
	var out Agent
	//FindOne should be enough to see if an agent is available with the right skill
	//Build filter for FindOne
	//filter := bson.D{{"skills",
	//		bson.D{{
	//			"$all",
	//			bson.A{"skill1", "skill3"},
	//					}},
	//				}}

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
	cursor, err := collection.Find(ctx, filter, options.Find().SetLimit(1))
	for cursor.Next(ctx) {
		err := cursor.Decode(&out)
		if err != nil {
			log.Fatal(err)
		}
	}

	//err := collection.FindOne(ctx , filter).Decode(&out)
	if err != nil {
		 return ErrSystemFindingAgent
	}
	litter.Dump(out)

	//If no agent found with AND condition of skills && tasks = 0
	//Then the system should find agents which were probably skillful but didnt get picked because they had tasks with them
	if out.AgentId == "" {
		//onlySkillFilter := bson.D{{"skills",
		//		bson.D{{
		//			"$all",
		//			bson.A{skills},
		//					}},
		//				}}


		skillMatchWithLowTask := bson.D{{"$and", []bson.D{
			bson.D{{"skills",
				bson.D{{
					"$all",
					bson.A{skills},
				}},
			}},
			bson.D{{"tasks", bson.D{{"$elemMatch", bson.A{"priority","low"}}}}},
		}}}

		cur, err := collection.Find(ctx, skillMatchWithLowTask)
		for cur.Next(ctx) {
			err := cursor.Decode(&out)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Got agent with low task")
		}

		if err != nil {
			fmt.Errorf("Error in finding matching skill and low agents %v\n",err)
		}

		//return ErrFailedToFindAnyAgentWithSkill
	}


	return nil
}

func (a *Agent) isAgentStructEmpty() bool {
	return a.Id == ""
}
