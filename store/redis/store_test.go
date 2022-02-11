package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)


func initializeTest(t *testing.T) (*gofr.Gofr,*gofr.Context){
	app := gofr.New()
	c := gofr.NewContext(nil, nil, app)
	c.Context = context.Background()

	// initializing the seeder
	seeder := datastore.NewSeeder(&app.DataStore, "../../db")
	seeder.RefreshRedis(t, "store")

	return app,c
}


func TestSetWithError(t *testing.T) {
	app,c:=initializeTest(t)

	app.Redis.Close()

	expected := "redis: client is closed"
	store := New()

	resp := store.Set(c, "key", "value")

	assert.Equal(t, expected, resp.Error())
}

func TestSet(t *testing.T) {
	_,c:=initializeTest(t)
	store := New()

	err := store.Set(c, "someKey123", "someValue123")
	if err != nil {
		t.Errorf("FAILED, Expected no error, Got: %v", err)
	}
}

func TestGet(t *testing.T) {
	_,c:=initializeTest(t)
	tests := []struct {
		desc string
		key  string
		resp string
		err  error
	}{
		{"get success", "someKey123", "someValue123", nil},
		{"get fail", "someKey", "", errors.DB{}},
	}

	for i, tc := range tests {
		store := New()
		resp, err := store.Get(c, tc.key)

		assert.Equal(t, tc.resp, resp, "TEST[%d], failed.\n%s", i, tc.desc)

		assert.IsType(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestDelete(t *testing.T) {
	_,c:=initializeTest(t)
	tests := []struct {
		desc string
		key  string
		err  error
	}{
		{"delete success", "someKey123", nil},
	}

	for i, tc := range tests {
		store := New()
		err := store.Delete(c, tc.key)

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}