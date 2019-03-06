//go:generate go run github.com/99designs/gqlgen

package graphql_server

import (
	"context"
	"fmt"

	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/omaressameldin/grpc-graphql-demo/graphql-server/custom_models"
	v1 "github.com/omaressameldin/grpc-graphql-demo/grpc-server/pkg/api/v1"
	"google.golang.org/grpc"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

type Resolver struct {
	TodoClient *grpc.ClientConn
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Todo() TodoResolver {
	return &todoResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateTodo(ctx context.Context, input NewTodo) (*custom_models.Todo, error) {
	c := v1.NewToDoServiceClient(r.TodoClient)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := time.Now().In(time.UTC)
	reminder, _ := ptypes.TimestampProto(t)

	req1 := v1.CreateRequest{
		Api: apiVersion,
		ToDo: &v1.ToDo{
			Title: &v1.ToDo_TitleValue{
				TitleValue: input.Title,
			},
			Description: &v1.ToDo_DescriptionValue{
				DescriptionValue: input.Description,
			},
			Reminder: &v1.ToDo_ReminderValue{
				ReminderValue: reminder,
			},
		},
	}

	res1, err := c.Create(ctx, &req1)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}

	todo := custom_models.BuildTodo(res1.GetToDo())
	return todo, nil
}

func (r *mutationResolver) UpdateTodo(ctx context.Context, input UpdateTodo) (*custom_models.Todo, error) {
	c := v1.NewToDoServiceClient(r.TodoClient)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	updatedTodo := v1.ToDo{
		Id: int64(input.TodoID),
	}
	if input.Title != nil {
		updatedTodo.Title = &v1.ToDo_TitleValue{
			TitleValue: *input.Title,
		}
	}
	if input.Description != nil {
		updatedTodo.Description = &v1.ToDo_DescriptionValue{
			DescriptionValue: *input.Description,
		}
	}

	req1 := v1.UpdateRequest{
		Api:  apiVersion,
		ToDo: &updatedTodo,
	}

	res1, err := c.Update(ctx, &req1)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}

	todo := custom_models.BuildTodo(res1.GetToDo())
	return todo, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Todos(ctx context.Context) ([]custom_models.Todo, error) {
	c := v1.NewToDoServiceClient(r.TodoClient)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req1 := v1.ReadAllRequest{
		Api: apiVersion,
	}

	res1, err := c.ReadAll(ctx, &req1)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}

	todos := []custom_models.Todo{}

	for _, todo := range res1.GetToDos() {
		todos = append(todos, *custom_models.BuildTodo(todo))
	}

	return todos, nil
}

func (r *queryResolver) Todo(ctx context.Context, input ReadTodo) (*custom_models.Todo, error) {
	c := v1.NewToDoServiceClient(r.TodoClient)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req1 := v1.ReadRequest{
		Api: apiVersion,
		Id:  int64(input.TodoID),
	}

	res1, err := c.Read(ctx, &req1)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}

	todo := custom_models.BuildTodo(res1.GetToDo())
	return todo, nil
}

type todoResolver struct{ *Resolver }

func (r *todoResolver) User(ctx context.Context, obj *custom_models.Todo) (*User, error) {
	return &User{ID: obj.UserID, Name: fmt.Sprintf("user :%d", obj.UserID)}, nil
}
