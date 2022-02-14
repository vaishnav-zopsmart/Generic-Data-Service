package main

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"google.golang.org/grpc"

	"developer.zopsmart.com/go/gofr/pkg/gofr/request"

	grpc2 "github.com/mcafee/generic-data-service/handler/grpc"
)

func TestIntegration(t *testing.T) {
	body := []byte(`{"name":"user1"}`)

	go main()
	time.Sleep(time.Second * 5)

	tests := []struct {
		desc       string
		method     string
		endpoint   string
		statusCode int
		body       []byte
	}{
		{"get success", http.MethodPost, "config", http.StatusCreated, body},
		{"get non existent entity", http.MethodGet, "config/name", http.StatusOK, nil},
		{"unregistered update route", http.MethodDelete, "config/name", http.StatusNoContent, []byte(`{}`)},
	}

	for i, tc := range tests {
		req, _ := request.NewMock(tc.method, "http://localhost:9098/"+tc.endpoint, bytes.NewBuffer(tc.body))
		c := http.Client{}

		resp, err := c.Do(req)
		if err != nil {
			t.Errorf("TEST[%v] Failed.\tHTTP request encountered Err: %v\n%s", i, err, tc.desc)
			continue
		}

		if resp.StatusCode != tc.statusCode {
			t.Errorf("TEST[%v] Failed.\tExpected %v\tGot %v\n%s", i, tc.statusCode, resp.StatusCode, tc.desc)
		}

		_ = resp.Body.Close()
	}
}

func TestGRPCClient(t *testing.T) {
	go main()
	time.Sleep(time.Second * 5)

	conn, err := grpc.Dial("localhost:10000", grpc.WithInsecure())
	if err != nil {
		t.Errorf("did not connect: %s", err)
		return
	}

	defer conn.Close()

	c := grpc2.NewGenericDataServiceClient(conn)

	grpc.NewServer()

	_, err = c.SetKey(context.TODO(), &grpc2.Data{Key: "1", Value: "user1"})
	assert.NoError(t, err)

	_, err = c.GetKey(context.TODO(), &grpc2.Key{Key: "1"})
	assert.NoError(t, err)

	_, err = c.DeleteKey(context.TODO(), &grpc2.Key{Key: "1"})
	assert.NoError(t, err)
}
