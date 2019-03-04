FROM golang:1.11-alpine
WORKDIR /usr/src/app
RUN apk add git
RUN go get github.com/99designs/gqlgen
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o graphql-server ./server/.

FROM scratch
COPY --from=0 /usr/src/app/graphql-server /graphql-server
