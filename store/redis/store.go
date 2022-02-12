package redis

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"github.com/mcafee/generic-data-service/store"
)

// model is the type on which all the core layer's functionality is implemented on
type redisStore struct{}

// New is factory function for store
//nolint:revive // model should not be used without proper initilization with required dependency
func New() store.Storer {
	return redisStore{}
}

// Get returns the value for a given key, throws an error, if something goes wrong
func (r redisStore) Get(ctx *gofr.Context, key string) (string, error) {
	// fetch the Redis client
	rc := ctx.Redis

	value, err := rc.Get(ctx.Context, key).Result()
	if err != nil {
		ctx.Logger.Info(err)
		return "", errors.DB{Err: err}
	}

	return value, nil
}

// Set accepts a key-value pair, and sets those in Redis, if expiration is non-zero value, it sets a expiration(TTL)
// on those keys, if expiration is 0, then the keys have no expiration time
func (r redisStore) Set(ctx *gofr.Context, key, value string) error {
	// fetch the Redis client
	rc := ctx.Redis

	sc := rc.Set(ctx.Context, key, value, 0)
	if sc != nil && sc.Err() != nil {
		return errors.DB{Err: sc.Err()}
	}

	return nil
}

// Delete deletes a key from Redis, returns the error if it fails to delete
func (r redisStore) Delete(ctx *gofr.Context, key string) error {
	// fetch the Redis client
	rc := ctx.Redis
	return rc.Del(ctx.Context, key).Err()
}
