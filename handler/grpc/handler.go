package grpc

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
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

func (h handler) Get(ctx *gofr.Context, k *Key) (*Response, error) {
	value, err := h.st.Get(ctx, k.Key)
	if err != nil {
		return nil, errors.EntityNotFound{Entity: "value", ID: k.Key}
	}

	resp := &Response{
		Response: value,
	}

	return resp, nil
}

func (h handler) SetKey(ctx *gofr.Context, d *Data) (*Response, error) {
	err := h.st.Set(ctx, d.Key, d.Value)
	if err != nil {
		return nil, err
	}

	resp := &Response{
		Response: "Successful",
	}

	return resp, nil
}

func (h handler) DeleteKey(ctx *gofr.Context, k *Key) (*Response, error) {
	err := h.st.Delete(ctx, k.Key)
	if err != nil {
		return nil, err
	}

	resp := &Response{
		Response: "Deleted successfully",
	}

	return resp, nil
}
