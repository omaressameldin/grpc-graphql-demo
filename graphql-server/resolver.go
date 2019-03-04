//go:generate go run github.com/99designs/gqlgen

package graphql_server

import (
	"context"
	"fmt"

	"github.com/omaressameldin/grpc-graphql-demo/graphql-server/custom_models"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{
	todos []custom_models.Todo
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
	todo := &custom_models.Todo{
		Title:   input.Title,
		Description: input.Description,
		ID:     len(r.todos) + 1,
		UserID: input.UserID,
	}
	r.todos = append(r.todos, *todo)
	return todo, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Todos(ctx context.Context) ([]custom_models.Todo, error) {
	return r.todos, nil
}

type todoResolver struct{ *Resolver }

func (r *todoResolver) User(ctx context.Context, obj *custom_models.Todo) (*User, error) {
	return &User{ID: obj.UserID, Name: fmt.Sprintf("user :%d", obj.UserID)}, nil
}
