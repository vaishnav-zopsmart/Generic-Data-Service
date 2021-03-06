package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/mcafee/generic-data-service/stores"
)

func initializeTest(t *testing.T) (*stores.MockStorer, *gofr.Gofr) {
	ctrl := gomock.NewController(t)
	service := stores.NewMockStorer(ctrl)
	app := gofr.New()

	return service, app
}

func getContext(method string, body []byte, pathParams map[string]string, app *gofr.Gofr) *gofr.Context {
	const url = "/config"

	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, url, bytes.NewReader(body))
	req := request.NewHTTPRequest(r)
	res := responder.NewContextualResponder(w, r)
	ctx := gofr.NewContext(res, req, app)

	ctx.SetPathParams(pathParams)

	return ctx
}

func TestConfig_GetKey(t *testing.T) {
	mockStore, app := initializeTest(t)

	err := errors.DB{Err: errors.Error("redis: nil")}

	op := map[string]string{"1": "user1"}

	mockStore.EXPECT().Get(gomock.Any(), "1").Return("user1", nil)
	mockStore.EXPECT().Get(gomock.Any(), "2").Return("", err)

	tests := []struct {
		desc   string
		key    string
		output interface{}
		err    error
	}{
		{"success case", "1", op, nil},
		{"error from stores", "2", nil, err},
		{"missing param", "", nil, errors.MissingParam{Param: []string{"key"}}},
	}

	for _, tc := range tests {
		h := New(mockStore)

		ctx := getContext(http.MethodGet, nil, map[string]string{"key": tc.key}, app)

		output, err := h.Get(ctx)

		assert.Equal(t, tc.output, output)

		assert.Equal(t, tc.err, err)
	}
}

func TestConfig_SetKey(t *testing.T) {
	mockStore, app := initializeTest(t)

	err := errors.DB{Err: errors.Error("redis: nil")}
	mp1 := map[string]string{"1": "user1"}
	mp2 := map[string]string{"2": "abcd"}

	mockStore.EXPECT().Set(gomock.Any(), "1", "user1").Return(nil)
	mockStore.EXPECT().Set(gomock.Any(), "2", "abcd").Return(err)

	tests := []struct {
		desc   string
		input  map[string]string
		output interface{}
		err    error
	}{
		{"success case", mp1, "Successful", nil},
		{"error from stores", mp2, nil, err},
	}

	for _, tc := range tests {
		h := New(mockStore)

		body, err := json.Marshal(tc.input)
		if err != nil {
			t.Errorf("Received unexpected error:\n%+v", err)

			return
		}

		ctx := getContext(http.MethodPost, body, nil, app)

		output, err := h.Set(ctx)

		assert.Equal(t, tc.output, output)

		assert.Equal(t, tc.err, err)
	}
}

func TestConfig_SetKeyBindError(t *testing.T) {
	mockStore, app := initializeTest(t)

	expErr := errors.InvalidParam{Param: []string{"body"}}

	input := map[string]int{"1": 1}

	h := New(mockStore)

	body, err := json.Marshal(input)
	if err != nil {
		t.Errorf("Received unexpected error:\n%+v", err)

		return
	}

	ctx := getContext(http.MethodPost, body, nil, app)

	output, err := h.Set(ctx)

	assert.Nil(t, output)

	assert.Equal(t, expErr, err)
}

func TestConfig_DeleteKey(t *testing.T) {
	mockStore, app := initializeTest(t)

	err := errors.DB{Err: errors.Error("redis: nil")}

	mockStore.EXPECT().Delete(gomock.Any(), "1").Return(nil)
	mockStore.EXPECT().Delete(gomock.Any(), "2").Return(err)

	tests := []struct {
		desc   string
		key    string
		output interface{}
		err    error
	}{
		{"success case", "1", "Deleted successfully", nil},
		{"error from stores", "2", nil, err},
		{"missing param", "", nil, errors.MissingParam{Param: []string{"key"}}},
	}

	for _, tc := range tests {
		h := New(mockStore)

		ctx := getContext(http.MethodDelete, nil, map[string]string{"key": tc.key}, app)

		output, err := h.Delete(ctx)

		assert.Equal(t, tc.output, output)

		assert.Equal(t, tc.err, err)
	}
}
