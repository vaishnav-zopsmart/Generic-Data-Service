package grpc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/mcafee/generic-data-service/store"
)

func initializeTest(t *testing.T) (*store.MockStorer, *gofr.Gofr) {
	ctrl := gomock.NewController(t)
	service := store.NewMockStorer(ctrl)
	app := gofr.New()

	return service, app
}

func TestConfig_GetKey(t *testing.T) {
	mockStore, app := initializeTest(t)

	err := errors.EntityNotFound{Entity: "value", ID: "2"}

	mockStore.EXPECT().Get(gomock.Any(), "1").Return("user1", nil)
	mockStore.EXPECT().Get(gomock.Any(), "2").Return("", err)

	tests := []struct {
		desc string
		key  string
		resp *Data
		err  error
	}{
		{"success case", "1", &Data{Key: "1", Value: "user1"}, nil},
		{"error from store", "2", nil, err},
		{"missing param", "", nil, errors.MissingParam{Param: []string{"key"}}},
	}

	for i, tc := range tests {
		req := httptest.NewRequest(http.MethodGet, "http://dummy/"+tc.key, nil)
		r := request.NewHTTPRequest(req)
		ctx := gofr.NewContext(nil, r, app)

		grcpHandler := New(mockStore)

		resp, err := grcpHandler.Get(ctx, &Key{Key: tc.key})

		assert.Equal(t, tc.resp, resp, "Test[%v] failed.", i)

		assert.Equal(t, tc.err, err, "Test[%v] failed.", i)
	}
}

func TestConfig_SetKey(t *testing.T) {
	mockStore, app := initializeTest(t)

	err := errors.DB{Err: errors.Error("redis: nil")}
	mp1 := &Data{Key: "1", Value: "user1"}
	mp2 := &Data{Key: "1", Value: "abcd"}

	mockStore.EXPECT().Set(gomock.Any(), mp1.Key, mp1.Value).Return(nil)
	mockStore.EXPECT().Set(gomock.Any(), mp2.Key, mp2.Value).Return(err)

	tests := []struct {
		desc   string
		input  *Data
		output *Response
		err    error
	}{
		{"success case", mp1, &Response{Response: "Successful"}, nil},
		{"error from store", mp2, nil, err},
	}

	for _, tc := range tests {
		req := httptest.NewRequest(http.MethodPost, "http://dummy", nil)
		r := request.NewHTTPRequest(req)
		ctx := gofr.NewContext(nil, r, app)

		grcpHandler := New(mockStore)

		output, err := grcpHandler.SetKey(ctx, tc.input)

		assert.Equal(t, tc.output, output)

		assert.Equal(t, tc.err, err)
	}
}

func TestConfig_DeleteKey(t *testing.T) {
	mockStore, app := initializeTest(t)

	err := errors.EntityNotFound{Entity: "value", ID: "2"}

	mockStore.EXPECT().Delete(gomock.Any(), "1").Return(nil)
	mockStore.EXPECT().Delete(gomock.Any(), "2").Return(err)

	tests := []struct {
		desc string
		key  string
		resp *Response
		err  error
	}{
		{"success case", "1", &Response{Response: "Deleted successfully"}, nil},
		{"error from store", "2", nil, err},
		{"missing param", "", nil, errors.MissingParam{Param: []string{"key"}}},
	}

	for i, tc := range tests {
		req := httptest.NewRequest(http.MethodGet, "http://dummy/"+tc.key, nil)
		r := request.NewHTTPRequest(req)
		ctx := gofr.NewContext(nil, r, app)

		grcpHandler := New(mockStore)

		resp, err := grcpHandler.DeleteKey(ctx, &Key{Key: tc.key})

		assert.Equal(t, tc.resp, resp, "Test[%v] failed.", i)

		assert.Equal(t, tc.err, err, "Test[%v] failed.", i)
	}
}
