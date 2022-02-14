package memory

import (
	"context"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"github.com/stretchr/testify/assert"
)

const (
	key        = "id1"
	value      = "user"
	invalidKey = "id2"
)

func TestGet(t *testing.T) {
	ctx := initializeTest()

	st := store{
		mp: make(map[string]string),
	}

	st.mp[key] = value
	defer delete(st.mp, key)

	tests := []struct {
		desc  string
		key   string
		value string
		err   error
	}{
		{"correct case", key, value, nil},
		{"invalid key", invalidKey, "", errors.EntityNotFound{Entity: "config", ID: invalidKey}},
	}

	for _, tc := range tests {
		value, err := st.Get(ctx, tc.key)

		assert.Equal(t, tc.err, err)

		assert.Equal(t, tc.value, value)
	}
}

func TestSet(t *testing.T) {
	ctx := initializeTest()

	s := New()

	err := s.Set(ctx, key, value)

	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	ctx := initializeTest()
	st := store{
		mp: make(map[string]string),
	}

	st.mp[key] = value
	defer delete(st.mp, key)

	tests := []struct {
		desc  string
		key   string
		value string
		err   error
	}{
		{"correct case", key, value, nil},
		{"invalid key", invalidKey, "", errors.EntityNotFound{Entity: "config", ID: invalidKey}},
	}

	for _, tc := range tests {
		err := st.Delete(ctx, tc.key)

		assert.Equal(t, tc.err, err)
	}
}

func initializeTest() *gofr.Context {
	app := gofr.New()
	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()

	return ctx
}
