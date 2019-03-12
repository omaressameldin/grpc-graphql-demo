//go:generate go run github.com/99designs/gqlgen

package graphql_server

import (
	"context"

	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/omaressameldin/grpc-graphql-demo/graphql-server/custom_models"
	v1 "github.com/omaressameldin/grpc-graphql-demo/grpc-server/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateTodo(ctx context.Context, input NewTodo) (*custom_models.Todo, error) {
	c := v1.NewToDoServiceClient(r.TodoClient)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req1 := v1.CreateRequest{
		Api: apiVersion,
		ToDo: &v1.ToDo{
			Title: &v1.ToDo_TitleValue{
				TitleValue: input.Title,
			},
			Description: &v1.ToDo_DescriptionValue{
				DescriptionValue: input.Description,
			},
			IsDone: &v1.ToDo_IsDoneValue{
				IsDoneValue: false,
			},
		},
	}
	if input.Reminder != nil {
		reminder, err := ptypes.TimestampProto(*input.Reminder)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "reminder field has invalid format-> "+err.Error())
		}
		req1.ToDo.Reminder = &v1.ToDo_ReminderValue{
			ReminderValue: reminder,
		}
	}

	res1, err := c.Create(ctx, &req1)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}

	todo, err := custom_models.BuildTodo(res1.GetToDo())
	if err != nil {
		return nil, err
	}
	defer sendRemainingToChannel(c, ctx)
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

	if input.IsDone != nil {
		updatedTodo.IsDone = &v1.ToDo_IsDoneValue{
			IsDoneValue: *input.IsDone,
		}
		defer sendRemainingToChannel(c, ctx)
	}

	if input.Reminder != nil {
		reminder, err := ptypes.TimestampProto(*input.Reminder)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "reminder field has invalid format-> "+err.Error())
		}
		updatedTodo.Reminder = &v1.ToDo_ReminderValue{
			ReminderValue: reminder,
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

	todo, err := custom_models.BuildTodo(res1.GetToDo())
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *mutationResolver) DeleteTodo(ctx context.Context, input DeleteTodo) (bool, error) {
	c := v1.NewToDoServiceClient(r.TodoClient)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req1 := v1.DeleteRequest{
		Api: apiVersion,
		Id:  int64(input.TodoID),
	}

	res1, err := c.Delete(ctx, &req1)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}
	defer sendRemainingToChannel(c, ctx)
	return res1.GetDeleted(), nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Todos(ctx context.Context, input *AllTodos) ([]custom_models.Todo, error) {
	c := v1.NewToDoServiceClient(r.TodoClient)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req1 := v1.ReadAllRequest{
		Api: apiVersion,
	}
	if input != nil {
		req1.JustRemaining = input.JustRemaining
	}

	res1, err := c.ReadAll(ctx, &req1)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}

	todos := []custom_models.Todo{}

	for _, todo := range res1.GetToDos() {
		t, err := custom_models.BuildTodo(todo)
		if err != nil {
			return nil, err
		}
		todos = append(todos, *t)
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

	todo, err := custom_models.BuildTodo(res1.GetToDo())
	if err != nil {
		return nil, err
	}
	return todo, nil
}

var remainingTodosChannel chan int

func sendRemainingToChannel(c v1.ToDoServiceClient, ctx context.Context) {

	req1 := v1.ReadAllRequest{
		Api:           apiVersion,
		JustRemaining: true,
	}

	res1, err := c.ReadAll(ctx, &req1)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}
	if remainingTodosChannel != nil {
		remainingTodosChannel <- len(res1.GetToDos())
	}
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) RemainingTodos(ctx context.Context) (<-chan int, error) {
	remainingTodosChannel = make(chan int, 1)
	go func() {
		<-ctx.Done()
	}()

	return remainingTodosChannel, nil
}
