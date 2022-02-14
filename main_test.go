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

func TestHTTPClientRedis(t *testing.T) {
	t.Setenv("BACKEND_STORE", "redis")

	go main()
	time.Sleep(time.Second * 5)

	body := []byte(`{"name":"user1"}`)

	tests := []struct {
		desc       string
		config     string
		method     string
		endpoint   string
		statusCode int
		body       []byte
	}{
		{"POST success redis", "redis", http.MethodPost, "config", http.StatusCreated, body},
		{"GET success redis", "redis", http.MethodGet, "config/name", http.StatusOK, nil},
		{"Delete Success redis", "redis", http.MethodDelete, "config/name", http.StatusNoContent, []byte(`{}`)},
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

func TestHTTPClientDynamoDB(t *testing.T) {
	t.Setenv("BACKEND_STORE", "dynamodb")

	go main()
	time.Sleep(time.Second * 5)

	body := []byte(`{"name":"user1"}`)

	tests := []struct {
		desc       string
		config     string
		method     string
		endpoint   string
		statusCode int
		body       []byte
	}{
		{"POST success dynamodb", "dynamodb", http.MethodPost, "config", http.StatusCreated, body},
		{"GET success dynamodb", "dynamodb", http.MethodGet, "config/name", http.StatusOK, nil},
		{"Delete success dynamodb", "dynamodb", http.MethodDelete, "config/name", http.StatusNoContent, []byte(`{}`)},
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

func TestGRPCClientRedis(t *testing.T) {
	t.Setenv("BACKEND_STORE", "redis")

	go main()
	time.Sleep(time.Second * 5)

	conn, err := grpc.Dial("localhost:10000", grpc.WithInsecure())
	if err != nil {
		t.Errorf("did not connect: %s", err)
		return
	}

	defer conn.Close()

	c := grpc2.NewGenericDataServiceClient(conn)

	_, err = c.Set(context.TODO(), &grpc2.Data{Key: "1", Value: "user1"})
	assert.NoError(t, err)

	_, err = c.Get(context.TODO(), &grpc2.Key{Key: "1"})
	assert.NoError(t, err)

	_, err = c.Delete(context.TODO(), &grpc2.Key{Key: "1"})
	assert.NoError(t, err)
}

func TestGRPCClientDynamoDB(t *testing.T) {
	t.Setenv("BACKEND_STORE", "dynamodb")

	go main()
	time.Sleep(time.Second * 5)

	conn, err := grpc.Dial("localhost:10000", grpc.WithInsecure())
	if err != nil {
		t.Errorf("did not connect: %s", err)
		return
	}

	defer conn.Close()

	c := grpc2.NewGenericDataServiceClient(conn)

	_, err = c.Set(context.TODO(), &grpc2.Data{Key: "1", Value: "user1"})
	assert.NoError(t, err)

	_, err = c.Get(context.TODO(), &grpc2.Key{Key: "1"})
	assert.NoError(t, err)

	_, err = c.Delete(context.TODO(), &grpc2.Key{Key: "1"})
	assert.NoError(t, err)
}

func TestInvalidBackendStore(t *testing.T) {
	t.Setenv("BACKEND_STORE", "")

	go main()
	time.Sleep(time.Second * 5)

	conn, err := grpc.Dial("localhost:10000", grpc.WithInsecure())
	if err != nil {
		t.Errorf("did not connect: %s", err)
		return
	}

	defer conn.Close()

	c := grpc2.NewGenericDataServiceClient(conn)

	_, err = c.Set(context.TODO(), &grpc2.Data{Key: "1", Value: "user1"})
	assert.Error(t, err, "Expected error: Error while dialing dial tcp ")
}
