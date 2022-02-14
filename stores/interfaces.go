package stores

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

// Storer is an abstraction for the core layer
type Storer interface {
	Get(ctx *gofr.Context, key string) (string, error)
	Set(ctx *gofr.Context, key, value string) error
	Delete(ctx *gofr.Context, key string) error
}
