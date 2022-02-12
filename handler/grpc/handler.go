package grpc

import (
	"context"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/log"
	"github.com/mcafee/generic-data-service/store"
)

type handler struct {
	UnimplementedGenericDataServiceServer
	st store.Storer
}

// New is factory function for GRPC Handler
//nolint:revive // handler should not be used without proper initilization with required dependency
func New(s store.Storer) handler {
	return handler{
		st: s,
	}
}

func (h handler) Get(ctx context.Context, k *Key) (*Data, error) {
	if k.Key == "" {
		return nil, errors.MissingParam{Param: []string{"key"}}
	}

	c:=getContext(ctx)

	value, err := h.st.Get(c, k.Key)
	if err != nil {
		return nil, errors.EntityNotFound{Entity: "value", ID: k.Key}
	}

	resp := &Data{
		Key:   k.Key,
		Value: value,
	}

	return resp, nil
}

func (h handler) SetKey(ctx context.Context, d *Data) (*Response, error) {
	c:=getContext(ctx)

	err := h.st.Set(c, d.Key, d.Value)
	if err != nil {
		return nil, err
	}

	resp := &Response{
		Response: "Successful",
	}

	return resp, nil
}


func (h handler) DeleteKey(ctx context.Context, k *Key) (*Response, error) {
	c:=getContext(ctx)

	if k.Key==""{
		return nil,errors.MissingParam{Param: []string{"key"}}
	}

	err := h.st.Delete(c, k.Key)
	if err != nil {
		return nil, err
	}

	resp := &Response{
		Response: "Deleted successfully",
	}

	return resp, nil
}

func getContext(ctx context.Context) *gofr.Context{
	logger:=log.NewCorrelationLogger("")
	return &gofr.Context{Context:ctx,Logger: logger}
}