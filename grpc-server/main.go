package main

import (
	"fmt"
	"os"

	cmd "./cmd/server"
)

func main() {
	if err := cmd.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cmd.CloseServer()
}
