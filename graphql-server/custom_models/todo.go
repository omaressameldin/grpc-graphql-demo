package custom_models

import (
	"time"

	v1 "github.com/omaressameldin/grpc-graphql-demo/grpc-server/pkg/api/v1"
)

type Todo struct {
	ID          int
	Title       string
	Description string
	IsDone      bool
	UserID      int
	Reminder    time.Time
}

func BuildTodo(todo *v1.ToDo) *Todo {
	return &Todo{
		ID:          int(todo.GetId()),
		Description: todo.GetDescriptionValue(),
		Title:       todo.GetTitleValue(),
	}
}
