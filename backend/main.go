package main

import (
	"context"
	"dotnbox/server"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	server := server.New(logger)
	ctx := context.Background()
	server.Run(ctx, ":8080")
}
