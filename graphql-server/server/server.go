package main

import (
	"log"
	"net/http"
	"os"
	"google.golang.org/grpc"

	"github.com/99designs/gqlgen/handler"
	"github.com/omaressameldin/grpc-graphql-demo/graphql-server"
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

func connectToService() (*grpc.ClientConn, error) {
	return grpc.Dial("grpc-server:3000", grpc.WithInsecure())
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

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(graphql_server.NewExecutableSchema(graphql_server.Config{Resolvers: resolver})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))


	defer resolver.TodoClient.Close()
}
