FROM golang:1.11-alpine
WORKDIR /usr/src/app
RUN apk add git
RUN go get github.com/boltdb/bolt/...
RUN go get github.com/golang/protobuf/ptypes
RUN go get google.golang.org/grpc
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o grpc-server .

FROM scratch
COPY --from=0 /usr/src/app/grpc-server /grpc-server
