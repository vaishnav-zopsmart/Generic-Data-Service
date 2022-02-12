package main

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/mcafee/generic-data-service/handler"
	"github.com/mcafee/generic-data-service/handler/grpc"
	"github.com/mcafee/generic-data-service/store"
	"github.com/mcafee/generic-data-service/store/dynamoDB"
	"github.com/mcafee/generic-data-service/store/redis"
)

func main() {
	app := gofr.New()

	var s store.Storer

	backendStore := app.Config.Get("BACKEND_STORE")

	switch backendStore {
	case "redis":
		s = redis.New()
	case "dynamoDB":
		table := app.Config.Get("DYNAMODB_TABLE")
		s = dynamoDB.New(table)
	default:
		return
	}

	h := handler.New(s)

	// Specifying the different routes supported by this services
	app.GET("/config/{key}", h.GetKey)
	app.POST("/config", h.SetKey)
	app.DELETE("/config/{key}", h.DeleteKey)

	grpcHandler := grpc.New(s)
	grpc.RegisterGenericDataServiceServer(app.Server.GRPC.Server(), grpcHandler)

	app.Start()

}
