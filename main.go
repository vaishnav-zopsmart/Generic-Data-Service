package main

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"github.com/mcafee/generic-data-service/handler"
	"github.com/mcafee/generic-data-service/handler/grpc"
	"github.com/mcafee/generic-data-service/stores"
	"github.com/mcafee/generic-data-service/stores/dynamodb"
	"github.com/mcafee/generic-data-service/stores/redis"
)

func main() {
	app := gofr.New()

	var s stores.Storer

	backendStore := app.Config.Get("BACKEND_STORE")

	switch backendStore {
	case "redis":
		s = redis.New()
	case "dynamodb":
		table := app.Config.Get("DYNAMODB_TABLE")
		s = dynamodb.New(table)
	default:
		return
	}

	h := handler.New(s)

	// Specifying the different routes supported by this services
	app.GET("/config/{key}", h.Get)
	app.POST("/config", h.Set)
	app.DELETE("/config/{key}", h.Delete)

	grpcHandler := grpc.New(s, app)
	grpc.RegisterGenericDataServiceServer(app.Server.GRPC.Server(), grpcHandler)

	app.Start()
}
