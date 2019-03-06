FROM golang:1.11-alpine
WORKDIR /usr/src/app
RUN apk add git
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o grpc-server .

FROM scratch
COPY --from=0 /usr/src/app/grpc-server /grpc-server
