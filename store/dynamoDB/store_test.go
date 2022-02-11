package dynamoDB

import (
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/mcafee/generic-data-service/store"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	app := gofr.New()

	table := "genericStore"
	deleteTableInput := &dynamodb.DeleteTableInput{TableName: aws.String(table)}

	_, err := app.DynamoDB.DeleteTable(deleteTableInput)
	if err != nil {
		app.Logger.Errorf("error in deleting table, %v", err)
	}

	createTableInput := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("id"), AttributeType: aws.String("S")},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("id"), KeyType: aws.String("HASH")},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{ReadCapacityUnits: aws.Int64(10), WriteCapacityUnits: aws.Int64(5)},
		TableName:             aws.String(table),
	}

	_, err = app.DynamoDB.CreateTable(createTableInput)
	if err != nil {
		app.Logger.Errorf("Failed creation of table %v, %v", table, err)
	}

	os.Exit(m.Run())
}

func initializeTest(t *testing.T) (*gofr.Context, store.Storer) {
	app := gofr.New()

	// RefreshTables
	seeder := datastore.NewSeeder(&app.DataStore, "../../db")
	seeder.RefreshDynamoDB(t, "genericStore")

	ctx := gofr.NewContext(nil, nil, app)

	store := New("genericStore")

	return ctx, store
}

func TestGet(t *testing.T) {
	ctx, store := initializeTest(t)

	resp, err := store.Get(ctx, "1")
	if err != nil {
		t.Errorf("Failed\tExpected %v\nGot %v\n", nil, err)
	}

	assert.Equal(t, "Ponting", resp)
}

func TestGet_Error(t *testing.T) {
	app := gofr.New()

	ctx := gofr.NewContext(nil, nil, app)
	store := New("dummy")

	_, err := store.Get(ctx, "1")

	assert.IsType(t, errors.DB{}, err)
}

func TestCreate(t *testing.T) {
	ctx, store := initializeTest(t)

	key, value := "7", "John"

	err := store.Set(ctx, key, value)
	if err != nil {
		t.Errorf("Failed\tExpected %v\nGot %v\n", nil, err)
	}
}

func TestDelete(t *testing.T) {
	ctx, store := initializeTest(t)

	err := store.Delete(ctx, "1")
	if err != nil {
		t.Errorf("Failed\tExpected %v\nGot %v\n", nil, err)
	}
}
