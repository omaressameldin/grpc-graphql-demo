FROM znly/protoc:0.4.0
WORKDIR /usr/src/app
COPY ./proto ./proto
RUN mkdir -p pkg/api/v1
RUN protoc --proto_path=proto/v1  --go_out=plugins=grpc:pkg/api/v1 todo-service.proto

FROM golang:1.11-alpine
WORKDIR /usr/src/app
COPY --from=0 /usr/src/app/pkg ./pkg