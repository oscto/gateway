package main

import (
	"context"
	"go-micro.dev/v4/api"
	"log"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := api.NewApi()
	if err := srv.Run(ctx); err != nil {
		log.Fatal("failed to server run ", err)
	}

}
