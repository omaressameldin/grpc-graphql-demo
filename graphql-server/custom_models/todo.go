package custom_models

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	v1 "github.com/omaressameldin/grpc-graphql-demo/grpc-server/pkg/api/v1"
)

type Todo struct {
	ID          int
	Title       string
	Description string
	IsDone      bool
	Reminder    time.Time
}

func BuildTodo(todo *v1.ToDo) (*Todo, error) {
	t := &Todo{
		ID:          int(todo.GetId()),
		Description: todo.GetDescriptionValue(),
		Title:       todo.GetTitleValue(),
		IsDone:      todo.GetIsDoneValue(),
	}

	reminder, err := ptypes.Timestamp(todo.GetReminderValue())
	if err != nil {
		return nil, fmt.Errorf("reminder field has invalid format-> " + err.Error())
	}
	t.Reminder = reminder

	return t, nil
}
