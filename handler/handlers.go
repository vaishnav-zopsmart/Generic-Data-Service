package handler

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"github.com/mcafee/generic-data-service/store"
)

type config struct {
	st store.Storer
}

// New is factory function for config
//nolint:revive // handler should not be used without proper initilization with required dependency
func New(s store.Storer) config {
	return config{
		st: s,
	}
}

// GetKey is a handler function of type gofr.Handler, it fetches keys
func (c config) GetKey(ctx *gofr.Context) (interface{}, error) {
	// fetch the path parameter as specified in the route
	key := ctx.PathParam("key")
	if key == "" {
		return nil, errors.MissingParam{Param: []string{"key"}}
	}

	value, err := c.st.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	resp := make(map[string]string)
	resp[key] = value

	return resp, nil
}

// SetKey is a handler function of type gofr.Handler, it sets keys
func (c config) SetKey(ctx *gofr.Context) (interface{}, error) {
	input := make(map[string]string)

	err := ctx.Bind(&input)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	for key, value := range input {
		err = c.st.Set(ctx, key, value)
		if err != nil {
			return nil, err
		}
	}

	return "Successful", nil
}

// DeleteKey is a handler function of type gofr.Handler, it deletes keys
func (c config) DeleteKey(ctx *gofr.Context) (interface{}, error) {
	// fetch the path parameter as specified in the route
	key := ctx.PathParam("key")
	if key == "" {
		return nil, errors.MissingParam{Param: []string{"key"}}
	}

	if err := c.st.Delete(ctx, key); err != nil {
		ctx.Logger.Errorf("err: ", err)
		return nil, err
	}

	return "Deleted successfully", nil
}
