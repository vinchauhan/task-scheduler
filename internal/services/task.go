package services

import (
	"context"
	"errors"
	"fmt"
)

var (
	//ErrTaskCannotBeAssigned is for Task that cannot be assigned because of rules
	ErrTaskCannotBeAssigned = errors.New("TASK CANNOT BE ASSIGNED")
)

//CreateTask method creates a task on the database
func (s *Service) CreateTask(ctx context.Context, taskID string, priority string, skills []string) error {

	var subquery = "and"
	//If skills for the passed task is nil > The system cannot decide which agent to be assigned.
	if len(skills) == 0 {
		return ErrTaskCannotBeAssigned
	}

	//Check the skill mapping table to see if the skills for the task matches any agent.
	for i, skill := range skills {
		subquery = fmt.Sprintf("skill=$%d ", i+1) + subquery
		fmt.Printf("index , value %d %s", i, skill)
	}
	query := "SELECT * from skillmapping where " + subquery
	fmt.Printf("query is %s", query)
	return nil
}
