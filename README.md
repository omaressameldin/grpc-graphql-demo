# Readme

## What this is
- A Todo app using GRPC and GraphQl
- the todo app has a GRPC service for creating, updating, deleting, and reading tasks
- it's using boltDB for saving tasks

## How to use
- install docker and docker-compose and make sure docker does not need sudo to run (`sudo groupadd docker` `sudo gpasswd -a $USER docker`)
- once everything is installed use `docker --version` to make sure that docker is running
- run `docker-compose -f docker-compose.dep.yml up --build`
- run `docker-compose up --build` at the root of the project to launch all services
- open [localhost:9181](http://localhost:9181) to go to graphql playground

## Technologies used
- golang
- GRPC
- graphQL
- docker
- docker-compose
