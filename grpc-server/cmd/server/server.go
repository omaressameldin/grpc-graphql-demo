package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/omaressameldin/grpc-graphql-demo/grpc-server/pkg/service/v1"
	"github.com/omaressameldin/grpc-graphql-demo/grpc-server/server"
)

// Config is configuration for Server
type Config struct {
	// gRPC server start parameters section
	// gRPC is TCP port to listen by gRPC server
	GRPCPort     string
	DatabasePath string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// get configuration
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfg.DatabasePath, "db-path", "", "Database path")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	v1API := v1.NewToDoServiceServer(cfg.DatabasePath)

	return server.RunServer(ctx, v1API, cfg.GRPCPort)
}

// CloseServer closes all connections such as database connection
func CloseServer() error {
	return v1.CloseConnection()
}
