package main

import (
	"api-service/servers"
	"api-service/utils"
	"context"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()
	logger := utils.SugarLogger()

	group, _ := errgroup.WithContext(ctx)
	ginServer := servers.NewGinServer()
	if err := ginServer.Serve(group); err != nil {
		logger.Errorf("gin server error: %v", err)
	}

	if err := group.Wait(); err != nil {
		logger.Errorf("group error: %v", err)
	}
}
