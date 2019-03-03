FROM znly/protoc:0.4.0
WORKDIR /usr/src/app
COPY ./proto ./proto
RUN mkdir -p pkg/api/v1
RUN protoc --proto_path=proto/v1  --go_out=plugins=grpc:pkg/api/v1 todo-service.proto

FROM golang:1.11-alpine
WORKDIR /usr/src/app
RUN apk add git
RUN go get github.com/boltdb/bolt/...
RUN go get github.com/golang/protobuf/ptypes
RUN go get google.golang.org/grpc
COPY . .
COPY --from=0 /usr/src/app/pkg/api ./pkg/api
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o grpc-server .

FROM scratch
COPY --from=1 /usr/src/app/grpc-server /grpc-server
