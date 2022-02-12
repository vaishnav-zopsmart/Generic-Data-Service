package dynamodb

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/mcafee/generic-data-service/store"
)

type dynamoDBstore struct {
	tableName string
}

// New factory function for person store
func New(t string) store.Storer {
	return dynamoDBstore{tableName: t}
}

func (s dynamoDBstore) Get(ctx *gofr.Context, key string) (string, error) {
	input := &dynamodb.GetItemInput{
		AttributesToGet: []*string{aws.String("value")},
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(key)},
		},

		TableName: aws.String(s.tableName),
	}

	out, err := ctx.DynamoDB.GetItem(input)
	if err != nil {
		return "", errors.DB{Err: err}
	}

	v := struct {
		Value string `json:"value"`
	}{}

	err = dynamodbattribute.UnmarshalMap(out.Item, &v)
	if err != nil {
		return "", errors.DB{Err: err}
	}

	return v.Value, nil
}

func (s dynamoDBstore) Set(ctx *gofr.Context, key, value string) error {
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"id":    {S: aws.String(key)},
			"value": {S: aws.String(value)},
		},
		TableName: aws.String(s.tableName),
	}

	_, err := ctx.DynamoDB.PutItem(input)

	return err
}

func (s dynamoDBstore) Delete(ctx *gofr.Context, key string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(key)},
		},
		TableName: aws.String(s.tableName),
	}

	_, err := ctx.DynamoDB.DeleteItem(input)

	return err
}
