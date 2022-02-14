package grpc

import (
	"context"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"github.com/mcafee/generic-data-service/stores"
)

type handler struct {
	UnimplementedGenericDataServiceServer
	s   stores.Storer
	app *gofr.Gofr
}

// New is factory function for GRPC Handler
//nolint:revive // handler should not be used without proper initilization with required dependency
func New(s stores.Storer, app *gofr.Gofr) handler {
	return handler{
		s:   s,
		app: app,
	}
}

func (h handler) Get(ctx context.Context, k *Key) (*Data, error) {
	if k.Key == "" {
		return nil, errors.MissingParam{Param: []string{"key"}}
	}

	c := getContext(ctx, h.app)

	value, err := h.s.Get(c, k.Key)
	if err != nil {
		return nil, errors.EntityNotFound{Entity: "value", ID: k.Key}
	}

	resp := &Data{
		Key:   k.Key,
		Value: value,
	}

	return resp, nil
}

func (h handler) Set(ctx context.Context, d *Data) (*Response, error) {
	c := getContext(ctx, h.app)

	err := h.s.Set(c, d.Key, d.Value)
	if err != nil {
		return nil, err
	}

	resp := &Response{
		Response: "Successful",
	}

	return resp, nil
}

func (h handler) Delete(ctx context.Context, k *Key) (*Response, error) {
	c := getContext(ctx, h.app)

	if k.Key == "" {
		return nil, errors.MissingParam{Param: []string{"key"}}
	}

	err := h.s.Delete(c, k.Key)
	if err != nil {
		return nil, err
	}

	resp := &Response{
		Response: "Deleted successfully",
	}

	return resp, nil
}

// getContext returns a gofr Context
func getContext(ctx context.Context, app *gofr.Gofr) *gofr.Context {
	logger := app.Logger

	return &gofr.Context{Context: ctx, Logger: logger, Gofr: app}
}
