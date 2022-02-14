package memory

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"github.com/mcafee/generic-data-service/stores"
)

type store struct {
	mp map[string]string
}

func New() stores.Storer {
	return store{mp: make(map[string]string)}
}

func (s store) Get(ctx *gofr.Context, key string) (string, error) {
	value, ok := s.mp[key]
	if !ok {
		return "", errors.EntityNotFound{Entity: "config", ID: key}
	}

	return value, nil
}

func (s store) Set(ctx *gofr.Context, key, value string) error {
	s.mp[key] = value

	return nil
}

func (s store) Delete(ctx *gofr.Context, key string) error {
	_, ok := s.mp[key]
	if !ok {
		return errors.EntityNotFound{Entity: "config", ID: key}
	}

	delete(s.mp, key)

	return nil
}
