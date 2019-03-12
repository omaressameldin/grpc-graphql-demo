# Readme

## What this is
- A Todo app using [GRPC](https://github.com/grpc/grpc-go) and [GraphQL](https://github.com/99designs/gqlgen)
- the todo app has a GRPC service for creating, updating, deleting, and reading tasks
- it's using [BoltDB](https://github.com/boltdb/bolt) for saving tasks

## How to use
- install docker and docker-compose and make sure docker does not need sudo to run (`sudo groupadd docker` `sudo gpasswd -a $USER docker`)
- once everything is installed use `docker --version` to make sure that docker is running
- run `docker-compose -f docker-compose.dep.yml up --build`
- run `docker-compose up --build` at the root of the project to launch all services
- open [localhost:9181](http://localhost:9181) to go to graphql playground

## Sample Requests
### Queries:
```graphql
query findTodos {
  todos(input: {justRemaining: true}) {
    id
    description
    title
    isDone
    reminder
  }
}
```
```graphql
query findTodo {
  todo(input: {todoId: "11"}) {
    title
    description
    reminder
  }
}
```
### Mutations:
```graphql
mutation createTodo {
  createTodo(input:{description: "the thirty-first todo", title:"todo 31", reminder: "2019-03-07 16:15:00"}) {
    id
    title
    description
    reminder
  }
}
```
```graphql
mutation updateTodo {
  updateTodo(input: {
  	todoId:"26",
  	title:"todo 26"
  	description: "the twenty-sixth todo"
   	reminder:"2019-03-07 16:15:00"
    isDone: true

  }){
    id
    title
    description
    reminder
  }
}
```
```graphql
mutation deleteTodo {
  deleteTodo(input: {todoId: "10"})
}
# delete will only work if you add a token header, e.g:
#  {
#    "Token": "123"
#  }
```

### Subscriptions:
```graphql
subscription remainingTodos {
  remainingTodos
}
```
## Technologies used
- Golang
- GRPC
- graphQL
- Docker
- Docker-compose
- BoltDB
