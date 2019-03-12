package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/grpc"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	graphql_server "github.com/omaressameldin/grpc-graphql-demo/graphql-server"
)

func newGraphQLResolver() (*graphql_server.Resolver, error) {
	todoClient, err := connectToService()
	if err != nil {
		return nil, err
	}

	return &graphql_server.Resolver{
		TodoClient: todoClient,
	}, nil
}

const defaultPort = "8080"
const authorizationKey = "Token"

type contextKey string

func connectToService() (*grpc.ClientConn, error) {
	return grpc.Dial("grpc-server:3000", grpc.WithInsecure())
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		auth := r.Header.Get(authorizationKey)
		if auth != "" {
			// Write your fancy token introspection logic here and if valid user then pass appropriate key in header
			// IMPORTANT: DO NOT HANDLE UNAUTHORISED USER HERE
			ctx = context.WithValue(ctx, contextKey(authorizationKey), auth)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	resolver, err := newGraphQLResolver()
	if err != nil {
		log.Fatalf("couldnt connect to server: %v", err)
	}

	c := graphql_server.Config{Resolvers: resolver}

	c.Directives.AuthenticationToken = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		token := ctx.Value(contextKey(authorizationKey))
		if token != nil {
			return next(ctx)
		}
		return nil, fmt.Errorf("You are not authorized to do that action")
	}

	rootHandler := handler.GraphQL(graphql_server.NewExecutableSchema(c))

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", authMiddleware(rootHandler))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	defer resolver.TodoClient.Close()
}
